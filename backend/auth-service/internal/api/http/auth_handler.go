package http

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	gojwt "github.com/golang-jwt/jwt/v5" // Added for direct access to golang-jwt types
	"github.com/russian-steam/auth-service/internal/domain"
	"github.com/russian-steam/auth-service/internal/pkg/kafka" // Added Kafka
	"github.com/russian-steam/auth-service/internal/pkg/jwt"   // For local jwt package errors (ErrTokenExpired)
	"github.com/russian-steam/auth-service/internal/service"
)

type AuthHTTPHandler struct {
	userService    *service.UserService
	tokenService   *service.TokenService
	verifyService  *service.VerificationService
	eventPublisher kafka.EventPublisher // Added EventPublisher
}

func NewAuthHTTPHandler(us *service.UserService, ts *service.TokenService, vs *service.VerificationService, ep kafka.EventPublisher) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		userService:   us,
		tokenService:  ts,
		verifyService: vs,
		eventPublisher: ep, // Store EventPublisher
	}
}

func (h *AuthHTTPHandler) RegisterRoutes(router *gin.Engine) {
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/register", h.Register)
		authGroup.POST("/login", h.Login)
		authGroup.POST("/refresh-token", h.RefreshToken)
		authGroup.POST("/logout", h.Logout) // Requires auth
		authGroup.POST("/verify-email", h.VerifyEmail)
	}
}

// Helper to format validation errors
func formatValidationErrors(err error) []APIError {
    var ve validator.ValidationErrors
    if errors.As(err, &ve) {
        out := make([]APIError, len(ve))
        for i, fe := range ve {
            out[i] = APIError{
                Code:  "VALIDATION_ERROR",
                Title: "Validation failed for field '" + fe.Field() + "' with rule '" + fe.Tag() + "'",
                Source: &ErrorSource{Pointer: "/data/attributes/" + strings.ToLower(fe.Field())}, // Example, adjust as needed
            }
        }
        return out
    }
    return []APIError{{Code: "INVALID_REQUEST_BODY", Title: "Invalid request body", Detail: err.Error()}}
}


func (h *AuthHTTPHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Errors: formatValidationErrors(err)})
		return
	}

	user, err := h.userService.RegisterUser(req.Username, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, ErrorResponse{Errors: []APIError{{Code: "USER_ALREADY_EXISTS", Title: "User already exists"}}})
		} else {
			log.Printf("Error registering user: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Errors: []APIError{{Code: "INTERNAL_SERVER_ERROR", Title: "Failed to register user"}}})
		}
		return
	}

    // Publish UserRegisteredEvent
    registeredEventData := kafka.UserRegisteredEventData{
        UserID:                user.ID,
        Username:              user.Username,
        Email:                 user.Email,
        Status:                string(user.Status),
        RegistrationTimestamp: user.CreatedAt,
    }
    if errPub := h.eventPublisher.Publish(c.Request.Context(), "com.platform.auth.user.registered.v1", registeredEventData); errPub != nil {
        log.Printf("Failed to publish user.registered event for user %s: %v", user.ID, errPub)
        // Non-critical for registration flow, log and continue
    }

    // Generate and publish EmailVerificationRequestedEvent
    // Assuming VerificationCodeExpiry is defined in service package or accessible
    // For this diff, we'll use a local example value if not directly available from service constants.
    // In a real scenario, service.VerificationCodeExpiry should be used.
    const verificationCodeDefaultExpiry = 15 * time.Minute
    verificationCodeActualExpiry := verificationCodeDefaultExpiry // Replace with actual constant if available

    _, rawVerificationCode, err := h.verifyService.GenerateEmailVerificationCode(user.ID, user.Email)
    if err != nil {
        log.Printf("Failed to generate verification code for user %s: %v", user.ID, err)
    } else {
        emailVerificationEventData := kafka.EmailVerificationRequestedEventData{
            UserID:           user.ID,
            Email:            user.Email,
            VerificationCode: rawVerificationCode,
            ExpiresAt:        time.Now().UTC().Add(verificationCodeActualExpiry),
        }
        if errPub := h.eventPublisher.Publish(c.Request.Context(), "com.platform.auth.email.verification_requested.v1", emailVerificationEventData); errPub != nil {
            log.Printf("Failed to publish email.verification_requested event for user %s: %v", user.ID, errPub)
        }
    }

	c.JSON(http.StatusCreated, RegisterResponse{
        Data: UserResponse{
            ID:        user.ID,
            Username:  user.Username,
            Email:     user.Email,
            Status:    string(user.Status),
            CreatedAt: user.CreatedAt,
        },
        Meta: MessageMeta{Message: "Registration successful. Please check your email to verify your account."},
    })
}

func (h *AuthHTTPHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Errors: formatValidationErrors(err)})
		return
	}

	user, err := h.userService.AuthenticateUser(req.Login, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) || errors.Is(err, service.ErrUserNotFound) || errors.Is(err, service.ErrPasswordVerification) {
			c.JSON(http.StatusUnauthorized, ErrorResponse{Errors: []APIError{{Code: "INVALID_CREDENTIALS", Title: "Invalid login or password"}}})
		} else {
			log.Printf("Error authenticating user: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Errors: []APIError{{Code: "INTERNAL_SERVER_ERROR", Title: "Failed to authenticate user"}}})
		}
		return
	}

    if user.Status == domain.StatusPendingVerification {
         c.JSON(http.StatusForbidden, ErrorResponse{Errors: []APIError{{Code: "EMAIL_NOT_VERIFIED", Title: "Email not verified", Detail: "Please verify your email before logging in."}}})
         return
    }
    if user.Status != domain.StatusActive {
         c.JSON(http.StatusForbidden, ErrorResponse{Errors: []APIError{{Code: "ACCOUNT_NOT_ACTIVE", Title: "Account not active", Detail: "Your account is currently " + string(user.Status)}}})
         return
    }

	accessToken, rawRefreshToken, accessTokenJTI, accessExp, _, err := h.tokenService.GenerateTokens(user) // Assuming GenerateTokens returns JTI & refresh token expiry
	if err != nil {
		log.Printf("Error generating tokens for user %s: %v", user.ID, err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Errors: []APIError{{Code: "TOKEN_GENERATION_FAILED", Title: "Failed to generate tokens"}}})
		return
	}

    // Publish UserLoginSuccessEvent
    loginSuccessData := kafka.UserLoginSuccessEventData{
        UserID:         user.ID,
        SessionID:      accessTokenJTI, // Or Refresh Token JTI if that's your session identifier
        IPAddress:      c.ClientIP(),
        UserAgent:      c.Request.UserAgent(),
        LoginTimestamp: time.Now().UTC(),
    }
    if errPub := h.eventPublisher.Publish(c.Request.Context(), "com.platform.auth.user.login_success.v1", loginSuccessData); errPub != nil {
        log.Printf("Failed to publish user.login_success event for user %s: %v", user.ID, errPub)
    }

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(time.Until(accessExp).Seconds()),
        RefreshToken: rawRefreshToken,
		UserID:       user.ID,
		Username:     user.Username,
		Roles:        []string{"user"},
	})
}

func (h *AuthHTTPHandler) RefreshToken(c *gin.Context) {
    var req RefreshTokenRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{Errors: formatValidationErrors(err)})
        return
    }

    rtClaims, err := h.tokenService.JwtService().ValidateToken(req.RefreshToken)
    if err != nil {
        code := "INVALID_REFRESH_TOKEN"
        title := "Refresh token is invalid or expired"
        if errors.Is(err, jwt.ErrTokenExpired) {
            code = "REFRESH_TOKEN_EXPIRED"
            title = "Refresh token has expired"
        }
        c.JSON(http.StatusUnauthorized, ErrorResponse{Errors: []APIError{{Code: code, Title: title}}})
        return
    }

    user, err := h.userService.GetUserByID(rtClaims.UserID)
    if err != nil {
        log.Printf("Refresh: User %s not found from refresh token: %v", rtClaims.UserID, err)
        c.JSON(http.StatusUnauthorized, ErrorResponse{Errors: []APIError{{Code: "USER_NOT_FOUND", Title: "User associated with token not found"}}})
        return
    }

    newAccessToken, _, _, newAccessExp, _, err := h.tokenService.ValidateAndRefreshTokens(req.RefreshToken, user)
    if err != nil {
        log.Printf("Error validating/refreshing token: %v", err)
        code := "TOKEN_REFRESH_FAILED"
        title := "Failed to refresh token"
        if errors.Is(err, service.ErrRefreshTokenNotFound) { code = "REFRESH_TOKEN_NOT_FOUND"; title = "Refresh token not found"}
        if errors.Is(err, service.ErrRefreshTokenRevoked) { code = "REFRESH_TOKEN_REVOKED"; title = "Refresh token revoked"}
        if errors.Is(err, service.ErrRefreshTokenExpired) { code = "REFRESH_TOKEN_EXPIRED"; title = "Refresh token expired"}

        c.JSON(http.StatusUnauthorized, ErrorResponse{Errors: []APIError{{Code: code, Title: title}}})
        return
    }

    c.JSON(http.StatusOK, RefreshTokenResponse{
        AccessToken: newAccessToken,
        TokenType:   "Bearer",
        ExpiresIn:   int64(time.Until(newAccessExp).Seconds()),
    })
}

func (h *AuthHTTPHandler) Logout(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        c.JSON(http.StatusUnauthorized, ErrorResponse{Errors: []APIError{{Code: "MISSING_AUTH_HEADER", Title: "Authorization header is missing"}}})
        return
    }
    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
        c.JSON(http.StatusUnauthorized, ErrorResponse{Errors: []APIError{{Code: "INVALID_AUTH_HEADER", Title: "Authorization header is malformed"}}})
        return
    }
    accessToken := parts[1]

    claims, err := h.tokenService.JwtService().ValidateToken(accessToken)

    // We want to blacklist even if expired, as client clock might be off or it might be replayed.
    // So, we only care about structural validity to get JTI and Expiry.
    // If ValidateToken returns ErrTokenExpired, claims might still be populated.
    // If it's another error (e.g. malformed), claims might be nil.

    jtiToBlacklist := ""
    expiryOfTokenToBlacklist := time.Now().Add(1 * time.Minute) // Default short expiry if real one not found

    if claims != nil { // claims were parsed
        jtiToBlacklist = claims.ID
        if claims.ExpiresAt != nil && !claims.ExpiresAt.Time.IsZero() {
             expiryOfTokenToBlacklist = claims.ExpiresAt.Time
        }
    } else if errors.Is(err, jwt.ErrTokenExpired) { // Check if our local jwt service returned ErrTokenExpired
        // Token was expired, but ValidateToken didn't return claims. Try parsing unverified.
        // Use gojwt.Parser here
        parsedToken, _, parseErr := new(gojwt.Parser).ParseUnverified(accessToken, &gojwt.RegisteredClaims{})
        if parseErr == nil {
            if unverifiedClaims, ok := parsedToken.Claims.(*gojwt.RegisteredClaims); ok {
                jtiToBlacklist = unverifiedClaims.ID
                if unverifiedClaims.ExpiresAt != nil && !unverifiedClaims.ExpiresAt.Time.IsZero() {
                    expiryOfTokenToBlacklist = unverifiedClaims.ExpiresAt.Time
                }
            }
        }
    } else if err != nil { // Token is invalid for reasons other than expiry, and claims are nil
        log.Printf("Logout: Invalid access token presented: %v", err)
        c.JSON(http.StatusUnauthorized, ErrorResponse{Errors: []APIError{{Code: "INVALID_ACCESS_TOKEN", Title: "Invalid access token for logout"}}})
        return
    }


    if jtiToBlacklist != "" {
        if err := h.tokenService.BlacklistAccessToken(jtiToBlacklist, expiryOfTokenToBlacklist); err != nil {
            log.Printf("Failed to blacklist access token JTI %s: %v", jtiToBlacklist, err)
            // Non-fatal for logout.
        }
    } else {
        log.Println("Logout: Could not extract JTI from access token for blacklisting.")
    }

    // For MVP, client is expected to discard refresh token.
    // A more complete solution might involve revoking the specific refresh token if sent via cookie/body.

    // Publish SessionRevokedEvent
    if claims != nil && claims.UserID != "" && jtiToBlacklist != "" { // Ensure we have necessary info
        sessionRevokedData := kafka.SessionRevokedEventData{
            UserID:           claims.UserID,
            SessionID:        jtiToBlacklist,
            RevocationReason: "user_logout",
            RevokedAt:        time.Now().UTC(),
        }
        if errPub := h.eventPublisher.Publish(c.Request.Context(), "com.platform.auth.session.revoked.v1", sessionRevokedData); errPub != nil {
            log.Printf("Failed to publish session.revoked event for JTI %s: %v", jtiToBlacklist, errPub)
        }
    }


    c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHTTPHandler) VerifyEmail(c *gin.Context) {
    var req VerifyEmailRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{Errors: formatValidationErrors(err)})
        return
    }

    user, err := h.verifyService.VerifyEmailCode(req.VerificationCode, req.Email)
    if err != nil {
        code := "VERIFICATION_FAILED"
        title := "Email verification failed"
        status := http.StatusBadRequest
        if errors.Is(err, service.ErrVerificationCodeNotFound) { code = "VERIFICATION_CODE_NOT_FOUND"; title = "Verification code not found or already used."}
        if errors.Is(err, service.ErrVerificationCodeExpired) { code = "VERIFICATION_CODE_EXPIRED"; title = "Verification code has expired."}
        if errors.Is(err, service.ErrVerificationCodeInvalid) { code = "VERIFICATION_CODE_INVALID"; title = "Invalid verification code."}

        log.Printf("Email verification failed for %s: %v", req.Email, err)
        c.JSON(status, ErrorResponse{Errors: []APIError{{Code: code, Title: title}}})
        return
    }

    log.Printf("User %s email %s verified successfully.", user.ID, user.Email)

    // Publish UserEmailVerifiedEvent
    emailVerifiedEventData := kafka.UserEmailVerifiedEventData{
        UserID:                user.ID,
        Email:                 user.Email,
        VerificationTimestamp: time.Now().UTC(),
    }
    if errPub := h.eventPublisher.Publish(c.Request.Context(), "com.platform.auth.user.email_verified.v1", emailVerifiedEventData); errPub != nil {
        log.Printf("Failed to publish user.email_verified event for user %s: %v", user.ID, errPub)
    }

    c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully. You can now log in."})
}

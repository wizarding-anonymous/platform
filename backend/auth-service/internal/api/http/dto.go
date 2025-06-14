package http

import "time"

// Standard error response structure
type ErrorResponse struct {
	Errors []APIError `json:"errors"`
}

type APIError struct {
	Code    string          `json:"code"`
	Title   string          `json:"title"`
	Detail  string          `json:"detail,omitempty"`
	Source  *ErrorSource    `json:"source,omitempty"`
}

type ErrorSource struct {
	Pointer string `json:"pointer,omitempty"` // JSON Pointer (RFC6901) to the offending part of the request document
}


// --- Request DTOs ---
type RegisterRequest struct {
	Username        string `json:"username" binding:"required,min=3,max=32"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=8,max=72"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type LoginRequest struct {
	Login    string `json:"login" binding:"required"` // Can be username or email
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type VerifyEmailRequest struct {
	VerificationCode string `json:"verification_code" binding:"required,len=6"` // Assuming 6 char code
	Email            string `json:"email" binding:"required,email"` // To ensure code is for this email
}

// --- Response DTOs ---
type UserResponse struct {
	ID              string    `json:"id"`
	Username        string    `json:"username"`
	Email           string    `json:"email"`
	Status          string    `json:"status"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type RegisterResponse struct {
	Data UserResponse `json:"data"`
	Meta MessageMeta  `json:"meta"`
}

type MessageMeta struct {
	Message string `json:"message"`
}

type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"` // Usually "Bearer"
	ExpiresIn    int64     `json:"expires_in"` // Seconds until access token expiry
	RefreshToken string    `json:"refresh_token,omitempty"` // May be sent in HttpOnly cookie instead
	UserID       string    `json:"user_id"`
	Username     string    `json:"username"`
	Roles        []string  `json:"roles"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"` // Seconds until new access token expiry
}

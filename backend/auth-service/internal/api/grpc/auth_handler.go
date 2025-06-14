package grpc

import (
	"context"
	"log"
    // "time" // No longer directly used, timestamppb handles it
    "errors" // Added for errors.Is

	"github.com/russian-steam/auth-service/internal/service"
	"github.com/russian-steam/auth-service/internal/pkg/jwt"
	pb "github.com/russian-steam/auth-service/proto/auth/v1" // Alias for generated proto
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthGRPCHandler struct {
	pb.UnimplementedAuthServiceServer // Embed for forward compatibility
	tokenService *service.TokenService
}

func NewAuthGRPCHandler(ts *service.TokenService) *AuthGRPCHandler {
	return &AuthGRPCHandler{tokenService: ts}
}

func (h *AuthGRPCHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	if req.AccessToken == "" {
		return &pb.ValidateTokenResponse{IsValid: false, ErrorMessage: "Access token is required"}, nil
	}

	claims, err := h.tokenService.JwtService().ValidateToken(req.AccessToken)
	if err != nil {
		log.Printf("ValidateToken: Token validation failed: %v", err)
		errorMessage := "Invalid token"
        if errors.Is(err, jwt.ErrTokenExpired) {
            errorMessage = "Token expired"
        }
		return &pb.ValidateTokenResponse{IsValid: false, ErrorMessage: errorMessage}, nil
	}

	// Check JTI blacklist
	isBlacklisted, err := h.tokenService.IsAccessTokenBlacklisted(claims.ID)
	if err != nil {
		log.Printf("ValidateToken: Failed to check JTI blacklist: %v", err)
		// Consider returning gRPC internal error status
		return &pb.ValidateTokenResponse{IsValid: false, ErrorMessage: "Internal server error while checking token status"}, nil
	}
	if isBlacklisted {
		log.Printf("ValidateToken: Token JTI %s is blacklisted", claims.ID)
		return &pb.ValidateTokenResponse{IsValid: false, ErrorMessage: "Token has been invalidated (logged out)"}, nil
	}

    var expiresAtProto *timestamppb.Timestamp
    if claims.ExpiresAt != nil {
        // Ensure claims.ExpiresAt.Time is not zero before converting
        if !claims.ExpiresAt.Time.IsZero() {
            expiresAtProto = timestamppb.New(claims.ExpiresAt.Time)
        }
    }

	return &pb.ValidateTokenResponse{
		IsValid:   true,
		UserId:    claims.UserID,
		Username:  claims.Username,
		Email:     claims.Email,
		Roles:     claims.Roles,
		ExpiresAt: expiresAtProto,
	}, nil
}

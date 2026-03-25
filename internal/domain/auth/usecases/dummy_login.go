package usecases

import (
	"context"
	"time"

	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/common"
	"github.com/avito-internships/test-backend-1-EmotionlessDev/internal/domain/auth/dto"
	"github.com/golang-jwt/jwt/v5"
)

type DummyLogin struct {
	jwtSecret string
}

func NewDummyLogin(jwtSecret string) *DummyLogin {
	return &DummyLogin{jwtSecret: jwtSecret}
}

func (uc *DummyLogin) Execute(ctx context.Context, role string) (dto.TokenResponse, error) {
	var fixedUserID string

	switch role {
	case "admin":
		fixedUserID = "11111111-1111-1111-1111-111111111111"
	case "user":
		fixedUserID = "22222222-2222-2222-2222-222222222222"
	default:
		return dto.TokenResponse{}, common.ErrInvalidRole
	}

	claims := jwt.MapClaims{
		"user_id": fixedUserID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
		"iat":     time.Now().Unix(),
	}

	tokenJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := tokenJwt.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return dto.TokenResponse{}, err
	}

	return dto.TokenResponse{
		Token:     signedToken,
		UserID:    fixedUserID,
		Role:      role,
		CreatedAt: time.Now().UTC(),
	}, nil
}

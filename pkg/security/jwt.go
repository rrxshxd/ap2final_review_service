package security

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
	"strings"
	"time"
)

type JWTProvider struct {
	secretKey       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type Claims struct {
	UserID *string
	Role   *string
}

func NewJWTProvider(
	secretKey string,
	accessTokenTTL, refreshTokenTTL time.Duration,
) *JWTProvider {
	return &JWTProvider{
		secretKey:       secretKey,
		AccessTokenTTL:  accessTokenTTL,
		RefreshTokenTTL: refreshTokenTTL,
	}
}

func (p *JWTProvider) GenerateAccessToken(userID string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(p.AccessTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(p.secretKey))
}

func (p *JWTProvider) GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(p.RefreshTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(p.secretKey))
}

func (p *JWTProvider) VerifyAndParseClaims(tokenStr string) (Claims, error) {
	jwtClaims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenStr, jwtClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(p.secretKey), nil
	})

	if err != nil {
		return Claims{}, err
	}

	claims := Claims{}

	userID, ok := jwtClaims["user_id"].(string)
	if ok {
		claims.UserID = &userID
	}

	role, ok := jwtClaims["role"].(string)
	if ok {
		claims.Role = &role
	}

	return claims, nil
}

func TokenFromCtx(ctx context.Context) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}

	authHeader := md["authorization"]
	if len(authHeader) == 0 {
		return "", false
	}

	tokenStr := strings.TrimPrefix(authHeader[0], "Bearer ")

	return tokenStr, true
}

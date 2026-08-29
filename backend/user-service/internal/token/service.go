package token

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	accessTokenType  = "access"
	refreshTokenSize = 32
)

type Service struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

type AccessToken struct {
	Value     string
	ExpiresAt time.Time
}

type RefreshToken struct {
	Value     string
	Hash      string
	ExpiresAt time.Time
}

type jwtHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

type accessClaims struct {
	Subject   string `json:"sub"`
	TokenType string `json:"typ"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

func NewService(secret string, accessTTL, refreshTTL time.Duration) (*Service, error) {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return nil, fmt.Errorf("jwt secret is required")
	}
	if accessTTL <= 0 {
		return nil, fmt.Errorf("access token ttl must be positive")
	}
	if refreshTTL <= 0 {
		return nil, fmt.Errorf("refresh token ttl must be positive")
	}

	return &Service{
		secret:     []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}, nil
}

func (s *Service) NewAccessToken(userID string, now time.Time) (*AccessToken, error) {
	expiresAt := now.Add(s.accessTTL)
	claims := accessClaims{
		Subject:   userID,
		TokenType: accessTokenType,
		IssuedAt:  now.Unix(),
		ExpiresAt: expiresAt.Unix(),
	}

	value, err := s.sign(claims)
	if err != nil {
		return nil, err
	}

	return &AccessToken{
		Value:     value,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *Service) NewRefreshToken(now time.Time) (*RefreshToken, error) {
	randomBytes := make([]byte, refreshTokenSize)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	value := base64.RawURLEncoding.EncodeToString(randomBytes)

	return &RefreshToken{
		Value:     value,
		Hash:      s.HashRefreshToken(value),
		ExpiresAt: now.Add(s.refreshTTL),
	}, nil
}

func (s *Service) HashRefreshToken(value string) string {
	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(value))

	return hex.EncodeToString(mac.Sum(nil))
}

func (s *Service) sign(claims accessClaims) (string, error) {
	headerBytes, err := json.Marshal(jwtHeader{
		Algorithm: "HS256",
		Type:      "JWT",
	})
	if err != nil {
		return "", fmt.Errorf("marshal jwt header: %w", err)
	}

	claimsBytes, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("marshal jwt claims: %w", err)
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerBytes)
	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsBytes)
	signingInput := encodedHeader + "." + encodedClaims

	mac := hmac.New(sha256.New, s.secret)
	mac.Write([]byte(signingInput))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return signingInput + "." + signature, nil
}

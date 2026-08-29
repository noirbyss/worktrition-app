package authn

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var errInvalidAccessToken = errors.New("invalid access token")

type jwtHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

type accessClaims struct {
	Subject   string `json:"sub"`
	TokenType string `json:"typ"`
	ExpiresAt int64  `json:"exp"`
}

func verifyAccessToken(token string, secret []byte, now time.Time) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", errInvalidAccessToken
	}

	var header jwtHeader
	if err := decodeJWTPart(parts[0], &header); err != nil {
		return "", errInvalidAccessToken
	}
	if header.Algorithm != "HS256" {
		return "", errInvalidAccessToken
	}

	signingInput := parts[0] + "." + parts[1]
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(signingInput))
	expectedSignature := mac.Sum(nil)

	actualSignature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return "", errInvalidAccessToken
	}
	if !hmac.Equal(actualSignature, expectedSignature) {
		return "", errInvalidAccessToken
	}

	var claims accessClaims
	if err := decodeJWTPart(parts[1], &claims); err != nil {
		return "", errInvalidAccessToken
	}
	if claims.Subject == "" {
		return "", errInvalidAccessToken
	}
	if claims.TokenType != "" && claims.TokenType != "access" {
		return "", errInvalidAccessToken
	}
	if claims.ExpiresAt <= now.Unix() {
		return "", fmt.Errorf("access token expired")
	}

	return claims.Subject, nil
}

func decodeJWTPart(part string, target any) error {
	data, err := base64.RawURLEncoding.DecodeString(part)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, target)
}

package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Header represents standard JWT header
type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

// MapClaims represents custom JWT claims map
type MapClaims map[string]any

// JWT manager for signing and verifying tokens using HMAC-SHA256
type JWT struct {
	secret []byte
}

// New creates a new JWT helper with given secret key
func New(secret string) *JWT {
	return &JWT{secret: []byte(secret)}
}

// Sign creates a signed JWT token string with given claims and expiration duration
func (j *JWT) Sign(claims MapClaims, duration time.Duration) (string, error) {
	header := Header{
		Alg: "HS256",
		Typ: "JWT",
	}

	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("failed to marshal header: %w", err)
	}

	if claims == nil {
		claims = make(MapClaims)
	}
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["iat"] = time.Now().Unix()

	claimsBytes, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("failed to marshal claims: %w", err)
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerBytes)
	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsBytes)

	signingString := encodedHeader + "." + encodedClaims
	signature := j.signString(signingString)

	return signingString + "." + signature, nil
}

// Verify parses and verifies a JWT token string, returning claims if valid
func (j *JWT) Verify(tokenStr string) (MapClaims, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	signingString := parts[0] + "." + parts[1]
	expectedSig := j.signString(signingString)

	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		return nil, fmt.Errorf("invalid token signature")
	}

	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode claims: %w", err)
	}

	var claims MapClaims
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	if expVal, ok := claims["exp"]; ok {
		var expUnix int64
		switch v := expVal.(type) {
		case float64:
			expUnix = int64(v)
		case int64:
			expUnix = v
		}

		if time.Now().Unix() > expUnix {
			return nil, fmt.Errorf("token has expired")
		}
	}

	return claims, nil
}

func (j *JWT) signString(str string) string {
	mac := hmac.New(sha256.New, j.secret)
	mac.Write([]byte(str))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

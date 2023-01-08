package oauth2

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
)

const minSecretLength = 1

func TestNewHS256(t *testing.T) {
	// Test 1: Check that the function returns an error if an empty string is passed as the secret key
	_, err := newHS256("", time.Hour, "user1")
	if err == nil {
		t.Errorf("Expected an error when secret is an empty string")
	}

	// Test 2: Check that the function returns an error if a zero value is passed as the token expiration
	_, err = newHS256("secret", 0, "user1")
	if err == nil {
		t.Errorf("Expected an error when expiry is zero")
	}

	// Test 3: Check that the function returns an error if an empty string is passed as the user ID
	_, err = newHS256("secret", time.Hour, "")
	if err == nil {
		t.Errorf("Expected an error when user ID is an empty string")
	}

	// Test 4: Check that the function returns a token in the correct format
	token, err := newHS256("secret", time.Hour, "user1")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !strings.HasPrefix(token, "ey") {
		t.Errorf("Expected token to have prefix 'ey', but got %s", token)
	}

	// Test 5: Check that the function returns an error if the token cannot be signed
	secret := strings.Repeat("a", minSecretLength-1)
	_, err = newHS256(secret, time.Hour, "user1")
	if err == nil {
		t.Errorf("Expected an error when secret is shorter than %d characters", minSecretLength)
	}
}

func TestParseHS256(t *testing.T) {
	// Test 1: Check that the function returns an error if an empty string is passed as the token
	_, err := parseHS256("", "secret")
	if err == nil {
		t.Errorf("Expected an error when token is an empty string")
	}

	// Test 2: Check that the function returns an error if an empty string is passed as the secret key
	_, err = parseHS256("token", "")
	if err == nil {
		t.Errorf("Expected an error when secret is an empty string")
	}

	// Test 3: Check that the function returns an error if the token is invalid
	_, err = parseHS256("invalid_token", "secret")
	if err == nil {
		t.Errorf("Expected an error when token is invalid")
	}

	// Test 4: Check that the function returns an error if the token has expired
	secret := "secret"
	validToken, _ := newHS256(secret, -time.Hour, "user1")
	_, err = parseHS256(validToken, secret)
	if err == nil {
		t.Errorf("Expected an error when token has expired")
	}

	// Test 5: Check that the function returns an error if the secret key is incorrect
	validToken, _ = newHS256(secret, time.Hour, "user1")
	_, err = parseHS256(validToken, "incorrect_secret")
	if err == nil {
		t.Errorf("Expected an error when secret is incorrect")
	}

	// Test 6: Check that the function returns valid claims in the jwt.MapClaims format if a valid token is passed
	validToken, _ = newHS256(secret, time.Hour, "user1")

	var claims jwt.Claims
	claims, err = parseHS256(validToken, secret)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check that claims are returned in the jwt.MapClaims format
	if _, ok := claims.(jwt.MapClaims); !ok {
		t.Errorf("Expected claims to be of type jwt.MapClaims, but got %T", claims)
	}

	// Check that the correct claims are returned
	exp := claims.(jwt.MapClaims)["exp"].(float64)
	if int64(exp) != time.Now().Add(time.Hour).Unix() {
		t.Errorf("Expected exp claim to be %d, but got %f", time.Now().Add(time.Hour).Unix(), exp)
	}

	iat := claims.(jwt.MapClaims)["iat"].(float64)
	if int64(iat) != time.Now().Unix() {
		t.Errorf("Expected iat claim to be %d, but got %f", time.Now().Unix(), iat)
	}

	nbf := claims.(jwt.MapClaims)["nbf"].(float64)
	if int64(nbf) != time.Now().Unix() {
		t.Errorf("Expected nbf claim to be %d, but got %f", time.Now().Unix(), nbf)
	}

	userID := claims.(jwt.MapClaims)["user_id"].(string)
	if userID != "user1" {
		t.Errorf("Expected user_id claim to be %s, but got %s", "user1", userID)
	}
}

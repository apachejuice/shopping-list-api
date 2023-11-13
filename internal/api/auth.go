package api

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/palantir/stacktrace"
)

// Authentication with keycloak
type Authenticator struct {
	// the below variables are fetched from configuration
	KcUrl            string `json:"kc_url"`
	KcRealm          string `json:"kc_realm"`
	KcRealmPublicKey string `json:"kc_realm_public_key"`
	KcClientId       string `json:"kc_client_id"`
	KcClientSecret   string `json:"kc_client_secret"`
}

func (a *Authenticator) CheckToken(userToken string) (bool, *ApiError) {
	isValid, err := a.validateSignedToken(userToken)
	if err != nil {
		return false, NewApiError(stacktrace.Propagate(err, "Failed to validate token"), false)
	}

	return isValid, nil
}

// (isValid, error)
func (a *Authenticator) validateSignedToken(tokenString string) (bool, *ApiError) {
	publicKey, err := parseKeycloakRSAPublicKey(a.KcRealmPublicKey)
	if err != nil {
		return false, NewApiError(stacktrace.Propagate(err, "Failed to parse keycloak realm public key"), false)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, NewApiError(stacktrace.NewError("unexpected signing method: %v", token.Header["alg"]), false)
		}

		return publicKey, nil
	})

	if err != nil {
		return false, NewApiError(stacktrace.Propagate(err, "Error parsing or validating token"), false)
	}

	if !token.Valid {
		return false, NewApiError(stacktrace.NewError("Token is invalid"), false)
	}

	claims := token.Claims.(jwt.MapClaims)
	exp := int64(claims["exp"].(float64))
	if time.Unix(exp, 0).After(time.Now()) {
		return true, nil
	}

	return false, NewApiError(stacktrace.NewError("Token is expired"), false)
}

func parseKeycloakRSAPublicKey(base64Encoded string) (*rsa.PublicKey, error) {
	buf, err := base64.StdEncoding.DecodeString(base64Encoded)
	if err != nil {
		return nil, err
	}
	parsedKey, err := x509.ParsePKIXPublicKey(buf)
	if err != nil {
		return nil, err
	}
	publicKey, ok := parsedKey.(*rsa.PublicKey)
	if ok {
		return publicKey, nil
	}
	return nil, fmt.Errorf("unexpected key type %T", publicKey)
}

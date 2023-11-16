package api

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"apachejuice.dev/apachejuice/shopping-list-api/internal/logging"
	"apachejuice.dev/apachejuice/shopping-list-api/internal/repo"
	"github.com/Nerzal/gocloak/v13"
	"github.com/golang-jwt/jwt/v4"
	"github.com/palantir/stacktrace"
)

type AuthenticatorConfig struct {
	// the below variables are fetched from configuration
	KcUrl            string `json:"kc_url"`
	KcRealm          string `json:"kc_realm"`
	KcRealmPublicKey string `json:"kc_realm_public_key"`
	KcClientId       string `json:"kc_client_id"`
	KcClientSecret   string `json:"kc_client_secret"`
}

// Authentication with keycloak
type Authenticator struct {
	config     AuthenticatorConfig
	kcInstance *gocloak.GoCloak
}

func NewAuthenticator(config AuthenticatorConfig) Authenticator {
	a := Authenticator{config: config}
	a.kcInstance = gocloak.NewClient(a.config.KcUrl)
	go a.logSanityCheck()

	return a
}

func (a *Authenticator) logSanityCheck() {
	path, _ := url.JoinPath(a.config.KcUrl, "realms", a.config.KcRealm, ".well-known/openid-configuration")
	resp, err := a.kcInstance.RestyClient().R().Get(path)

	if err != nil {
		logging.Error(err, "-")
		log.Fatal(err)
	}

	var rspSchema struct {
		Issuer string `json:"issuer"`
	}

	err = json.Unmarshal(resp.Body(), &rspSchema)
	if err != nil {
		logging.Error(err, "-")
		log.Fatal(err)
	}

	logging.Info("Keycloak sanity check done, realm address %s", rspSchema.Issuer)
}

func (a *Authenticator) CheckToken(ctx context.Context, userToken string) (bool, string, *ApiError) {
	isValid, userId, err := a.validateSignedToken(ctx, userToken)
	if err != nil {
		return false, userId, NewApiError(stacktrace.Propagate(err, "Failed to validate token"), false)
	}

	return isValid, userId, nil
}

// (isValid, userId, error)
func (a *Authenticator) validateSignedToken(ctx context.Context, tokenString string) (bool, string, *ApiError) {
	publicKey, err := parseKeycloakRSAPublicKey(a.config.KcRealmPublicKey)
	if err != nil {
		return false, "", NewApiError(stacktrace.Propagate(err, "Failed to parse keycloak realm public key"), false)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, NewApiError(stacktrace.NewError("unexpected signing method: %v", token.Header["alg"]), false)
		}

		return publicKey, nil
	})

	if err != nil {
		if err.Error() == "Token is expired" {
			logging.Info("Short-circuiting API access, token expired")
			return false, "", NewApiError(stacktrace.NewError("Token is expired"), false)
		}

		return false, "", NewApiError(stacktrace.Propagate(err, "Error parsing or validating token"), false)
	}

	if !token.Valid {
		return false, "", NewApiError(stacktrace.NewError("Token is invalid"), false)
	}

	claims := token.Claims.(jwt.MapClaims)
	sub := claims["sub"].(string)

	// OK, token not expired and is correctly formatted. Check with keycloak:
	userinfo, err := a.kcInstance.GetUserInfo(ctx, tokenString, a.config.KcRealm)
	if err != nil {
		if aerr, ok := err.(gocloak.APIError); ok && aerr.Code == http.StatusUnauthorized {
			// User does not exist at all
			return false, sub, NewApiError(stacktrace.Propagate(err, "User does not exist"), false)
		}

		// Other error; not necessarily our fault but indicates some fault with keycloak, mark as 500
		return false, sub, NewApiError(stacktrace.Propagate(err, "Unable to get user info from keycloak"), true)
	}

	ok, err := repo.HasUserWithId(ctx, *userinfo.Sub)
	if !ok && err == nil {
		// User does not exist; create it
		err = repo.CreateUser(ctx, *userinfo.Sub, *userinfo.PreferredUsername, time.Now().UTC())
		if err != nil {
			return false, sub, NewApiError(stacktrace.Propagate(err, "Failed to create new user"), true)
		}
		return true, sub, nil
	} else if ok {
		return true, sub, nil
	}

	return false, sub, NewApiError(stacktrace.NewError("Unknown error"), true)
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

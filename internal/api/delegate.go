package api

import (
	"context"
	"net/http"

	"apachejuice.dev/apachejuice/shopping-list-api/internal/apispec"
	"apachejuice.dev/apachejuice/shopping-list-api/internal/repo"
	"apachejuice.dev/apachejuice/shopping-list-api/internal/translate"
	"github.com/palantir/stacktrace"
)

// Recieves api calls that are readily authenticated. Handles database calls and translation
type apiDelegate struct {
}

func (d apiDelegate) getMe(ctx context.Context, userId string) (*apispec.User, *ApiError) {
	dbUser, err := repo.GetUserWithId(ctx, userId)
	if err != nil {
		return nil, NewApiError(stacktrace.Propagate(err, "Unable to retrieve user from database"), true)
	}

	// No record found
	if dbUser == nil {
		return nil, NewApiErrorWithCode(stacktrace.NewError("No user found"), http.StatusNotFound)
	}

	return translate.UserToJson(dbUser), nil
}

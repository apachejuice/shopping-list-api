package api

import (
	"context"
	"net/http"

	"apachejuice.dev/apachejuice/shopping-list-api/internal/apispec"
	"apachejuice.dev/apachejuice/shopping-list-api/internal/model"
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

func (d apiDelegate) getLists(ctx context.Context, userId string) ([]*apispec.ShoppingList, *ApiError) {
	lists, err := repo.GetListsForUser(ctx, userId)
	if err != nil {
		return nil, NewApiError(stacktrace.Propagate(err, "Unable to retrieve lists from database"), true)
	}

	responseLists := make([]*apispec.ShoppingList, 0)
	for _, list := range lists {
		responseLists = append(responseLists, translate.ListToJson(list))
	}

	return responseLists, nil
}

func (d apiDelegate) getListId(ctx context.Context, listId, userId string) (*apispec.ShoppingList, *ApiError) {
	list, err := repo.GetListById(ctx, listId)
	if err != nil {
		return nil, NewApiError(stacktrace.Propagate(err, "Unable to get list by id"), true)
	}

	if list.CreatorID != userId {
		return nil, NewApiErrorWithCode(stacktrace.NewError("List not owned by user"), http.StatusUnauthorized)
	}

	return translate.ListToJson(list), nil
}

func (d apiDelegate) createList(ctx context.Context, list *apispec.ShoppingList, userId string) (*apispec.ShoppingList, *ApiError) {
	dbList := translate.JsonToList(list, userId)
	result, err := repo.AddList(ctx, dbList)
	if err != nil {
		return nil, NewApiError(stacktrace.Propagate(err, "Unable to insert list to database"), true)
	}

	if list.Items != nil {
		dbItems := make(model.ShoppingListItemSlice, 0)
		for _, item := range *list.Items {
			dbItem := translate.JsonToItem(&item, dbList.ListID)
			dbItems = append(dbItems, dbItem)
		}

		err = repo.AddItems(ctx, dbItems, dbList)
		if err != nil {
			return nil, NewApiError(err, true)
		}
	}

	return translate.ListToJson(result), nil
}

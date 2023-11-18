package repo

import (
	"context"

	"apachejuice.dev/apachejuice/shopping-list-api/internal/model"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func GetListsForUser(ctx context.Context, userId string) (model.ShoppingListSlice, error) {
	lists, err := model.ShoppingLists(
		selectWithCreatorId(userId)...,
	).All(ctx, db)

	return lists, processErr(err)
}

func HasList(ctx context.Context, id string) (bool, error) {
	ok, err := model.ShoppingLists(
		selectWithListId(id)...,
	).Exists(ctx, db)

	return ok, processErr(err)
}

func GetListById(ctx context.Context, id string) (*model.ShoppingList, error) {
	list, err := model.ShoppingLists(
		selectWithListId(id)...,
	).One(ctx, db)

	return list, processErr(err)
}

func AddList(ctx context.Context, list *model.ShoppingList) (*model.ShoppingList, error) {
	list.ListID = uuid.New().String()
	return list, list.Insert(ctx, db, boil.Infer())
}

package repo

import (
	"context"

	"apachejuice.dev/apachejuice/shopping-list-api/internal/model"
	"github.com/google/uuid"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func GetItemsForList(ctx context.Context, listId string) (model.ShoppingListItemSlice, error) {
	items, err := model.ShoppingListItems(
		selectItemWithListId(listId)...,
	).All(ctx, db)

	return items, processErr(err)
}

func AddItems(ctx context.Context, items model.ShoppingListItemSlice, list *model.ShoppingList) error {
	for _, item := range items {
		item.ItemID = uuid.New().String()
		err := item.Insert(ctx, db, boil.Infer())
		if err != nil {
			return err
		}
	}

	return nil
}

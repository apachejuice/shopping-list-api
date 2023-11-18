package translate

import (
	"apachejuice.dev/apachejuice/shopping-list-api/internal/apispec"
	"apachejuice.dev/apachejuice/shopping-list-api/internal/model"
)

func ItemToJson(item *model.ShoppingListItem) *apispec.ShoppingListItem {
	return &apispec.ShoppingListItem{
		Name: item.ItemName,
	}
}

func JsonToItem(item *apispec.ShoppingListItem, listId string) *model.ShoppingListItem {
	return &model.ShoppingListItem{
		ItemID:   item.Id.String(),
		ItemName: item.Name,
		ListID:   listId,
	}
}

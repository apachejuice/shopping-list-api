package translate

import (
	"apachejuice.dev/apachejuice/shopping-list-api/internal/apispec"
	"apachejuice.dev/apachejuice/shopping-list-api/internal/model"
	"github.com/google/uuid"
)

func ListToJson(list *model.ShoppingList) *apispec.ShoppingList {
	id := uuid.MustParse(list.ListID)
	itemCount := list.ItemCount

	return &apispec.ShoppingList{
		ItemCount: &itemCount,
		Id:        &id,
		Name:      list.ListName,
	}
}

func JsonToList(list *apispec.ShoppingList, userId string) *model.ShoppingList {
	itemCount := 0
	if list.Items != nil {
		itemCount = len(*list.Items)
	}

	return &model.ShoppingList{
		ItemCount: itemCount,
		ListName:  list.Name,
		CreatorID: userId,
	}
}

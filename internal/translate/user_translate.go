package translate

import (
	"apachejuice.dev/apachejuice/shopping-list-api/internal/apispec"
	"apachejuice.dev/apachejuice/shopping-list-api/internal/model"
	"github.com/google/uuid"
)

func UserToJson(user *model.User) *apispec.User {
	return &apispec.User{
		Name: user.Username,
		Id:   uuid.MustParse(user.UserID),
	}
}

func JsonToUser(user *apispec.User) *model.User {
	return &model.User{
		UserID:   user.Id.String(),
		Username: user.Name,
	}
}

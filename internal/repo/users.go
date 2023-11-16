package repo

import (
	"context"
	"time"

	"apachejuice.dev/apachejuice/shopping-list-api/internal/model"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func HasUserWithId(ctx context.Context, userId string) (bool, error) {
	result, err := model.Users(
		selectWithUserId(userId)...,
	).Exists(ctx, db)

	return result, processErr(err)
}

func GetUserWithId(ctx context.Context, userId string) (*model.User, error) {
	user, err := model.Users(
		selectWithUserId(userId)...,
	).One(ctx, db)

	return user, processErr(err)
}

func CreateUser(ctx context.Context, userId string, username string, createdAt time.Time) error {
	user := &model.User{
		UserID:    userId,
		Username:  username,
		CreatedAt: createdAt,
	}

	return user.Insert(ctx, db, boil.Infer())
}

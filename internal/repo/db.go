package repo

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"apachejuice.dev/apachejuice/shopping-list-api/internal/logging"
	"apachejuice.dev/apachejuice/shopping-list-api/internal/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var (
	db *sql.DB
)

func Connect() {
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_PASS")
	dbname := os.Getenv("MYSQL_DBNAME")

	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?parseTime=true", user, pass, dbname))
	if err != nil {
		log.Fatal(err)
	}

	logging.Info("Connected to database")
}

func processErr(err error) error {
	if err == sql.ErrNoRows {
		return nil
	}

	return err
}

func selectWithUserId(userId string) []qm.QueryMod {
	// DO NOT USE ? AS A PLACEHOLDER: %s for table and column names (MySQL expects backticks or nothing) and %q for string values
	return []qm.QueryMod{qm.Select("*"), model.UserWhere.UserID.EQ(userId)}
}

func selectWithCreatorId(creatorId string) []qm.QueryMod {
	return []qm.QueryMod{qm.Select("*"), model.ShoppingListWhere.CreatorID.EQ(creatorId)}
}

func selectWithListId(listId string) []qm.QueryMod {
	return []qm.QueryMod{qm.Select("*"), model.ShoppingListWhere.ListID.EQ(listId), qm.Limit(1)} // only one row can ever be returned
}

func selectItemWithListId(listId string) []qm.QueryMod {
	return []qm.QueryMod{qm.Select("*"), model.ShoppingListItemWhere.ListID.EQ(listId)}
}

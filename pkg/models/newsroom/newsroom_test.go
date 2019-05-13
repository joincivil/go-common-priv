package newsroom_test

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	. "github.com/joincivil/go-common-priv/pkg/models/article"
	. "github.com/joincivil/go-common-priv/pkg/models/newsroom"
	"github.com/joincivil/go-common-priv/pkg/models/testutils"
	"testing"
)

const (
	postgresstr = "host=localhost port=5432 user=docker password=docker dbname=civil_crawler sslmode=disable"
)

func TestCreateNewsroom(t *testing.T) {
	db, err := gorm.Open("postgres", postgresstr)
	if err != nil {
		fmt.Println(err)
		t.Fatal("could not connect to db")
	}
	defer db.Close()

	cleaner := testutils.DeleteCreatedEntities(db)
	defer cleaner()

	testutils.MigrateModels(db)

	newsroom := Newsroom{
		Name:    "Newsroom1",
		Address: "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
	}

	if err := db.Create(&newsroom).Error; err != nil {
		t.Errorf("should have created a newsroom")
	}

	newsroomb := Newsroom{
		Name:    "Newsroom2",
		Address: "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
	}

	if err := db.Create(&newsroomb).Error; err == nil {
		t.Errorf("should have thrown an error because of the duplicate address")
	}
}

func TestAddArticle(t *testing.T) {
	db, err := gorm.Open("postgres", postgresstr)
	if err != nil {
		fmt.Println(err)
		t.Fatal("could not connect to db")
	}
	defer db.Close()

	cleaner := testutils.DeleteCreatedEntities(db)
	defer cleaner()

	testutils.MigrateModels(db)

	newsroom := Newsroom{
		Name:    "Newsroom1",
		Address: "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
	}

	if err := db.Create(&newsroom).Error; err != nil {
		t.Errorf("failed to make newsroom")
	}

	var count1 int
	db.Model(&Article{}).Where("newsroom_address = ?", "0x8c722B8AC728aDd7780a66017e8daDBa530EE261").Count(&count1)

	articleMeta := json.RawMessage(`{"title": "Worlds Greatest Article"}`)

	result := db.Model(&newsroom).Association("Articles").Append(Article{
		ArticleMetadata: postgres.Jsonb{articleMeta},
	})

	if result.Error != nil {
		t.Errorf("should not fail")
	}

	var count2 int
	db.Model(&Article{}).Where("newsroom_address = ?", "0x8c722B8AC728aDd7780a66017e8daDBa530EE261").Count(&count2)

	if count1 >= count2 {
		t.Errorf("expected another article with the newsrooms address to be added")
	}
}

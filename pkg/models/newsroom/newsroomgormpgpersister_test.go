package newsroom_test

import (
	"fmt"
	"github.com/joincivil/go-common-priv/pkg/models/article"
	"github.com/joincivil/go-common-priv/pkg/models/newsroom"
	"github.com/joincivil/go-common-priv/pkg/models/testutils"
	"testing"
)

const (
	postgressHost     = "localhost"
	postgressPort     = 5432
	postgressUser     = "docker"
	postgressPassword = "docker"
	dbname            = "civil_crawler"
)

func TestCreateNewsroom(t *testing.T) {
	pg, err := newsroom.NewGormPGPersister(postgressHost, postgressPort, postgressUser, postgressPassword, dbname)

	defer pg.DB.Close()

	cleaner := testutils.DeleteCreatedEntities(pg.DB)
	defer cleaner()

	testutils.MigrateModels(pg.DB)

	if err != nil {
		fmt.Println(err)
		t.Errorf("threw an error making the persister")
	}

	newsrooma := &newsroom.Newsroom{
		Name:    "Newsroom1",
		Address: "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
	}

	if newsrooma.ID != 0 {
		t.Errorf("shouldn't have an id yet")
	}

	if err := pg.CreateNewsroom(newsrooma); err != nil {
		t.Errorf("should have created a newsroom")
	}

	if newsrooma.ID == 0 {
		t.Errorf("should have an ID now")
	}

	newsroomb := &newsroom.Newsroom{
		Name:    "Newsroom2",
		Address: "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
	}

	if err := pg.CreateNewsroom(newsroomb); err == nil {
		t.Errorf("should have thrown an error because of the duplicate address")
	}
}

func TestUpdateNewsroom(t *testing.T) {
	pg, err := newsroom.NewGormPGPersister(postgressHost, postgressPort, postgressUser, postgressPassword, dbname)

	defer pg.DB.Close()

	cleaner := testutils.DeleteCreatedEntities(pg.DB)
	defer cleaner()

	testutils.MigrateModels(pg.DB)

	if err != nil {
		fmt.Println(err)
		t.Errorf("threw an error making the persister")
	}

	newsrooma := &newsroom.Newsroom{
		Name:    "Newsroom1",
		Address: "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
	}

	if err := pg.CreateNewsroom(newsrooma); err != nil {
		t.Errorf("should have created a newsroom")
	}

}

func TestAddArticle(t *testing.T) {
	pg, err := newsroom.NewGormPGPersister(postgressHost, postgressPort, postgressUser, postgressPassword, dbname)

	defer pg.DB.Close()

	cleaner := testutils.DeleteCreatedEntities(pg.DB)
	defer cleaner()

	testutils.MigrateModels(pg.DB)

	if err != nil {
		fmt.Println(err)
		t.Errorf("threw an error making the persister")
	}

	newsrooma := &newsroom.Newsroom{
		Name:    "Newsroom1",
		Address: "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
	}

	if err := pg.CreateNewsroom(newsrooma); err != nil {
		t.Errorf("should have created a newsroom")
	}

	articleMeta := &article.ArticleMetadata{
		Title:        "new stufff",
		CanonicalURL: "https://newstuff.bz/newarticle",
	}

	narticle := &article.Article{
		ArticleMetadata: *articleMeta,
		NewsroomAddress: "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
	}

	if err := pg.AddArticle(newsrooma.ID, narticle); err != nil {
		t.Errorf("failed to add article")
	}

	articles, articleErr := pg.GetArticlesForNewsroom(newsrooma.ID)

	if articleErr != nil {
		fmt.Println(articleErr)
		t.Errorf("couldnt get articles")
	}

	if len(articles) != 1 {
		t.Errorf("article wasnt added")
	}
}

func TestNewsroomByID(t *testing.T) {
	pg, err := newsroom.NewGormPGPersister(postgressHost, postgressPort, postgressUser, postgressPassword, dbname)

	defer pg.DB.Close()

	cleaner := testutils.DeleteCreatedEntities(pg.DB)
	defer cleaner()

	testutils.MigrateModels(pg.DB)

	if err != nil {
		fmt.Println(err)
		t.Errorf("threw an error making the persister")
	}

	newsrooma := &newsroom.Newsroom{
		Name:    "Newsroom1",
		Address: "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
	}

	if err := pg.CreateNewsroom(newsrooma); err != nil {
		t.Errorf("should have created a newsroom")
	}

	foundNewsroom, lookuperr := pg.NewsroomByID(newsrooma.ID)

	if lookuperr != nil {
		fmt.Println(err)
		t.Errorf("threw an error looking up the newsroom")
	}

	if foundNewsroom.Address != "0x8c722B8AC728aDd7780a66017e8daDBa530EE261" {
		t.Errorf("newsroom data is incorrect")
	}

	if foundNewsroom.ID != newsrooma.ID {
		t.Errorf("isn't the same newsroom")
	}

}

package article_test

import (
	"fmt"
	"testing"
	"time"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/joincivil/go-common-priv/pkg/models/article"
	"github.com/joincivil/go-common-priv/pkg/models/testutils"
	carticle "github.com/joincivil/go-common/pkg/article"

	gormutils "github.com/joincivil/go-common-priv/pkg/utils/gorm"
)

func testFunc(persister article.Persister) {
}

func TestGormInterface(t *testing.T) {
	// Ensure the GORM persister implements the Persister interface
	creds := testutils.GetTestDBConnection()
	pg, _ := article.NewGormPGPersister(creds.Host, creds.Port, creds.User, creds.Password, creds.Dbname)
	testFunc(pg)
}

func TestCreateArticle(t *testing.T) {
	creds := testutils.GetTestDBConnection()
	// Test out the WithDB GormPGPersister constructor
	db, err := gormutils.NewGormPGConnection(creds.Host, creds.Port, creds.User,
		creds.Password, creds.Dbname, 2, 2, 10*time.Second)
	if err != nil {
		fmt.Println(err)
		t.Errorf("threw an error creating the db conn")
	}

	pg, err := article.NewGormPGPersisterWithDB(db)

	if err != nil {
		fmt.Println(err)
		t.Errorf("threw an error making the persister")
	}

	defer pg.DB.Close()

	cleaner := testutils.DeleteCreatedEntities(pg.DB)
	defer cleaner()

	testutils.MigrateModels(pg.DB) // nolint: errcheck

	articleMeta := &carticle.Metadata{
		Title:        "new stufff",
		CanonicalURL: "https://newstuff.bz/newarticle",
	}

	narticle := &carticle.Article{
		ArticleMetadata: *articleMeta,
		NewsroomAddress: "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
	}

	if narticle.ID != 0 {
		t.Errorf("new article shouldnt have an id yet")
	}

	pg.CreateArticle(narticle) //nolint:errcheck

	if narticle.ID == 0 {
		t.Errorf("an id should be assigned to the narticle after save")
	}
}

func TestArticleByID(t *testing.T) {
	creds := testutils.GetTestDBConnection()
	pg, err := article.NewGormPGPersister(creds.Host, creds.Port, creds.User, creds.Password, creds.Dbname)

	if err != nil {
		fmt.Println(err)
		t.Errorf("threw an error making the persister")
	}

	defer pg.DB.Close()

	testutils.MigrateModels(pg.DB) // nolint: errcheck

	cleaner := testutils.DeleteCreatedEntities(pg.DB)
	defer cleaner()

	articleMeta := &carticle.Metadata{
		Title:        "new stufff",
		CanonicalURL: "https://newstuff.bz/newarticle",
	}

	narticle := &carticle.Article{
		ArticleMetadata: *articleMeta,
		NewsroomAddress: "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
	}

	pg.CreateArticle(narticle) // nolint: errcheck

	foundarticle, lookuperr := pg.ArticleByID(narticle.ID)

	if lookuperr != nil {
		fmt.Println(err)
		t.Errorf("threw an error looking up the article")
	}

	if foundarticle.ArticleMetadata.Title != articleMeta.Title {
		t.Errorf("article metadata is incorrect")
	}
}

func TestUpdateArticle(t *testing.T) {
	creds := testutils.GetTestDBConnection()
	pg, err := article.NewGormPGPersister(creds.Host, creds.Port, creds.User, creds.Password, creds.Dbname)

	if err != nil {
		fmt.Println(err)
		t.Errorf("threw an error making the persister")
	}

	testutils.MigrateModels(pg.DB) // nolint: errcheck

	defer pg.DB.Close()

	cleaner := testutils.DeleteCreatedEntities(pg.DB)
	defer cleaner()

	articleMeta := &carticle.Metadata{
		Title:        "new stufff",
		CanonicalURL: "https://newstuff.bz/newarticle",
	}

	narticle := &carticle.Article{
		ArticleMetadata: *articleMeta,
		NewsroomAddress: "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
	}

	pg.CreateArticle(narticle) // nolint: errcheck

	foundarticle, _ := pg.ArticleByID(narticle.ID)

	blockData := testutils.MakeFakeReceipt()

	foundarticle.BlockData = blockData

	if err := pg.UpdateArticle(foundarticle); err != nil {
		fmt.Println(err)
		t.Errorf("error saving article")
	}

	refetchArticle, _ := pg.ArticleByID(narticle.ID)
	if refetchArticle.BlockData.TxHash == (ethCommon.Hash{}) {
		t.Errorf("should have saved the new tx receipt")
	}
}

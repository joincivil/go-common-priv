package article_test

import (
	"fmt"
	"testing"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/joincivil/go-common-priv/pkg/models/article"
	"github.com/joincivil/go-common-priv/pkg/models/testutils"
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
	pg, err := article.NewGormPGPersister(creds.Host, creds.Port, creds.User, creds.Password, creds.Dbname)

	if err != nil {
		fmt.Println(err)
		t.Errorf("threw an error making the persister")
	}

	defer pg.DB.Close()

	cleaner := testutils.DeleteCreatedEntities(pg.DB)
	defer cleaner()

	testutils.MigrateModels(pg.DB) // nolint: errcheck

	articleMeta := &article.Metadata{
		Title:        "new stufff",
		CanonicalURL: "https://newstuff.bz/newarticle",
	}

	narticle := &article.Article{
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

	articleMeta := &article.Metadata{
		Title:        "new stufff",
		CanonicalURL: "https://newstuff.bz/newarticle",
	}

	narticle := &article.Article{
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

	articleMeta := &article.Metadata{
		Title:        "new stufff",
		CanonicalURL: "https://newstuff.bz/newarticle",
	}

	narticle := &article.Article{
		ArticleMetadata: *articleMeta,
		NewsroomAddress: "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
	}

	pg.CreateArticle(narticle) // nolint: errcheck

	foundarticle, _ := pg.ArticleByID(narticle.ID)

	blockData := ethTypes.Receipt{
		GasUsed: 7000000,
	}

	foundarticle.BlockData = blockData

	if err := pg.UpdateArticle(foundarticle); err != nil {
		fmt.Println(err)
		t.Errorf("error saving article")
	}
}

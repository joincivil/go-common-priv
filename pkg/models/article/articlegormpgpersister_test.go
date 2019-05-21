package article_test

import (
	"github.com/joincivil/go-common-priv/pkg/models/article"
	"github.com/joincivil/go-common-priv/pkg/models/testutils"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"fmt"
	"testing"
)

const (
	postgressHost = "localhost"
	postgressPort=5432
	postgressUser="docker"
	postgressPassword="docker"
	dbname="civil_crawler"
)

func TestCreateArticle(t *testing.T) {
	pg, err := article.NewArticleGormPGPersister(postgressHost, postgressPort, postgressUser, postgressPassword, dbname)

	defer pg.DB.Close()

	cleaner := testutils.DeleteCreatedEntities(pg.DB)
	defer cleaner()
	
	testutils.MigrateModels(pg.DB)

	if err != nil {
		fmt.Println(err)
		t.Errorf("threw an error making the persister")
	}

	articleMeta := &article.ArticleMetadata{
		Title: "new stufff",
		CanonicalURL: "https://newstuff.bz/newarticle",
	}

	narticle := &article.Article{
		ArticleMetadata: *articleMeta,
		NewsroomAddress: "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
	}

	if narticle.ID != 0 {
		t.Errorf("new article shouldnt have an id yet")
	}

	pg.CreateArticle(narticle)

	if narticle.ID == 0 {
		t.Errorf("an id should be assigned to the narticle after save")
	}
}

func TestArticleByID(t *testing.T) {
	pg, err := article.NewArticleGormPGPersister(postgressHost, postgressPort, postgressUser, postgressPassword, dbname)

	defer pg.DB.Close()

	cleaner := testutils.DeleteCreatedEntities(pg.DB)
	defer cleaner()
	
	testutils.MigrateModels(pg.DB)

	if err != nil {
		fmt.Println(err)
		t.Errorf("threw an error making the persister")
	}

	articleMeta := &article.ArticleMetadata{
		Title: "new stufff",
		CanonicalURL: "https://newstuff.bz/newarticle",
	}

	narticle := &article.Article{
		ArticleMetadata: *articleMeta,
		NewsroomAddress: "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
	}

	pg.CreateArticle(narticle)

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
	pg, err := article.NewArticleGormPGPersister(postgressHost, postgressPort, postgressUser, postgressPassword, dbname)

	defer pg.DB.Close()

	cleaner := testutils.DeleteCreatedEntities(pg.DB)
	defer cleaner()
	
	testutils.MigrateModels(pg.DB)

	if err != nil {
		fmt.Println(err)
		t.Errorf("threw an error making the persister")
	}

	articleMeta := &article.ArticleMetadata{
		Title: "new stufff",
		CanonicalURL: "https://newstuff.bz/newarticle",
	}

	narticle := &article.Article{
		ArticleMetadata: *articleMeta,
		NewsroomAddress: "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
	}

	pg.CreateArticle(narticle)

	foundarticle, _ := pg.ArticleByID(narticle.ID)

	blockData := ethTypes.Receipt {
		GasUsed: 7000000,
	}
	
	foundarticle.BlockData = blockData

	if err := pg.UpdateArticle(foundarticle); err != nil {
		fmt.Println(err)
		t.Errorf("error saving article")
	}
}

package newsroom_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/joincivil/go-common-priv/pkg/models/newsroom"
	"github.com/joincivil/go-common-priv/pkg/models/testutils"
	gormutils "github.com/joincivil/go-common-priv/pkg/utils/gorm"
	carticle "github.com/joincivil/go-common/pkg/article"
)

func testFunc(persister newsroom.Persister) {
}

func TestGormInterface(t *testing.T) {
	// Ensure the GORM persister implements the Persister interface
	creds := testutils.GetTestDBConnection()
	pg, _ := newsroom.NewGormPGPersister(creds.Host, creds.Port, creds.User, creds.Password, creds.Dbname)
	testFunc(pg)
}

func TestCreateNewsroom(t *testing.T) {
	creds := testutils.GetTestDBConnection()
	// Test out the WithDB GormPGPersister constructor
	db, err := gormutils.NewGormPGConnection(creds.Host, creds.Port, creds.User,
		creds.Password, creds.Dbname, 2, 2, 10*time.Second)
	if err != nil {
		fmt.Println(err)
		t.Errorf("threw an error creating the db conn")
	}

	pg, err := newsroom.NewGormPGPersisterWithDB(db)

	if err != nil {
		fmt.Println(err)
		t.Errorf("threw an error making the persister")
	}

	testutils.MigrateModels(pg.DB) // nolint: errcheck

	defer pg.DB.Close()

	cleaner := testutils.DeleteCreatedEntities(pg.DB)
	defer cleaner()

	meta := &newsroom.Meta{
		Index: true,
		Claim: false,
	}

	newsrooma := &newsroom.Newsroom{
		Name:    "Newsroom1",
		Address: "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
		Meta:    meta,
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
	creds := testutils.GetTestDBConnection()
	pg, err := newsroom.NewGormPGPersister(creds.Host, creds.Port, creds.User, creds.Password, creds.Dbname)

	if err != nil {
		fmt.Println(err)
		t.Errorf("threw an error making the persister")
	}

	testutils.MigrateModels(pg.DB) // nolint: errcheck

	defer pg.DB.Close()

	cleaner := testutils.DeleteCreatedEntities(pg.DB)
	defer cleaner()

	newsrooma := &newsroom.Newsroom{
		Name:    "Newsroom1",
		Address: "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
	}

	if err := pg.CreateNewsroom(newsrooma); err != nil {
		t.Errorf("should have created a newsroom")
	}

}

func TestAddArticle(t *testing.T) {
	creds := testutils.GetTestDBConnection()
	pg, err := newsroom.NewGormPGPersister(creds.Host, creds.Port, creds.User, creds.Password, creds.Dbname)

	if err != nil {
		fmt.Println(err)
		t.Errorf("threw an error making the persister")
	}

	testutils.MigrateModels(pg.DB) // nolint: errcheck

	defer pg.DB.Close()

	cleaner := testutils.DeleteCreatedEntities(pg.DB)
	defer cleaner()

	newsrooma := &newsroom.Newsroom{
		Name:    "Newsroom1",
		Address: "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
	}

	if err := pg.CreateNewsroom(newsrooma); err != nil {
		t.Errorf("should have created a newsroom")
	}

	articleMeta := &carticle.Metadata{
		Title:        "new stufff",
		CanonicalURL: "https://newstuff.bz/newarticle",
	}

	narticle := &carticle.Article{
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

func TestNewsrooms(t *testing.T) {
	creds := testutils.GetTestDBConnection()
	pg, err := newsroom.NewGormPGPersister(creds.Host, creds.Port, creds.User, creds.Password, creds.Dbname)

	if err != nil {
		fmt.Println(err)
		t.Errorf("threw an error making the persister")
	}

	testutils.MigrateModels(pg.DB) // nolint: errcheck

	defer pg.DB.Close()

	cleaner := testutils.DeleteCreatedEntities(pg.DB)
	defer cleaner()

	meta := &newsroom.Meta{
		Index: true,
		Claim: false,
	}

	newsrooma := &newsroom.Newsroom{
		Name:    "Newsroom1",
		Address: "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
		Meta:    meta,
	}

	if err1 := pg.CreateNewsroom(newsrooma); err1 != nil {
		t.Errorf("should have created a newsroom1")
	}

	newsrooma = &newsroom.Newsroom{
		Name:    "Newsroom2",
		Address: "0x9c722B8AF728aDd7780a66017e8daDBa530EE261",
	}

	if err1 := pg.CreateNewsroom(newsrooma); err1 != nil {
		t.Errorf("should have created a newsroom2")
	}

	newsrooma = &newsroom.Newsroom{
		Name:    "Newsroom3",
		Address: "0x9d822B8AF728aDd7780a66017e8daDBa530EE261",
	}

	if err1 := pg.CreateNewsroom(newsrooma); err1 != nil {
		t.Errorf("should have created a newsroom3")
	}

	newsrooms, err := pg.Newsrooms()
	if err != nil {
		t.Errorf("should have retrieved newsrooms: %v", err)
	}

	if len(newsrooms) != 3 {
		t.Errorf("should have retrieved 3 newsrooms: len: %v", len(newsrooms))
	}

	if newsrooms[0].Name != "Newsroom1" {
		t.Errorf("should have gotten Newsroom1")
	}
	if newsrooms[0].Address != "0x8c722B8AC728aDd7780a66017e8daDBa530EE261" {
		t.Errorf("should have gotten Newsroom1 address")
	}
	if newsrooms[0].Meta == nil {
		t.Errorf("should have gotten meta for Newsroom1 address")
	}
	if !newsrooms[0].Meta.Index {
		t.Errorf("should have gotten correct index meta value for Newsroom1")
	}
	if newsrooms[0].Meta.Claim {
		t.Errorf("should have gotten correct claim meta value for Newsroom1")
	}
	if newsrooms[1].Name != "Newsroom2" {
		t.Errorf("should have gotten Newsroom2")
	}
	if newsrooms[1].Address != "0x9c722b8Af728aDD7780A66017e8daDba530eE261" {
		t.Errorf("should have gotten Newsroom2 address")
	}
	if newsrooms[2].Name != "Newsroom3" {
		t.Errorf("should have gotten Newsroom3")
	}
	if newsrooms[2].Address != "0x9d822b8aF728AdD7780a66017E8dadBA530ee261" {
		t.Errorf("should have gotten Newsroom3 address")
	}

}

func TestNewsroomByID(t *testing.T) {
	creds := testutils.GetTestDBConnection()
	pg, err := newsroom.NewGormPGPersister(creds.Host, creds.Port, creds.User, creds.Password, creds.Dbname)

	if err != nil {
		fmt.Println(err)
		t.Errorf("threw an error making the persister")
	}

	testutils.MigrateModels(pg.DB) // nolint: errcheck

	defer pg.DB.Close()

	cleaner := testutils.DeleteCreatedEntities(pg.DB)
	defer cleaner()

	newsrooma := &newsroom.Newsroom{
		Name:    "Newsroom1",
		Address: "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
	}

	if err1 := pg.CreateNewsroom(newsrooma); err1 != nil {
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

func TestNewsroomByAddress(t *testing.T) {
	creds := testutils.GetTestDBConnection()
	pg, err := newsroom.NewGormPGPersister(creds.Host, creds.Port, creds.User, creds.Password, creds.Dbname)

	if err != nil {
		fmt.Println(err)
		t.Errorf("threw an error making the persister")
	}

	testutils.MigrateModels(pg.DB) // nolint: errcheck

	defer pg.DB.Close()

	cleaner := testutils.DeleteCreatedEntities(pg.DB)
	defer cleaner()

	addr := "0x8c722B8AC728aDd7780a66017e8daDBa530EE261"

	newsrooma := &newsroom.Newsroom{
		Name:    "Newsroom1",
		Address: strings.ToLower(addr), // Test normalization
	}

	if err1 := pg.CreateNewsroom(newsrooma); err1 != nil {
		t.Errorf("should have created a newsroom")
	}

	foundNewsroom, lookuperr := pg.NewsroomByAddress(addr)

	if lookuperr != nil {
		fmt.Println(err)
		t.Errorf("threw an error looking up the newsroom")
	}

	if foundNewsroom.Address != addr {
		t.Errorf("newsroom data is incorrect")
	}

	if foundNewsroom.ID != newsrooma.ID {
		t.Errorf("isn't the same newsroom")
	}

}

func TestGetLatestArticleForNewsroom(t *testing.T) {
	creds := testutils.GetTestDBConnection()
	pg, err := newsroom.NewGormPGPersister(creds.Host, creds.Port, creds.User, creds.Password, creds.Dbname)

	if err != nil {
		fmt.Println(err)
		t.Errorf("threw an error making the persister")
	}

	testutils.MigrateModels(pg.DB) // nolint: errcheck

	defer pg.DB.Close()

	cleaner := testutils.DeleteCreatedEntities(pg.DB)
	defer cleaner()

	newsrooma := &newsroom.Newsroom{
		Name:    "Newsroom1",
		Address: "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
	}

	if err1 := pg.CreateNewsroom(newsrooma); err1 != nil {
		t.Errorf("should have created a newsroom")
	}

	_, err = pg.GetLatestArticleForNewsroom(newsrooma.ID)
	if err == nil {
		t.Errorf("should have gotten error")
	}

	now := time.Now()

	articleMeta := &carticle.Metadata{
		Title:        "new stufff latest",
		CanonicalURL: "https://newstuff.bz/newarticle",
		RevisionDate: now.Add(30 * time.Second),
	}

	narticle := &carticle.Article{
		ArticleMetadata: *articleMeta,
		NewsroomAddress: "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
	}

	if err1 := pg.AddArticle(newsrooma.ID, narticle); err1 != nil {
		t.Errorf("failed to add latest article")
	}

	articleMeta = &carticle.Metadata{
		Title:        "new stufff old",
		CanonicalURL: "https://newstuff.bz/newarticle",
		RevisionDate: now,
	}

	narticle = &carticle.Article{
		ArticleMetadata: *articleMeta,
		NewsroomAddress: "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
	}

	if err1 := pg.AddArticle(newsrooma.ID, narticle); err1 != nil {
		t.Errorf("failed to add old article")
	}

	articleMeta = &carticle.Metadata{
		Title:        "new stufff mid",
		CanonicalURL: "https://newstuff.bz/newarticle",
		RevisionDate: now.Add(15 * time.Second),
	}

	narticle = &carticle.Article{
		ArticleMetadata: *articleMeta,
		NewsroomAddress: "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
	}

	if err1 := pg.AddArticle(newsrooma.ID, narticle); err1 != nil {
		t.Errorf("failed to add old article")
	}

	art, err := pg.GetLatestArticleForNewsroom(newsrooma.ID)
	if err != nil {
		t.Errorf("failed to get latest article: %v", err)
	}

	if art.ArticleMetadata.Title != "new stufff latest" {
		t.Errorf("failed to fetch the latest article: %v", art.ArticleMetadata.Title)
	}
}

func TestGetArticlesForNewsroomIndexedSinceDate(t *testing.T) {
	creds := testutils.GetTestDBConnection()
	pg, err := newsroom.NewGormPGPersister(creds.Host, creds.Port, creds.User, creds.Password, creds.Dbname)

	if err != nil {
		fmt.Println(err)
		t.Errorf("threw an error making the persister")
	}

	testutils.MigrateModels(pg.DB) // nolint: errcheck

	defer pg.DB.Close()

	cleaner := testutils.DeleteCreatedEntities(pg.DB)
	defer cleaner()

	newsrooma := &newsroom.Newsroom{
		Name:    "Newsroom1",
		Address: "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
	}

	if err1 := pg.CreateNewsroom(newsrooma); err1 != nil {
		t.Errorf("should have created a newsroom")
	}

	arts, err := pg.GetArticlesForNewsroomIndexedSinceDate(newsrooma.ID, time.Now())
	if err != nil {
		t.Errorf("shouldn't throw error on empty set %v", err)
	}
	if len(arts) != 0 {
		t.Errorf("shouldn't be any articles")
	}

	now := time.Now()

	articleMeta := &carticle.Metadata{
		Title:        "new stufff latest",
		CanonicalURL: "https://newstuff.bz/newarticle",
	}

	narticle := &carticle.Article{
		ArticleMetadata:  *articleMeta,
		NewsroomAddress:  "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
		IndexedTimestamp: now.Add(30 * time.Second),
	}

	if err1 := pg.AddArticle(newsrooma.ID, narticle); err1 != nil {
		t.Errorf("failed to add latest article")
	}

	articleMeta = &carticle.Metadata{
		Title:        "new stufff old",
		CanonicalURL: "https://newstuff.bz/newarticle",
	}

	narticle = &carticle.Article{
		ArticleMetadata:  *articleMeta,
		NewsroomAddress:  "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
		IndexedTimestamp: now.Add(40 * time.Second),
	}

	if err1 := pg.AddArticle(newsrooma.ID, narticle); err1 != nil {
		t.Errorf("failed to add old article")
	}

	articleMeta = &carticle.Metadata{
		Title:        "new stufff mid",
		CanonicalURL: "https://newstuff.bz/newarticle",
	}

	narticle = &carticle.Article{
		ArticleMetadata:  *articleMeta,
		NewsroomAddress:  "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
		IndexedTimestamp: now.Add(55),
	}

	if err1 := pg.AddArticle(newsrooma.ID, narticle); err1 != nil {
		t.Errorf("failed to add old article")
	}

	articleMeta = &carticle.Metadata{
		Title:        "new stufff mid",
		CanonicalURL: "https://newstuff.bz/newarticle",
	}

	narticle = &carticle.Article{
		ArticleMetadata:  *articleMeta,
		NewsroomAddress:  "0x8c722B8AC728aDd7780a66017e8daDBa530EE261",
		IndexedTimestamp: now.Add(-1 * time.Second),
	}

	if err1 := pg.AddArticle(newsrooma.ID, narticle); err1 != nil {
		t.Errorf("failed to add old article")
	}

	art, err := pg.GetArticlesForNewsroomIndexedSinceDate(newsrooma.ID, now)
	if err != nil {
		t.Errorf("failed to get latest article: %v", err)
	}

	if len(art) != 3 {
		t.Errorf("did not fetch all or only articles indexed after now")
	}
}

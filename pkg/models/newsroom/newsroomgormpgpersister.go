package newsroom

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/golang/glog"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"

	"github.com/joincivil/go-common-priv/pkg/models/article"
	carticle "github.com/joincivil/go-common/pkg/article"
)

var (
	// ErrNoArticles indicates that there were no articles found for the query
	ErrNoArticles = errors.New("no articles found")
)

const (
	// Could make this configurable later if needed
	maxOpenConns    = 50
	maxIdleConns    = 10
	connMaxLifetime = time.Second * 1800 // 30 mins
)

// Gorm is the newsroom schema
type Gorm struct {
	gorm.Model
	Name     string
	Address  string `gorm:"unique;not null"`
	Meta     postgres.Jsonb
	Articles []article.Gorm `gorm:"foreignkey:NewsroomAddress;association_foreignkey:Address"`
}

// TableName sets the name of the corresponding table in the db
func (Gorm) TableName() string {
	return "newsrooms"
}

// GormPGPersister is implements the Newsroom Persister interface
type GormPGPersister struct {
	DB *gorm.DB
}

// NewGormPGPersister takes information about the db and returns a newsroom persister that uses gorm and postgres
func NewGormPGPersister(host string, port int, user string, password string, dbname string) (*GormPGPersister, error) {
	newsroomGormPGPersister := &GormPGPersister{}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		return newsroomGormPGPersister, errors.Wrap(err, "error connecting to gorm")
	}
	newsroomGormPGPersister.DB = db

	newsroomGormPGPersister.DB.DB().SetMaxOpenConns(maxOpenConns)
	newsroomGormPGPersister.DB.DB().SetMaxIdleConns(maxIdleConns)
	newsroomGormPGPersister.DB.DB().SetConnMaxLifetime(connMaxLifetime)
	return newsroomGormPGPersister, nil
}

// NewGormPGPersisterWithDB uses an existing gorm.DB struct to create a new GormPGPersister.
// This is useful if we want to reuse existing connections
func NewGormPGPersisterWithDB(db *gorm.DB) (*GormPGPersister, error) {
	newsroomGormPGPersister := &GormPGPersister{}
	newsroomGormPGPersister.DB = db
	return newsroomGormPGPersister, nil
}

// CreateNewsroom takes a newsroom struct and saves it to the db
func (p *GormPGPersister) CreateNewsroom(newsroom *Newsroom) error {
	bys, err := json.Marshal(newsroom.Meta)
	if err != nil {
		return errors.Wrap(err, "error marshalling metadata")
	}

	newsroomGorm := Gorm{
		Name:    newsroom.Name,
		Address: newsroom.Address,
		Meta:    postgres.Jsonb{RawMessage: bys},
	}

	if err := p.DB.Create(&newsroomGorm).Error; err != nil {
		return err
	}

	newsroom.ID = newsroomGorm.ID

	return nil
}

// UpdateNewsroom takes a newsroom struct that has an id and updates it with new values
func (p *GormPGPersister) UpdateNewsroom(newsroom *Newsroom) error {
	newsroomGorm := Gorm{}

	if err := p.DB.First(&newsroomGorm, newsroom.ID).Error; err != nil {
		return err
	}

	newsroomGorm.Name = newsroom.Name
	newsroomGorm.Address = newsroom.Address

	bys, err := json.Marshal(newsroom.Meta)
	if err != nil {
		return errors.Wrap(err, "error marshalling metadata")
	}
	newsroomGorm.Meta = postgres.Jsonb{RawMessage: bys}

	err = p.DB.Save(&newsroomGorm).Error
	return err
}

// AddArticle adds an article to a newsroom with the given ID
func (p *GormPGPersister) AddArticle(newsroomID uint, newArticle *carticle.Article) error {
	newsroomGorm := Gorm{}

	articleGorm := article.Gorm{}
	if err := articleGorm.PopulateFromArticle(newArticle); err != nil {
		return err
	}

	if err := p.DB.First(&newsroomGorm, newsroomID).Error; err != nil {
		return err
	}

	if err := p.DB.Model(&newsroomGorm).Association("Articles").Append(&articleGorm).Error; err != nil {
		return err
	}

	return nil
}

// Newsrooms returns the list of newsrooms
func (p *GormPGPersister) Newsrooms() ([]*Newsroom, error) {
	newsroomGorms := []Gorm{}

	if err := p.DB.Find(&newsroomGorms).Error; err != nil {
		return nil, err
	}

	newsrooms := make([]*Newsroom, len(newsroomGorms))
	for ind, nr := range newsroomGorms {
		newsroom := &Newsroom{}
		newsroom.ID = nr.ID
		newsroom.Name = nr.Name
		newsroom.Address = nr.Address

		// Convert meta to map
		var meta *Meta
		err := json.Unmarshal(nr.Meta.RawMessage, &meta)
		if err != nil {
			log.Errorf("error unmarshalling meta: err: %v", err)
			continue
		}
		newsroom.Meta = meta

		newsrooms[ind] = newsroom
	}

	return newsrooms, nil
}

// NewsroomByID returns the newsroom with the given ID if its found
func (p *GormPGPersister) NewsroomByID(newsroomID uint) (*Newsroom, error) {
	newsroomGorm := Gorm{}

	if err := p.DB.First(&newsroomGorm, newsroomID).Error; err != nil {
		return nil, err
	}

	newsroom := &Newsroom{}
	newsroom.ID = newsroomGorm.ID
	newsroom.Name = newsroomGorm.Name
	newsroom.Address = newsroomGorm.Address

	var meta *Meta
	err := json.Unmarshal(newsroomGorm.Meta.RawMessage, &meta)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshalling meta")
	}
	newsroom.Meta = meta

	return newsroom, nil
}

// GetArticlesForNewsroom returns all the articles for a newsroom with the given ID
func (p *GormPGPersister) GetArticlesForNewsroom(newsroomID uint) ([]carticle.Article, error) {
	newsroomGorm := Gorm{}

	if err := p.DB.Preload("Articles").First(&newsroomGorm, newsroomID).Error; err != nil {
		return nil, err
	}

	return p.convertedArticles(newsroomGorm)
}

// GetArticlesForNewsroomIndexedSinceDate returns all articles for a newsroom indexed after the date
func (p *GormPGPersister) GetArticlesForNewsroomIndexedSinceDate(newsroomID uint, date time.Time) ([]carticle.Article, error) {
	newsroomGorm := Gorm{}
	if err := p.DB.Preload("Articles", "indexed_timestamp >= ?", date).First(&newsroomGorm, newsroomID).Error; err != nil {
		return nil, err
	}

	return p.convertedArticles(newsroomGorm)
}

// GetLatestArticleForNewsroom returns the latest article for a newsroom with the given ID
func (p *GormPGPersister) GetLatestArticleForNewsroom(newsroomID uint) (*carticle.Article, error) {
	newsroomGorm := Gorm{}

	sortFunc := func(db *gorm.DB) *gorm.DB {
		return db.Limit(1).Order("articles.article_metadata->>'RevisionDate' DESC")
	}

	if err := p.DB.Preload("Articles", sortFunc).First(&newsroomGorm, newsroomID).Error; err != nil {
		return nil, err
	}

	if len(newsroomGorm.Articles) == 0 {
		return nil, ErrNoArticles
	}

	art := newsroomGorm.Articles[0]
	convertedArticle, err := art.ConvertToArticle()
	if err != nil {
		return nil, err
	}

	return convertedArticle, nil
}

func (p *GormPGPersister) convertedArticles(newsroomGorm Gorm) ([]carticle.Article, error) {
	articles := make([]carticle.Article, len(newsroomGorm.Articles))
	for i, a := range newsroomGorm.Articles {
		convertedArticle, err := a.ConvertToArticle()
		if err != nil {
			return nil, err
		}
		articles[i] = *convertedArticle
	}

	return articles, nil
}

package newsroom

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/joincivil/go-common-priv/pkg/models/article"
	"github.com/pkg/errors"
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
	Address  string         `gorm:"unique;not null"`
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

// CreateNewsroom takes a newsroom struct and saves it to the db
func (p *GormPGPersister) CreateNewsroom(newsroom *Newsroom) error {
	newsroomGorm := Gorm{
		Name:    newsroom.Name,
		Address: newsroom.Address,
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

	err := p.DB.Save(&newsroomGorm).Error
	return err
}

// AddArticle adds an article to a newsroom with the given ID
func (p *GormPGPersister) AddArticle(newsroomID uint, newArticle *article.Article) error {
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

	return newsroom, nil
}

// GetArticlesForNewsroom returns all the articles for a newsroom with the given ID
func (p *GormPGPersister) GetArticlesForNewsroom(newsroomID uint) ([]article.Article, error) {
	newsroomGorm := Gorm{}

	if err := p.DB.Preload("Articles").First(&newsroomGorm, newsroomID).Error; err != nil {
		return nil, err
	}
	articles := make([]article.Article, len(newsroomGorm.Articles))

	for i, a := range newsroomGorm.Articles {
		convertedArticle, err := a.ConvertToArticle()
		if err != nil {
			return nil, err
		}
		articles[i] = *convertedArticle
	}

	return articles, nil
}

// GetLatestArticleForNewsroom returns the latest article for a newsroom with the given ID
func (p *GormPGPersister) GetLatestArticleForNewsroom(newsroomID uint) (*article.Article, error) {
	newsroomGorm := Gorm{}

	sortFunc := func(db *gorm.DB) *gorm.DB {
		return db.Limit(1).Order("articles.article_metadata->>'RevisionDate' DESC")
	}

	if err := p.DB.Preload("Articles", sortFunc).First(&newsroomGorm, newsroomID).Error; err != nil {
		return nil, err
	}

	if len(newsroomGorm.Articles) == 0 {
		return nil, errors.New("no article found")
	}

	art := newsroomGorm.Articles[0]
	convertedArticle, err := art.ConvertToArticle()
	if err != nil {
		return nil, err
	}

	return convertedArticle, nil
}

package newsroom

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/joincivil/go-common-priv/pkg/models/article"
	"time"
	"github.com/pkg/errors"
)

const (
	// Could make this configurable later if needed
	maxOpenConns    = 50
	maxIdleConns    = 10
	connMaxLifetime = time.Second * 1800 // 30 mins
)

type NewsroomGorm struct {
	gorm.Model
	Name     string
	Address  string            `gorm:"unique;not null"`
	Articles []article.ArticleGorm `gorm:"foreignkey:NewsroomAddress;association_foreignkey:Address"`
}

type NewsroomGormPGPersister struct {
	DB *gorm.DB
	version string
}

func NewNewsroomGormPGPersister(host string, port int, user string, password string, dbname string)  (*NewsroomGormPGPersister, error) {
	newsroomGormPGPersister := &NewsroomGormPGPersister{}
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

func(p *NewsroomGormPGPersister) CreateNewsroom(newsroom *Newsroom) error {
	newsroomGorm := NewsroomGorm{
		Name: newsroom.Name,
		Address: newsroom.Address,
	}

	if err := p.DB.Create(&newsroomGorm).Error; err != nil {
		return err
	}

	newsroom.ID = newsroomGorm.ID
	
	return nil
}

func(p *NewsroomGormPGPersister) UpdateNewsroom(newsroom *Newsroom) error {
	newsroomGorm := NewsroomGorm{}

	if err := p.DB.First(&newsroomGorm, newsroom.ID).Error; err != nil {
		return err
	}

	newsroomGorm.Name = newsroom.Name
	newsroomGorm.Address = newsroom.Address

	if err := p.DB.Save(&newsroomGorm).Error; err != nil {
		return err
	}
	return nil
}

func(p *NewsroomGormPGPersister) AddArticle(newsroomID uint, newArticle *article.Article) error {
	newsroomGorm := NewsroomGorm{}

	articleGorm := article.ArticleGorm{}
	if err := articleGorm.PopulateFromArticle(newArticle); err != nil {
		return err;
	}

	if err := p.DB.First(&newsroomGorm, newsroomID).Error; err != nil {
		return err
	}

	if err := p.DB.Model(&newsroomGorm).Association("Articles").Append(&articleGorm).Error; err != nil {
		return err
	}

	return nil
}

func(p *NewsroomGormPGPersister) NewsroomByID(newsroomID uint) (*Newsroom, error) {
	newsroomGorm := NewsroomGorm{}

	if err := p.DB.First(&newsroomGorm, newsroomID).Error; err != nil {
		return nil, err
	}

	newsroom := &Newsroom{}
	newsroom.ID = newsroomGorm.ID
	newsroom.Name = newsroomGorm.Name
	newsroom.Address = newsroomGorm.Address

	return newsroom, nil
}

func(p *NewsroomGormPGPersister) GetArticlesForNewsroom(newsroomID uint)  ([]article.Article, error) {
	newsroomGorm := NewsroomGorm{}

	if err := p.DB.First(&newsroomGorm, newsroomID).Error; err != nil {
		return nil, err
	}

	if err := p.DB.Model(&newsroomGorm).Related(&article.ArticleGorm{}).Error; err != nil {
		return nil, err
	}
	articles := make([]article.Article, len(newsroomGorm.Articles))

	for i, a := range newsroom.Articles {
		convertedArticle, err := a.ConvertToArticle()
		if err != nil {
			return nil, err
		}
		articles[i] = *convertedArticle
	}

	return articles, nil
}



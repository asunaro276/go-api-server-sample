package entities

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"gorm.io/gorm"
)

type Content struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string         `gorm:"type:varchar(200);not null" json:"title"`
	Body        string         `gorm:"type:text;not null" json:"body"`
	ContentType string         `gorm:"type:varchar(50);not null" json:"content_type"`
	Author      string         `gorm:"type:varchar(100);not null" json:"author"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Content) TableName() string {
	return "contents"
}

var (
	ErrInvalidTitle       = errors.New("タイトルは1文字以上200文字以下で入力してください")
	ErrInvalidBody        = errors.New("本文は1文字以上で入力してください")
	ErrInvalidContentType = errors.New("コンテンツタイプは article, blog, news, page のいずれかを指定してください")
	ErrInvalidAuthor      = errors.New("作成者名は1文字以上100文字以下で入力してください")
)

var validContentTypes = map[string]bool{
	"article": true,
	"blog":    true,
	"news":    true,
	"page":    true,
}

func NewContent(title, body, contentType, author string) (*Content, error) {
	content := &Content{
		Title:       strings.TrimSpace(title),
		Body:        strings.TrimSpace(body),
		ContentType: strings.TrimSpace(contentType),
		Author:      strings.TrimSpace(author),
	}

	if err := content.Validate(); err != nil {
		return nil, err
	}

	return content, nil
}

func (c *Content) Validate() error {
	titleLen := utf8.RuneCountInString(c.Title)
	if titleLen == 0 || titleLen > 200 {
		return ErrInvalidTitle
	}

	if utf8.RuneCountInString(c.Body) == 0 {
		return ErrInvalidBody
	}

	if !validContentTypes[c.ContentType] {
		return ErrInvalidContentType
	}

	authorLen := utf8.RuneCountInString(c.Author)
	if authorLen == 0 || authorLen > 100 {
		return ErrInvalidAuthor
	}

	return nil
}

func (c *Content) Update(title, body, contentType, author string) error {
	newContent := &Content{
		Title:       strings.TrimSpace(title),
		Body:        strings.TrimSpace(body),
		ContentType: strings.TrimSpace(contentType),
		Author:      strings.TrimSpace(author),
	}

	if err := newContent.Validate(); err != nil {
		return err
	}

	c.Title = newContent.Title
	c.Body = newContent.Body
	c.ContentType = newContent.ContentType
	c.Author = newContent.Author

	return nil
}

func (c *Content) IsDeleted() bool {
	return c.DeletedAt.Valid
}

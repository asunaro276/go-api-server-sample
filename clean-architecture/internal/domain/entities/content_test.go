package entities

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ContentTestSuite struct {
	suite.Suite
}

func (suite *ContentTestSuite) TestNewContent() {
	suite.Run("正常なコンテンツが作成できる", func() {
		content, err := NewContent("テストタイトル", "テスト本文", "article", "テスト作成者")

		assert.NoError(suite.T(), err)
		assert.NotNil(suite.T(), content)
		assert.Equal(suite.T(), "テストタイトル", content.Title)
		assert.Equal(suite.T(), "テスト本文", content.Body)
		assert.Equal(suite.T(), "article", content.ContentType)
		assert.Equal(suite.T(), "テスト作成者", content.Author)
	})

	suite.Run("空白文字が自動でトリミングされる", func() {
		content, err := NewContent("  テストタイトル  ", "  テスト本文  ", "  article  ", "  テスト作成者  ")

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "テストタイトル", content.Title)
		assert.Equal(suite.T(), "テスト本文", content.Body)
		assert.Equal(suite.T(), "article", content.ContentType)
		assert.Equal(suite.T(), "テスト作成者", content.Author)
	})
}

func (suite *ContentTestSuite) TestValidation() {
	suite.Run("タイトルのバリデーション", func() {
		suite.Run("空のタイトルはエラー", func() {
			_, err := NewContent("", "テスト本文", "article", "テスト作成者")
			assert.Equal(suite.T(), ErrInvalidTitle, err)
		})

		suite.Run("201文字のタイトルはエラー", func() {
			longTitle := ""
			for i := 0; i < 201; i++ {
				longTitle += "あ"
			}
			_, err := NewContent(longTitle, "テスト本文", "article", "テスト作成者")
			assert.Equal(suite.T(), ErrInvalidTitle, err)
		})

		suite.Run("200文字のタイトルは有効", func() {
			longTitle := ""
			for i := 0; i < 200; i++ {
				longTitle += "あ"
			}
			_, err := NewContent(longTitle, "テスト本文", "article", "テスト作成者")
			assert.NoError(suite.T(), err)
		})
	})

	suite.Run("本文のバリデーション", func() {
		suite.Run("空の本文はエラー", func() {
			_, err := NewContent("テストタイトル", "", "article", "テスト作成者")
			assert.Equal(suite.T(), ErrInvalidBody, err)
		})
	})

	suite.Run("コンテンツタイプのバリデーション", func() {
		validTypes := []string{"article", "blog", "news", "page"}
		for _, contentType := range validTypes {
			suite.Run(fmt.Sprintf("%sは有効", contentType), func() {
				_, err := NewContent("テストタイトル", "テスト本文", contentType, "テスト作成者")
				assert.NoError(suite.T(), err)
			})
		}

		suite.Run("無効なコンテンツタイプはエラー", func() {
			_, err := NewContent("テストタイトル", "テスト本文", "invalid", "テスト作成者")
			assert.Equal(suite.T(), ErrInvalidContentType, err)
		})
	})

	suite.Run("作成者のバリデーション", func() {
		suite.Run("空の作成者はエラー", func() {
			_, err := NewContent("テストタイトル", "テスト本文", "article", "")
			assert.Equal(suite.T(), ErrInvalidAuthor, err)
		})

		suite.Run("101文字の作成者はエラー", func() {
			longAuthor := ""
			for i := 0; i < 101; i++ {
				longAuthor += "あ"
			}
			_, err := NewContent("テストタイトル", "テスト本文", "article", longAuthor)
			assert.Equal(suite.T(), ErrInvalidAuthor, err)
		})

		suite.Run("100文字の作成者は有効", func() {
			longAuthor := ""
			for i := 0; i < 100; i++ {
				longAuthor += "あ"
			}
			_, err := NewContent("テストタイトル", "テスト本文", "article", longAuthor)
			assert.NoError(suite.T(), err)
		})
	})
}

func (suite *ContentTestSuite) TestUpdate() {
	suite.Run("正常に更新できる", func() {
		content, _ := NewContent("元のタイトル", "元の本文", "article", "元の作成者")

		err := content.Update("新しいタイトル", "新しい本文", "blog", "新しい作成者")

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), "新しいタイトル", content.Title)
		assert.Equal(suite.T(), "新しい本文", content.Body)
		assert.Equal(suite.T(), "blog", content.ContentType)
		assert.Equal(suite.T(), "新しい作成者", content.Author)
	})

	suite.Run("無効なデータで更新するとエラー", func() {
		content, _ := NewContent("元のタイトル", "元の本文", "article", "元の作成者")

		err := content.Update("", "新しい本文", "blog", "新しい作成者")

		assert.Equal(suite.T(), ErrInvalidTitle, err)
		// 元のデータが保持されることを確認
		assert.Equal(suite.T(), "元のタイトル", content.Title)
	})
}

func (suite *ContentTestSuite) TestIsDeleted() {
	suite.Run("削除されていないコンテンツはfalseを返す", func() {
		content, _ := NewContent("テストタイトル", "テスト本文", "article", "テスト作成者")
		assert.False(suite.T(), content.IsDeleted())
	})
}

func TestContentTestSuite(t *testing.T) {
	suite.Run(t, new(ContentTestSuite))
}

package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Store defines all functions to execute db queries and transactions
type Store interface {
	Querier
	ArticleTx(ctx context.Context, arg ArticleTxParams) (ArticleTxResult, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	db *sql.DB
	*Queries
}

// NewStore creates a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// ExecTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// ArticleTxParams contains the input parameters of the article
type ArticleTxParams struct {
	Title       string    `json:"title"`
	AuthorId    int32     `json:"author_id"`
	Content     string    `json:"content"`
	Images      []string  `json:"images"`
	Categories  []int32   `json:"categories"`
	Tags        []int32   `json:"tags"`
	PublishedAt time.Time `json:"published_at"`
}

// ArticleTxResult is the result of the article transaction
type ArticleTxResult struct {
	Id          int32          `json:"id"`
	Title       string         `json:"title" `
	AuthorId    int32          `json:"author_id" `
	Content     string         `json:"content" `
	Images      []createdImage `json:"images"`
	Categories  []Category     `json:"categories"`
	Tags        []Tag          `json:"tags"`
	PublishedAt time.Time      `json:"published_at"`
}

type createdImage struct {
	Id  int32  `json:"id"`
	Url string `json:"url"`
}

// ArticleTx performs a creation of article and related data.
// It creates the article, article category, article tag, and images within a database transaction
func (store *SQLStore) ArticleTx(ctx context.Context, arg ArticleTxParams) (ArticleTxResult, error) {
	var result ArticleTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		articleResult, err := q.CreateArticle(ctx, CreateArticleParams{
			Title:       arg.Title,
			AuthorID:    arg.AuthorId,
			Content:     arg.Content,
			PublishedAt: arg.PublishedAt,
		})
		if err != nil {
			return err
		}

		var images []createdImage
		for _, image := range arg.Images {
			articleImageResult, err := q.CreateImage(ctx, CreateImageParams{
				ArticleID: articleResult.ID,
				Url:       image,
			})
			if err != nil {
				return err
			}
			images = append(images, createdImage{
				Id:  articleImageResult.ID,
				Url: articleImageResult.Url,
			})
		}

		var categories []Category
		for _, categoryId := range arg.Categories {
			existingCategory, err := q.FindCategory(ctx, categoryId)
			if err != nil {
				return err
			}
			_, err = q.CreateArticleCategory(ctx, CreateArticleCategoryParams{
				ArticleID:  articleResult.ID,
				CategoryID: existingCategory.ID,
			})
			if err != nil {
				return err
			}
			categories = append(categories, existingCategory)
		}

		var tags []Tag
		for _, tagId := range arg.Tags {
			existingTag, err := q.FindTag(ctx, tagId)
			if err != nil {
				return err
			}
			_, err = q.CreateArticleTag(ctx, CreateArticleTagParams{
				ArticleID: articleResult.ID,
				TagID:     existingTag.ID,
			})
			if err != nil {
				return err
			}
			tags = append(tags, existingTag)
		}

		result = ArticleTxResult{
			Id:          articleResult.ID,
			AuthorId:    articleResult.AuthorID,
			Title:       articleResult.Title,
			Content:     articleResult.Content,
			Images:      images,
			Categories:  categories,
			Tags:        tags,
			PublishedAt: arg.PublishedAt,
		}
		return err
	})

	return result, err
}

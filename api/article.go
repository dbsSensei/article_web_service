package api

import (
	"encoding/json"
	"fmt"
	db "github.com/dbssensei/articlewebservice/db/sqlc"
	"github.com/dbssensei/articlewebservice/token"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type createArticleRequest struct {
	Title      string   `json:"title" binding:"required"`
	Content    string   `json:"content" binding:"required"`
	Images     []string `json:"images"`
	Categories []int32  `json:"categories"`
	Tags       []int32  `json:"tags"`
}

func (server *Server) createArticle(ctx *gin.Context) {
	var req createArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.ArticleTxParams{
		Title:       req.Title,
		AuthorId:    authPayload.UserId,
		Content:     req.Content,
		Images:      req.Images,
		Categories:  req.Categories,
		Tags:        req.Tags,
		PublishedAt: time.Now(),
	}

	article, err := server.store.ArticleTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, article)
}

// TODO update article
//type updateArticleRequest struct {
//	Title   string `json:"title" binding:"required"`
//	Content string `json:"content" binding:"required"`
//}
//
//func (server *Server) updateArticle(ctx *gin.Context) {
//	var req updateArticleRequest
//	if err := ctx.ShouldBindJSON(&req); err != nil {
//		ctx.JSON(http.StatusBadRequest, errorResponse(err))
//		return
//	}
//
//	articleId, _ := strconv.Atoi(ctx.Param("id"))
//	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
//
//	updatedArticle, err := server.store.Update(ctx, db.UpdateArticleParams{
//		ID:       int32(articleId),
//		AuthorID: authPayload.UserId,
//		Name: sql.NullString{
//			String: req.Name,
//			Valid:  len(req.Name) > 0,
//		},
//		Price: sql.NullInt32{
//			Int32: req.Price,
//			Valid: req.Price > 0,
//		},
//		Status: sql.NullString{
//			String: req.Status,
//			Valid:  len(req.Status) > 0,
//		},
//		CategoryID: sql.NullInt32{
//			Int32: req.CategoryID,
//			Valid: req.CategoryID > 0,
//		},
//	})
//	if err != nil {
//		if pqErr, ok := err.(*pq.Error); ok {
//			switch pqErr.Code.Name() {
//			case "unique_violation":
//				ctx.JSON(http.StatusForbidden, errorResponse(err))
//				return
//			}
//		}
//		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
//		return
//	}
//
//	rsp := updatedArticle
//	ctx.JSON(http.StatusOK, rsp)
//}
//

type getArticle struct {
	ID          int32     `json:"id"`
	Title       string    `json:"title"`
	AuthorID    int32     `json:"author_id"`
	Content     string    `json:"content"`
	PublishedAt time.Time `json:"published_at"`
	AuthorName  string    `json:"author_name"`
	TotalViews  int64     `json:"total_views"`
}

func (server *Server) getArticles(ctx *gin.Context) {
	articles, err := server.store.GetArticles(ctx, db.GetArticlesParams{
		AuthorName: ctx.Query("author"),
		Search:     ctx.Query("search"),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var getArticles []getArticle
	for _, article := range articles {
		getArticles = append(getArticles, getArticle{
			ID:          article.ID,
			Title:       article.Title,
			AuthorID:    article.AuthorID,
			Content:     article.Content,
			PublishedAt: article.PublishedAt,
			AuthorName:  article.AuthorName.String,
			TotalViews:  article.TotalViews,
		})
	}

	ctx.JSON(http.StatusOK, getArticles)
}

type fullArticle struct {
	Id          int32         `json:"id"`
	Title       string        `json:"title" `
	AuthorId    int32         `json:"author_id" `
	Content     string        `json:"content" `
	Images      []image       `json:"images"`
	Categories  []db.Category `json:"categories"`
	Tags        []db.Tag      `json:"tags"`
	TotalViews  int64         `json:"total_views"`
	PublishedAt time.Time     `json:"published_at"`
}

type image struct {
	Id  int32  `json:"id"`
	Url string `json:"url"`
}

func (server *Server) getArticle(ctx *gin.Context) {
	articleId, _ := strconv.Atoi(ctx.Param("id"))

	articleInCache, err := server.redis.Get(ctx, fmt.Sprintf("article:%d", articleId)).Result()
	if err != nil {
		fmt.Println(err)
	} else {
		// Deserialize the user JSON string back into a struct
		var parsedArticleInCache fullArticle
		err = json.Unmarshal([]byte(articleInCache), &parsedArticleInCache)
		if err != nil {
			panic(err)
		}

		ctx.JSON(http.StatusOK, parsedArticleInCache)
		return
	}

	article, err := server.store.GetArticle(ctx, int32(articleId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.store.CreateView(ctx, db.CreateViewParams{
		ArticleID: article.ID,
		ViewDate:  time.Now(),
	})
	if err != nil {
		fmt.Println("Error while create article view", err)
	}

	totalViews, err := server.store.CountArticleViewsByArticleId(ctx, article.ID)
	if err != nil {
		fmt.Println("Error while get article total view", err)
	}

	storedImages, err := server.store.GetImageByArticleId(ctx, article.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	var images []image
	for _, storedImage := range storedImages {
		images = append(images, image{
			Id:  storedImage.ID,
			Url: storedImage.Url,
		})
	}

	storedArticleCategories, err := server.store.GetArticleCategoryByArticleId(ctx, article.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	var categories []db.Category
	for _, articleCategory := range storedArticleCategories {
		categories = append(categories, db.Category{
			ID:   articleCategory.ID,
			Name: articleCategory.Name,
		})
	}

	storedArticleTags, err := server.store.GetArticleTagByArticleId(ctx, article.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	var tags []db.Tag
	for _, articleTag := range storedArticleTags {
		tags = append(tags, db.Tag{
			ID:   articleTag.ID,
			Name: articleTag.Name,
		})
	}

	rsp := fullArticle{
		Id:          article.ID,
		Title:       article.Title,
		AuthorId:    article.AuthorID,
		Content:     article.Content,
		Images:      images,
		Categories:  categories,
		Tags:        tags,
		TotalViews:  totalViews,
		PublishedAt: article.PublishedAt,
	}

	articleJSON, err := json.Marshal(rsp)
	if err != nil {
		fmt.Println(err)
	}

	err = server.redis.Set(ctx, fmt.Sprintf("article:%d", articleId), articleJSON, time.Second*60).Err()
	if err != nil {
		fmt.Println(err)
	}

	ctx.JSON(http.StatusOK, rsp)
}

// TODO delete article
//func (server *Server) deleteArticle(ctx *gin.Context) {
//	articleId, _ := strconv.Atoi(ctx.Param("id"))
//	err := server.store.DeleteArticle(ctx, int32(articleId))
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
//		return
//	}
//
//	rsp := deleteArticleResponse{
//		Status: "success",
//	}
//	ctx.JSON(http.StatusOK, rsp)
//}

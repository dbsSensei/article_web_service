package api

import (
	"fmt"
	"github.com/go-redis/redis/v8"

	db "github.com/dbssensei/articlewebservice/db/sqlc"
	"github.com/dbssensei/articlewebservice/token"
	"github.com/dbssensei/articlewebservice/util"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our marketplace service.
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
	redis      *redis.Client
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store db.Store, redis *redis.Client) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		redis:      redis,
	}

	gin.SetMode(server.config.GinMode)
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/tokens/renew_access", server.renewAccessToken)
	router.GET("/articles", server.getArticles)
	router.GET("/articles/:id", server.getArticle)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.PATCH("/users", server.updateUser)
	authRoutes.POST("/articles", server.createArticle)
	//authRoutes.PATCH("/articles/:id", server.updateArticle)
	//authRoutes.DELETE("/articles/:id", server.deleteArticle)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

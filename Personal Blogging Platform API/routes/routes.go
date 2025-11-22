package routes

import (
	"github.com/RohithBN/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PATCH", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Type"}
	r.Use(cors.New(config))

	r.GET("/health",handler.HealthCheck)

	apiRoutes := r.Group("/api")
	{
		apiRoutes.GET("/articles", handler.GetArticles)
		apiRoutes.POST("/articles", handler.CreateArticle)
		apiRoutes.DELETE("/article/:id",handler.DeleteArticle)
		apiRoutes.PUT("/article/:id",handler.UpdateArticle)
		apiRoutes.GET("/article/:id",handler.GetArticleByID)
	}


	return r
		
}

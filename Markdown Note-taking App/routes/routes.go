package routes

import (
	"github.com/RohithBN/handler"
	"github.com/gin-gonic/gin"
)




func SetupRoutes() *gin.Engine{
	r:=gin.Default()
	
	r.GET("/health", func(c *gin.Context){
		c.JSON(200, gin.H{
			"status":"OK",
			"message":"Markdown Note-taking App is running",
		})
	})

	apiRoutes:= r.Group("/api")
	{
		apiRoutes.POST("/notes",handler.CreateNote)
		apiRoutes.GET("/notes/:id",handler.GetNoteById)
		apiRoutes.GET("/notes/:id/grammar",handler.CheckGrammar)
		apiRoutes.POST("/notes/:id/attachments",handler.AddAttachmentToNote)
		apiRoutes.GET("/notes/:id/attachments",handler.GetNoteAttachments)
		apiRoutes.GET("/notes/:id/rendered",handler.GetRenderedNote)
	}
	return r
}
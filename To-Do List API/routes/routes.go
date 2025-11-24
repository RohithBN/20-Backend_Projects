package routes

import (
	"github.cim/RohithBN/handler"
	"github.cim/RohithBN/middleware"
	"github.com/gin-gonic/gin"
)


func SetupRoutes() *gin.Engine{

	r:= gin.Default()

	publicRoutes:= r.Group("/api")
	{
		publicRoutes.POST("/register",handler.Register)
		publicRoutes.POST("/login",handler.Login)

	}

	protectedRoutes:= r.Group("/api")

	protectedRoutes.Use(middleware.AuthMiddleware())
	{
		protectedRoutes.GET("/todos",handler.GetTodos)
		protectedRoutes.PATCH("/todos/:id",handler.UpdateStatus)
	    protectedRoutes.DELETE("/todos/:id",handler.DeleteTodo)
		protectedRoutes.POST("/todos",handler.CreateTodo)
		protectedRoutes.GET("/todos/:status",handler.GetTodosByStatus)
	}

	return r;


}
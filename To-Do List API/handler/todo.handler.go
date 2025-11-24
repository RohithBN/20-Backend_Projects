package handler

import (
    "fmt"
    "net/http"

    "github.cim/RohithBN/lib"
    "github.cim/RohithBN/types"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

func GetTodos(c *gin.Context) {
    var todos []types.ToDoItem
    
    // Get token from context
    userToken, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    // Type assert to *jwt.Token
    token, ok := userToken.(*jwt.Token)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token format"})
        return
    }

    // Extract claims
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse claims"})
        return
    }

    // Get user ID from claims
    userID, ok := claims["userId"]
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in token"})
        return
    }

    fmt.Println("UserId:", userID)

    query := "SELECT todo_id, task, status, created_at, created_by FROM todos WHERE created_by = ?"
    rows, err := lib.DB.Query(query, userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch todos"})
        return
    }
    defer rows.Close()

    for rows.Next() {
        var todo types.ToDoItem
        err := rows.Scan(&todo.TodoId, &todo.Task, &todo.Status, &todo.CreatedAt, &todo.CreatedBy)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse todo item"})
            return
        }
        todos = append(todos, todo)
    }

    c.JSON(http.StatusOK, gin.H{"todos": todos})
}

func CreateTodo(c *gin.Context) {
    var newTodo types.ToDoItem

    // Get token from context
    userToken, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    // Type assert to *jwt.Token
    token, ok := userToken.(*jwt.Token)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token format"})
        return
    }

    // Extract claims
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse claims"})
        return
    }

    // Get user ID from claims
    userID, ok := claims["userId"]
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in token"})
        return
    }

    if err := c.ShouldBindJSON(&newTodo); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    query := "INSERT INTO todos (task, status, created_by) VALUES (?, ?, ?)"
    result, err := lib.DB.Exec(query, newTodo.Task, newTodo.Status, userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create todo"})
        return
    }

    todoId, err := result.LastInsertId()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve todo ID"})
        return
    }

    newTodo.TodoId = uint(todoId)
    newTodo.CreatedBy = fmt.Sprintf("%v", userID)
    c.JSON(http.StatusCreated, gin.H{"message": "Todo created successfully", "todo": newTodo})
}



func UpdateStatus (c *gin.Context){

	userToken , exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "User not authenticated"})
		return
	}

	token , ok := userToken.(*jwt.Token)
	if !ok {
		c.JSON(500, gin.H{"error": "Invalid token format"})
		return
	}

	claims , ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(500, gin.H{"error": "Failed to parse claims"})
		return
	}
	
	userID , ok := claims["userId"]
	if !ok {
		c.JSON(500, gin.H{"error": "User ID not found in token"})
		return
	}
	
	fmt.Println("UserId:", userID)
	var StatusObject struct{
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&StatusObject); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}


	todoId := c.Param("id")
	
	query := "UPDATE todos SET status = ? WHERE todo_id = ?"
	_, err := lib.DB.Exec(query, StatusObject.Status, todoId)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update todo status"})
		return
	}
	
	c.JSON(200, gin.H{"message": "Todo status updated successfully"})

}

func DeleteTodo (c *gin.Context){

	todoId := c.Param("id")

	
	query := "DELETE FROM todos WHERE todo_id = ?"
	_, err := lib.DB.Exec(query, todoId)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete todo"})
		return
	}
	
	c.JSON(200, gin.H{"message": "Todo deleted successfully"})

}

func GetTodosByStatus (c *gin.Context){
	status := c.Param("status")
	var todos []types.ToDoItem
	
	query := "SELECT todo_id, task, status, created_at, created_by FROM todos WHERE status = ?"
	rows, err := lib.DB.Query(query, status)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch todos by status"})
		return
	}
	defer rows.Close()
	
	for rows.Next() {
		var todo types.ToDoItem
		err := rows.Scan(&todo.TodoId, &todo.Task, &todo.Status, &todo.CreatedAt, &todo.CreatedBy)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to parse todo item"})
			return
		}
		todos = append(todos, todo)
	}
	c.JSON(200, gin.H{"todos": todos})
}

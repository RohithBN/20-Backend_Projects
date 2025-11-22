package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/RohithBN/database"
	"github.com/RohithBN/models"
	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "API is running"})
}

func GetArticles(c *gin.Context) {
	query := "SELECT article_id, title, content, author, created_at, updated_at FROM articles"
	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch articles", "details": err.Error()})
		return
	}
	defer rows.Close()

	var articles []models.Article

	for rows.Next() {
		var article models.Article
		err := rows.Scan(&article.ArticleId, &article.Title, &article.Content, &article.Author, &article.CreatedAt, &article.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan article", "details": err.Error()})
			return
		}
		articles = append(articles, article)
	}

	// Check for errors during iteration
	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error iterating articles", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"articles": articles, "count": len(articles)})
}

func CreateArticle(c *gin.Context) {
	var newArticle models.Article
	if err := c.BindJSON(&newArticle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if newArticle.Title == "" || newArticle.Content == "" || newArticle.Author == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title, Content and Author are required"})
		return
	}

	query := "INSERT INTO articles (title, content, author) VALUES (?, ?, ?)"
	result, err := database.DB.Exec(query, newArticle.Title, newArticle.Content, newArticle.Author)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article", "details": err.Error()})
		return
	}

	articleID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve article ID"})
		return
	}

	newArticle.ArticleId = uint(articleID)
	c.JSON(http.StatusCreated, gin.H{"article": newArticle})
}

func DeleteArticle(c *gin.Context) {
	articleId := c.Param("id")

	if articleId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Article ID is required"})
		return
	}

	query := "DELETE FROM articles WHERE article_id = ?"
	result, err := database.DB.Exec(query, articleId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article", "details": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check deletion"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}

func UpdateArticle(c *gin.Context) {
	articleID := c.Param("id")

	if articleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Article ID is required"})
		return
	}

	var updatedArticle models.Article
	if err := c.BindJSON(&updatedArticle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	query := "UPDATE articles SET title = ?, content = ?, author = ?, updated_at = ? WHERE article_id = ?"
	result, err := database.DB.Exec(query, updatedArticle.Title, updatedArticle.Content, updatedArticle.Author, time.Now(), articleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article", "details": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check update"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article updated successfully"})
}

func GetArticleByID(c *gin.Context) {
	articleId := c.Param("id")

	if articleId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Article ID is required"})
		return
	}

	query := "SELECT article_id, title, content, author, created_at, updated_at FROM articles WHERE article_id = ?"
	row := database.DB.QueryRow(query, articleId)

	var article models.Article
	err := row.Scan(&article.ArticleId, &article.Title, &article.Content, &article.Author, &article.CreatedAt, &article.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch article", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"article": article})
}

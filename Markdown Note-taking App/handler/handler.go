package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/RohithBN/lib"
	"github.com/RohithBN/types"
	"github.com/gin-gonic/gin"
)

func CreateNote(c *gin.Context) {
	var note types.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Input validation
	if note.Title == "" || note.MarkdownContent == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title and MarkdownContent are required"})
		return
	}

	// Save note to database (let MySQL handle timestamps)
	query := "INSERT INTO notes (title, markdown_content) VALUES (?, ?)"
	result, err := lib.DB.Exec(query, note.Title, note.MarkdownContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
		return
	}

	// Get the inserted ID
	noteId, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve note ID"})
		return
	}

	note.Id = int(noteId)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Note created successfully",
		"note":    note,
	})
}

func GetNoteById(c *gin.Context) {
	id := c.Param("id")

	var note types.Note
	query := "SELECT id, title, markdown_content, created_at, updated_at FROM notes WHERE id = ?"
	err := lib.DB.QueryRow(query, id).Scan(
		&note.Id,
		&note.Title,
		&note.MarkdownContent,
		&note.CreatedAt,
		&note.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch note"})
		return
	}

	c.JSON(http.StatusOK, note)
}

func GetAllNotes(c *gin.Context) {
	query := "SELECT id, title, markdown_content, created_at, updated_at FROM notes ORDER BY created_at DESC"
	rows, err := lib.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
		return
	}
	defer rows.Close()

	var notes []types.Note
	for rows.Next() {
		var note types.Note
		err := rows.Scan(
			&note.Id,
			&note.Title,
			&note.MarkdownContent,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse notes"})
			return
		}
		notes = append(notes, note)
	}

	if len(notes) == 0 {
		c.JSON(http.StatusOK, gin.H{"notes": []types.Note{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notes": notes})
}

func UpdateNote(c *gin.Context) {
	id := c.Param("id")

	var note types.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validation
	if note.Title == "" || note.MarkdownContent == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title and MarkdownContent are required"})
		return
	}

	query := "UPDATE notes SET title = ?, markdown_content = ? WHERE id = ?"
	result, err := lib.DB.Exec(query, note.Title, note.MarkdownContent, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify update"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note updated successfully"})
}

func DeleteNote(c *gin.Context) {
	id := c.Param("id")

	query := "DELETE FROM notes WHERE id = ?"
	result, err := lib.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify deletion"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}
func AddAttachmentToNote(c *gin.Context) {
	file, err := c.FormFile("attachment")
	if err != nil {
		c.JSON(400, gin.H{"error": "Attachment file is required"})
		return
	}
	fmt.Println("Received file", file.Filename)

	//upload file to S3 here

	var attachment types.Attachment
	attachment.OriginalFileName = file.Filename
	attachment.StoredFileName = lib.GenerateUniqueFileName(file.Filename)
	attachment.FileUrl = "https://hardcoded-s3-url/" + attachment.StoredFileName
	attachment.MimeType = file.Header.Get("Content-Type")
	attachment.Size = file.Size

	// parse note id from path parameter
	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid note id"})
		return
	}
	attachment.NoteId = idInt

	query := "INSERT INTO attachments (note_id, original_file_name, stored_file_name, file_url, mime_type, size) VALUES (?,?,?,?,?,?)"
	_, err = lib.DB.Exec(query, attachment.NoteId, attachment.OriginalFileName, attachment.StoredFileName, attachment.FileUrl, attachment.MimeType, attachment.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add attachment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Attachment added successfully", "attachment": attachment})
}

func CheckGrammar(c *gin.Context) {
	id := c.Param("id")

	var note types.Note
	query := "SELECT id, title, markdown_content, created_at, updated_at FROM notes WHERE id = ?"
	err := lib.DB.QueryRow(query, id).Scan(
		&note.Id,
		&note.Title,
		&note.MarkdownContent,
		&note.CreatedAt,
		&note.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch note"})
		return
	}

	// Clean markdown
	cleanedMarkdownText := lib.CleanMarkdownLine(note.MarkdownContent)

	// Grammar check - use environment variable
	languageToolURL := os.Getenv("LANGUAGETOOL_URL")
	if languageToolURL == "" {
		languageToolURL = "http://localhost:8010" // Default for local dev
	}

	data := url.Values{}
	data.Set("text", cleanedMarkdownText)
	data.Set("language", "en-US")

	resp, err := http.PostForm(languageToolURL+"/v2/check", data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to grammar check service"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Grammar check service returned error"})
		return
	}

	var GrammarAPIResponse types.GrammarResponse
	if err := json.NewDecoder(resp.Body).Decode(&GrammarAPIResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse grammar check response"})
		return
	}

	grammarCorrections := lib.FormatGrammarResponse(&GrammarAPIResponse)

	c.JSON(http.StatusOK, gin.H{
		"note":                note,
		"grammar_corrections": grammarCorrections,
	})
}


func GetNoteAttachments(c *gin.Context) {
	id := c.Param("id")

	query := "SELECT id, note_id, original_file_name, stored_file_name, file_url, uploaded_at, mime_type, size FROM attachments WHERE note_id = ? ORDER BY uploaded_at DESC"
	rows, err := lib.DB.Query(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attachments"})
		return
	}
	defer rows.Close()

	var attachments []types.Attachment
	for rows.Next() {
		var attachment types.Attachment
		err := rows.Scan(
			&attachment.Id,
			&attachment.NoteId,
			&attachment.OriginalFileName,
			&attachment.StoredFileName,
			&attachment.FileUrl,
			&attachment.UploadedAt,
			&attachment.MimeType,
			&attachment.Size,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":	"Failed to parse attachments"})
			return
		}
		attachments = append(attachments, attachment)
	}

	if len(attachments) == 0 {
		c.JSON(http.StatusOK, gin.H{"attachments": []types.Attachment{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"attachments": attachments})
}


func GetRenderedNote(c *gin.Context) {
	id := c.Param("id")
	
	var note types.Note
	query := "SELECT id, title, markdown_content, created_at, updated_at FROM notes WHERE id = ?"
	err := lib.DB.QueryRow(query, id).Scan(
		&note.Id,
		&note.Title,
		&note.MarkdownContent,
		&note.CreatedAt,
		&note.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch note"})
		return
	}
	renderedHTML := lib.MarkdownToHTML(note.MarkdownContent)
	fmt.Println("Rendered MD->HTML")
	c.Data(http.StatusOK, "text/html; charset=utf-8", renderedHTML)
}
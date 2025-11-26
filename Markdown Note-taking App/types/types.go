package types


type Note struct {
	Id int `json:"id"`
	Title string `json:"title"`
	MarkdownContent string `json:"markdown_content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}


type Attachment struct {
	Id int `json:"id"`
	NoteId int `json:"note_id"`
	OriginalFileName string `json:"original_file_name"`
	StoredFileName string `json:"stored_file_name"`
	FileUrl string `json:"file_url"`
	UploadedAt string `json:"uploaded_at"`
	MimeType string `json:"mime_type"`
	Size int64 `json:"size"`
}


type GrammarResponse struct{
	Matches []struct{
		Message string `json:"message"`
		Sentence string `json:"sentence"`
		Replacements []struct{
			Value string `json:"value"`
		} `json:"replacements"`
		Context struct{
			Text string `json:"text"`
			Offset int `json:"offset"`
			Length int `json:"length"`
		} `json:"context"`
	} `json:"matches"`		
}

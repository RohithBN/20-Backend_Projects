package types


type ToDoItem struct{
	TodoId  uint   `json:"todo_id"`
	Task   string `json:"task"`
	Status string `json:"status"`
	CreatedAt string `json:"created_at"`
	CreatedBy string `json:"created_by"`
}

type User struct{
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	CreatedAt string `json:"created_at"`
	Todos []ToDoItem `json:"todos"`
}
package utils

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

type Gpx struct {
	Id          int    `json:"id"`
	ActivityId  int    `json:"activity_id"`
	UserID      string `json:"user_id"`
	Filename    string `json:"filename"`
	StoragePath string `json:"storage_path"`
	UploadDate  string `json:"upload_date"`
	Stats       string `json:"stats"`
}

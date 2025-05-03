package utils

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

type Gpx struct {
	ID          int    `json:"id"`
	ActivityID  int    `json:"activity_id"`
	UserID      string `json:"user_id"`
	Filename    string `json:"filename"`
	StoragePath string `json:"storage_path"`
	UploadDate  string `json:"upload_date"`
	Stats       string `json:"stats"`
}

type GpxInfo struct {
	ID         int    `json:"id"`
	ActivityID int    `json:"activity_id"`
	Filename   string `json:"filename"`
	UploadDate string `json:"upload_date"`
	Stats      string `json:"stats"`
}

type GlobalStats struct {
	TotalActivities int     `json:"total_activities"`
	TotalDistance   float64 `json:"total_distance"`
	TotalAscent     float64 `json:"total_ascent"`
	TotalDescent    float64 `json:"total_descent"`
	TotalTime       float64 `json:"total_time"`
}

type Activity struct {
	ID                int     `json:"id"`
	UserID            string  `json:"user_id"`
	Title             string  `json:"title"`
	Description       string  `json:"description"`
	StartTime         string  `json:"start_time"`
	EndTime           string  `json:"end_time"`
	CreatedAt         string  `json:"created_at"`
	Distance          float64 `json:"distance"`
	TotalAscent       float64 `json:"total_ascent"`
	TotalDescent      float64 `json:"total_descent"`
	StartingElevation float64 `json:"starting_elevation"`
	MaximumElevation  float64 `json:"maximum_elevation"`
	AverageSpeed      float64 `json:"average_speed"`
}

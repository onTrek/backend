package utils

import (
	"database/sql"
)

type User struct {
	ID        string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email     string `json:"email" example:"user@example.com"`
	Name      string `json:"name" example:"John Doe"`
	Password  string `json:"password" example:"strongPassword123"`
	CreatedAt string `json:"created_at" example:"2025-05-11T08:00:00Z"`
}

type UserID struct {
	ID string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type UserInfo struct {
	ID    string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email string `json:"email" example:"user@example.com"`
	Name  string `json:"name" example:"John Doe"`
}

type RegisterInput struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"strongPassword123"`
	Name     string `json:"name" example:"John Doe"`
}

type Gpx struct {
	ID          int    `json:"id" example:"1"`
	ActivityID  int    `json:"activity_id" example:"101"`
	UserID      string `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Filename    string `json:"filename" example:"MonteBianco.gpx"`
	StoragePath string `json:"storage_path" example:"/uploads/gpx/MonteBianco.gpx"`
	UploadDate  string `json:"upload_date" example:"2025-05-11T08:00:00Z"`
	Stats       string `json:"stats" example:"{\"distance\": 5.5, \"ascent\": 250, \"descent\": 200}"`
}

type GpxInfo struct {
	ID         int    `json:"id" example:"1"`
	ActivityID int    `json:"activity_id" example:"101"`
	Filename   string `json:"filename" example:"MonteBianco.gpx"`
	UploadDate string `json:"upload_date" example:"2025-05-11T08:00:00Z"`
	Stats      string `json:"stats" example:"{\"distance\": 5.5, \"ascent\": 250, \"descent\": 200}"`
}

type GlobalStats struct {
	TotalActivities int     `json:"total_activities" example:"100"`
	TotalDistance   float64 `json:"total_distance" example:"500.5"`
	TotalAscent     float64 `json:"total_ascent" example:"2500"`
	TotalDescent    float64 `json:"total_descent" example:"2000"`
	TotalTime       float64 `json:"total_time" example:"120.5"`
}

type Activity struct {
	ID                int     `json:"id" example:"1"`
	UserID            string  `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title             string  `json:"title" example:"Morning Run"`
	Description       string  `json:"description" example:"A refreshing morning run in the park"`
	StartTime         string  `json:"start_time" example:"2025-05-11T08:00:00Z"`
	EndTime           string  `json:"end_time" example:"2025-05-11T09:00:00Z"`
	CreatedAt         string  `json:"created_at" example:"2025-05-11T08:00:00Z"`
	Distance          float64 `json:"distance" example:"5.5"`
	TotalAscent       float64 `json:"total_ascent" example:"250"`
	TotalDescent      float64 `json:"total_descent" example:"200"`
	StartingElevation float64 `json:"starting_elevation" example:"100"`
	MaximumElevation  float64 `json:"maximum_elevation" example:"350"`
	AverageSpeed      float64 `json:"average_speed" example:"5.2"`
}

type Session struct {
	ID        int            `json:"id" example:"1"`
	CreatedBy string         `json:"created_by" example:"550e8400-e29b-41d4-a716-446655440000"`
	CreatedAt string         `json:"created_at" example:"2025-05-11T08:00:00Z"`
	ClosedAt  sql.NullString `json:"closed_at" example:"2025-05-11T09:00:00Z"`
}

type SessionDoc struct {
	ID        int    `json:"id" example:"1"`
	CreatedBy string `json:"created_by" example:"550e8400-e29b-41d4-a716-446655440000"`
	CreatedAt string `json:"created_at" example:"2025-05-11T08:00:00Z"`
	ClosedAt  struct {
		String string `json:"String" example:"2025-05-11T09:00:00Z"`
		Valid  bool   `json:"Valid" example:"true"`
	}
}

type SessionInfo struct {
	SessionID int     `json:"session_id" example:"1"`
	Latitude  float64 `json:"latitude" example:"40.7128"`
	Longitude float64 `json:"longitude" example:"-74.0060"`
	Altitude  float64 `json:"altitude" example:"10.5"`
	Accuracy  float64 `json:"accuracy" example:"5.0"`
	Time      string  `json:"time" example:"2025-05-11T08:00:00Z"`
}

type SessionInfoUpdate struct {
	Latitude  float64 `json:"latitude" example:"40.7128"`
	Longitude float64 `json:"longitude" example:"-74.0060"`
	Altitude  float64 `json:"altitude" example:"10.5"`
	Accuracy  float64 `json:"accuracy" example:"5.0"`
}

type MemberInfo struct {
	User        UserInfo          `json:"user"`
	SessionInfo SessionInfoUpdate `json:"session_info"`
	TimeStamp   string            `json:"time_stamp" example:"2025-05-11T08:00:00Z"`
}
type SessionInfoResponse struct {
	CreatedBy UserInfo       `json:"created_by"`
	CreatedAt string         `json:"created_at" example:"2025-05-11T08:00:00Z"`
	ClosedAt  sql.NullString `json:"closed_at" example:"2025-05-11T09:00:00Z"`
	Members   []MemberInfo   `json:"members"`
}

type SessionInfoResponseDoc struct {
	CreatedBy UserInfo `json:"created_by"`
	CreatedAt string   `json:"created_at" example:"2025-05-11T08:00:00Z"`
	ClosedAt  struct {
		String string `json:"String" example:"2025-05-11T09:00:00Z"`
		Valid  bool   `json:"Valid" example:"true"`
	}
	Members []MemberInfo `json:"members"`
}

type SuccessResponse struct {
	Message string `json:"message" example:"Success"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"An error occurred"`
}

type ActivityInput struct {
	StartTime         string  `json:"startTime" example:"2025-05-11T08:00:00Z"`
	EndTime           string  `json:"endTime" example:"2025-05-11T09:00:00Z"`
	Distance          float64 `json:"distance" example:"5.5"`
	TotalAscent       float64 `json:"totalAscent" example:"250"`
	TotalDescent      float64 `json:"totalDescent" example:"200"`
	StartingElevation float64 `json:"startingElevation" example:"100"`
	MaximumElevation  float64 `json:"maximumElevation" example:"350"`
	AverageSpeed      float64 `json:"averageSpeed" example:"5.2"`
}

type Login struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"strongPassword123"`
}

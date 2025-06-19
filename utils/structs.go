package utils

import (
	"database/sql"
)

type User struct {
	ID        string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email     string `json:"email" example:"user@example.com"`
	Username  string `json:"username" example:"John Doe"`
	Password  string `json:"password" example:"strongPassword123"`
	CreatedAt string `json:"created_at" example:"2025-05-11T08:00:00Z"`
}

type UserToken struct {
	Token string `json:"token" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type UserInfo struct {
	ID       string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email    string `json:"email" example:"user@example.com"`
	Username string `json:"username" example:"John Doe"`
}

type UserEssentials struct {
	ID       string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username string `json:"username" example:"John Doe"`
}

type RegisterInput struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"strongPassword123"`
	Username string `json:"username" example:"John Doe"`
}

type Gpx struct {
	ID          int    `json:"id" example:"1"`
	UserID      string `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Filename    string `json:"filename" example:"MonteBianco.gpx"`
	StoragePath string `json:"storage_path" example:"/uploads/gpx/MonteBianco.gpx"`
	UploadDate  string `json:"upload_date" example:"2025-05-11T08:00:00Z"`
	Title       string `json:"title" example:"Monte Faggeto"`
}

type GpxInfo struct {
	ID         int      `json:"id" example:"1"`
	Filename   string   `json:"filename" example:"MonteBianco.gpx"`
	UploadDate string   `json:"upload_date" example:"2025-05-11T08:00:00Z"`
	Title      string   `json:"title" example:"Monte Faggeto"`
	Stats      GPXStats `json:"stats"`
}

type GpxInfoEssential struct {
	ID       int    `json:"id" example:"1"`
	Filename string `json:"filename" example:"MonteBianco.gpx"`
}

type GPXStats struct {
	Km          float64 `json:"km" example:"14.5"`
	Duration    string  `json:"duration" example:"06:30:00"`
	Ascent      int     `json:"ascent" example:"1000"`
	Descent     int     `json:"descent" example:"1000"`
	MaxAltitude int     `json:"max_altitude" example:"2500"`
	MinAltitude int     `json:"min_altitude" example:"1500"`
}

type GpxInfoDoc struct {
	Files []GpxInfo `json:"gpx_files"`
}

type Session struct {
	ID          int              `json:"session_id" example:"1"`
	CreatedBy   string           `json:"created_by" example:"550e8400-e29b-41d4-a716-446655440000"`
	Description string           `json:"description" example:"Morning hike with friends"`
	CreatedAt   string           `json:"created_at" example:"2025-05-11T08:00:00Z"`
	ClosedAt    sql.NullString   `json:"closed_at" example:"2025-05-11T09:00:00Z"`
	File        GpxInfoEssential `json:"file"`
}

type SessionDoc struct {
	ID          int    `json:"session_id" example:"1"`
	CreatedBy   string `json:"created_by" example:"550e8400-e29b-41d4-a716-446655440000"`
	Description string `json:"description" example:"Morning hike with friends"`
	CreatedAt   string `json:"created_at" example:"2025-05-11T08:00:00Z"`
	ClosedAt    struct {
		String string `json:"String" example:"2025-05-11T09:00:00Z"`
		Valid  bool   `json:"Valid" example:"true"`
	}
	File GpxInfoEssential `json:"file"`
}

type SessionInfo struct {
	SessionID   int     `json:"session_id" example:"1"`
	Description string  `json:"description" example:"Morning hike with friends"`
	Latitude    float64 `json:"latitude" example:"40.7128"`
	Longitude   float64 `json:"longitude" example:"-74.0060"`
	Altitude    float64 `json:"altitude" example:"10.5"`
	Accuracy    float64 `json:"accuracy" example:"5.0"`
	HelpRequest bool    `json:"help_request" example:"false"`
	GoingTo     string  `json:"going_to" example:"550e8400-e29b-41d4-a716-446655440000"`
	Time        string  `json:"time" example:"2025-05-11T08:00:00Z"`
	FileId      int     `json:"file_id" example:"1"`
}

type SessionInfoCreation struct {
	Description string  `json:"description" example:"Morning hike with friends"`
	Latitude    float64 `json:"latitude" example:"40.7128"`
	Longitude   float64 `json:"longitude" example:"-74.0060"`
	Altitude    float64 `json:"altitude" example:"10.5"`
	Accuracy    float64 `json:"accuracy" example:"5.0"`
	FileId      int     `json:"file_id" example:"1"`
}

type SessionId struct {
	ID string `json:"session_id" example:"1"`
}

type SessionInfoJoin struct {
	Latitude  float64 `json:"latitude" example:"40.7128"`
	Longitude float64 `json:"longitude" example:"-74.0060"`
	Altitude  float64 `json:"altitude" example:"10.5"`
	Accuracy  float64 `json:"accuracy" example:"5.0"`
}

type SessionInfoUpdate struct {
	Latitude      float64 `json:"latitude" example:"40.7128"`
	Longitude     float64 `json:"longitude" example:"-74.0060"`
	Altitude      float64 `json:"altitude" example:"10.5"`
	Accuracy      float64 `json:"accuracy" example:"5.0"`
	HelpRequested bool    `json:"help_request" example:"false"`
	GoingTo       string  `json:"going_to" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type MemberInfo struct {
	User        UserEssentials    `json:"user"`
	SessionInfo SessionInfoUpdate `json:"session_info"`
	TimeStamp   string            `json:"time_stamp" example:"2025-05-11T08:00:00Z"`
}
type SessionInfoResponse struct {
	CreatedBy   UserEssentials `json:"created_by"`
	Description string         `json:"description" example:"Morning hike with friends"`
	CreatedAt   string         `json:"created_at" example:"2025-05-11T08:00:00Z"`
	ClosedAt    sql.NullString `json:"closed_at" example:"2025-05-11T09:00:00Z"`
	Members     []MemberInfo   `json:"members"`
}

type SessionInfoResponseDoc struct {
	CreatedBy   UserEssentials `json:"created_by"`
	Description string         `json:"description" example:"Morning hike with friends"`
	CreatedAt   string         `json:"created_at" example:"2025-05-11T08:00:00Z"`
	ClosedAt    struct {
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

type Login struct {
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"strongPassword123"`
}

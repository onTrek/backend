package utils

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

type Group struct {
	ID          int              `json:"group_id" example:"1"`
	CreatedBy   string           `json:"created_by" example:"550e8400-e29b-41d4-a716-446655440000"`
	Description string           `json:"description" example:"Morning hike with friends"`
	CreatedAt   string           `json:"created_at" example:"2025-05-11T08:00:00Z"`
	File        GpxInfoEssential `json:"file"`
}

type GroupDoc struct {
	ID          int              `json:"group_id" example:"1"`
	CreatedBy   string           `json:"created_by" example:"550e8400-e29b-41d4-a716-446655440000"`
	Description string           `json:"description" example:"Morning hike with friends"`
	CreatedAt   string           `json:"created_at" example:"2025-05-11T08:00:00Z"`
	File        GpxInfoEssential `json:"file"`
}

type GroupInfo struct {
	GroupID     int     `json:"group_id" example:"1"`
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

type GroupInfoCreation struct {
	Description string `json:"description" example:"Morning hike with friends"`
	FileId      int    `json:"file_id" example:"1"`
}

type GroupId struct {
	ID string `json:"group_id" example:"1"`
}

type GroupInfoJoin struct {
	Latitude  float64 `json:"latitude" example:"40.7128"`
	Longitude float64 `json:"longitude" example:"-74.0060"`
	Altitude  float64 `json:"altitude" example:"10.5"`
	Accuracy  float64 `json:"accuracy" example:"5.0"`
}

type GroupInfoUpdate struct {
	Latitude      float64 `json:"latitude" example:"40.7128"`
	Longitude     float64 `json:"longitude" example:"-74.0060"`
	Altitude      float64 `json:"altitude" example:"10.5"`
	Accuracy      float64 `json:"accuracy" example:"5.0"`
	HelpRequested bool    `json:"help_request" example:"false"`
	GoingTo       string  `json:"going_to" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type MemberInfo struct {
	User          UserEssentials `json:"user"`
	Latitude      float64        `json:"latitude" example:"40.7128"`
	Longitude     float64        `json:"longitude" example:"-74.0060"`
	Altitude      float64        `json:"altitude" example:"10.5"`
	Accuracy      float64        `json:"accuracy" example:"5.0"`
	HelpRequested bool           `json:"help_request" example:"false"`
	GoingTo       string         `json:"going_to" example:"550e8400-e29b-41d4-a716-446655440000"`
	TimeStamp     string         `json:"time_stamp" example:"2025-05-11T08:00:00Z"`
}

type FileBody struct {
	FileId int `json:"file_id" example:"1"`
}

type GroupInfoResponse struct {
	CreatedBy   UserEssentials `json:"created_by"`
	Description string         `json:"description" example:"Morning hike with friends"`
	CreatedAt   string         `json:"created_at" example:"2025-05-11T08:00:00Z"`
	Members     []GroupMember  `json:"members"`
}

type GroupMember struct {
	ID       string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username string `json:"username" example:"John Doe"`
	Color    string `json:"color" example:"#e6194b"`
}

type GroupInfoResponseDoc struct {
	CreatedBy   UserEssentials `json:"created_by"`
	Description string         `json:"description" example:"Morning hike with friends"`
	CreatedAt   string         `json:"created_at" example:"2025-05-11T08:00:00Z"`
	Members     []GroupMember  `json:"members"`
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

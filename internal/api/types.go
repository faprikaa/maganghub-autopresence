package api

// DailyLogRequest represents the customizable fields for attendance submission
type DailyLogRequest struct {
	ActivityLog   string
	LessonLearned string
	Obstacles     string
}

// Attendance represents a single attendance record
type Attendance struct {
	ID             int    `json:"id"`
	ParticipantID  string `json:"participant_id"`
	Date           string `json:"date"`
	Status         string `json:"status"`
	ApprovalStatus string `json:"approval_status"`
	ReviewerID     string `json:"reviewer_id"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// AttendanceResponse represents the API response for attendance list
type AttendanceResponse struct {
	Data     []Attendance `json:"data"`
	Page     int          `json:"page"`
	PageSize int          `json:"page_size"`
	Total    int          `json:"total"`
}

// User represents user profile data
type User struct {
	ID                           string `json:"id"`
	Email                        string `json:"email"`
	PhoneNumber                  string `json:"phone_number"`
	Name                         string `json:"name"`
	UserType                     string `json:"user_type"`
	MentorName                   string `json:"mentor_name"`
	MentorPhoneNumber            string `json:"mentor_phone_number"`
	JobRole                      string `json:"job_role"`
	InternshipCompany            string `json:"internship_company"`
	InternshipCompanyAddress     string `json:"internship_company_address"`
	InternshipCompanyRegencyName string `json:"internship_company_regency_name"`
	InternshipCompanyPhone       string `json:"internship_company_phone"`
	InternshipCompanyEmail       string `json:"internship_company_email"`
	InternshipStartDate          string `json:"internship_start_date"`
	InternshipEndDate            string `json:"internship_end_date"`
	InternshipBatchName          string `json:"internship_batch_name"`
	CreatedAt                    string `json:"created_at"`
	UpdatedAt                    string `json:"updated_at"`
}

// UserResponse represents the API response for user profile
type UserResponse struct {
	Data User `json:"data"`
}

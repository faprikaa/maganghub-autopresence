package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

const MonevAPIURL = "https://monev.maganghub.kemnaker.go.id/api"

// Client handles API requests with authentication
type Client struct {
	httpClient    *http.Client
	cookies       []playwright.Cookie
	participantID string
}

// NewClient creates a new API client with cookies from playwright
func NewClient(cookies []playwright.Cookie) *Client {
	return &Client{
		httpClient: &http.Client{},
		cookies:    cookies,
	}
}

// formatCookies converts playwright cookies to a cookie header string
func (c *Client) formatCookies() string {
	var cookieParts []string
	for _, cookie := range c.cookies {
		cookieParts = append(cookieParts, fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))
	}
	return strings.Join(cookieParts, "; ")
}

// getAccessToken extracts access token from cookies
func (c *Client) getAccessToken() string {
	for _, cookie := range c.cookies {
		if cookie.Name == "accessToken" {
			return cookie.Value
		}
	}
	return ""
}

// setCommonHeaders sets common headers for API requests
func (c *Client) setCommonHeaders(req *http.Request) {
	req.Header.Set("accept", "application/json")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("referer", "https://monev.maganghub.kemnaker.go.id/dashboard")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36")
	req.Header.Set("cookie", c.formatCookies())

	// Add authorization header with access token
	if token := c.getAccessToken(); token != "" {
		req.Header.Set("authorization", "Bearer "+token)
	}
}

// GetMe fetches the current user profile and returns participant ID
func (c *Client) GetMe() (*User, error) {
	req, err := http.NewRequest("GET", MonevAPIURL+"/users/me", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setCommonHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var userResp UserResponse
	if err := json.Unmarshal(bodyText, &userResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Store participant ID for later use
	c.participantID = userResp.Data.ID

	return &userResp.Data, nil
}

// SubmitAttendance submits attendance with daily log
// Date is always today and Status is always PRESENT
func (c *Client) SubmitAttendance(req DailyLogRequest) (string, error) {
	// Date is always today
	today := time.Now().Format("2006-01-02")

	data := fmt.Sprintf(
		`{"date":"%s","status":"PRESENT","activity_log":"%s","lesson_learned":"%s","obstacles":"%s"}`,
		today, req.ActivityLog, req.LessonLearned, req.Obstacles,
	)

	httpReq, err := http.NewRequest("POST", MonevAPIURL+"/attendances/with-daily-log", strings.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	c.setCommonHeaders(httpReq)
	httpReq.Header.Set("content-type", "application/json")
	httpReq.Header.Set("origin", "https://monev.maganghub.kemnaker.go.id")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return string(bodyText), nil
}

// GetAttendances fetches attendance records for the current month
func (c *Client) GetAttendances() (*AttendanceResponse, error) {
	if c.participantID == "" {
		return nil, fmt.Errorf("participant ID not set, call GetMe() first")
	}

	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
	endDate := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location()).Format("2006-01-02")

	url := fmt.Sprintf("%s/attendances?participant_id=%s&start_date=%s&end_date=%s",
		MonevAPIURL, c.participantID, startDate, endDate)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setCommonHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var attendanceResp AttendanceResponse
	if err := json.Unmarshal(bodyText, &attendanceResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &attendanceResp, nil
}

// HasAttendedToday checks if attendance has been submitted for today
func (c *Client) HasAttendedToday() (bool, error) {
	attendances, err := c.GetAttendances()
	if err != nil {
		return false, err
	}

	today := time.Now().Format("2006-01-02")
	for _, att := range attendances.Data {
		if att.Date == today {
			return true, nil
		}
	}

	return false, nil
}

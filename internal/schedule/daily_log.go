package schedule

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// DailyLog represents a daily log entry for attendance
type DailyLog struct {
	Day           string `json:"day"`
	ActivityLog   string `json:"activity_log"`
	LessonLearned string `json:"lesson_learned"`
	Obstacles     string `json:"obstacles"`
}

// LoadDailyLogs loads daily logs from a JSON file
func LoadDailyLogs(filepath string) ([]DailyLog, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var logs []DailyLog
	if err := json.Unmarshal(data, &logs); err != nil {
		return nil, err
	}

	return logs, nil
}

// GetTodayLog returns the log entry for today's day of week
func GetTodayLog(logs []DailyLog) *DailyLog {
	today := time.Now().Weekday().String()

	var fallback *DailyLog
	for i := range logs {
		if logs[i].Day == today {
			return &logs[i]
		}
		if logs[i].Day == "Fallback" {
			fallback = &logs[i]
		}
	}

	// Return fallback entry if no match
	if fallback != nil {
		return fallback
	}

	// Last resort: return first entry
	if len(logs) > 0 {
		return &logs[0]
	}

	return nil
}

// ValidateDailyLogs checks if all log fields have at least 100 characters
// Returns warnings for any fields that are too short
func ValidateDailyLogs(logs []DailyLog) []string {
	var warnings []string
	minLength := 100

	for _, log := range logs {
		if len(log.ActivityLog) < minLength {
			warnings = append(warnings, fmt.Sprintf("[%s] activity_log terlalu pendek (%d karakter, minimal %d)", log.Day, len(log.ActivityLog), minLength))
		}
		if len(log.LessonLearned) < minLength {
			warnings = append(warnings, fmt.Sprintf("[%s] lesson_learned terlalu pendek (%d karakter, minimal %d)", log.Day, len(log.LessonLearned), minLength))
		}
		if len(log.Obstacles) < minLength {
			warnings = append(warnings, fmt.Sprintf("[%s] obstacles terlalu pendek (%d karakter, minimal %d)", log.Day, len(log.Obstacles), minLength))
		}
	}

	return warnings
}

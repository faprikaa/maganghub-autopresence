package schedule

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"
)

// DailyLog represents a daily log entry for attendance
type DailyLog struct {
	Day           string   `json:"day"`
	ActivityLog   []string `json:"activity_log"`
	LessonLearned []string `json:"lesson_learned"`
	Obstacles     []string `json:"obstacles"`
}

// ResolvedDailyLog is the selected daily log content ready for submission
type ResolvedDailyLog struct {
	Day           string
	ActivityLog   string
	LessonLearned string
	Obstacles     string
}

type dailyLogsByDay struct {
	Days map[string]dailyLogEntry `json:"days"`
}

type dailyLogEntry struct {
	ActivityLog   []string `json:"activity_log"`
	LessonLearned []string `json:"lesson_learned"`
	Obstacles     []string `json:"obstacles"`
}

// LoadDailyLogs loads daily logs from a JSON file
func LoadDailyLogs(filepath string) ([]DailyLog, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var byDay dailyLogsByDay
	if err := json.Unmarshal(data, &byDay); err != nil {
		return nil, fmt.Errorf("format daily_logs.json tidak valid: %w", err)
	}
	if len(byDay.Days) == 0 {
		return nil, fmt.Errorf("format daily_logs.json tidak valid: field days wajib ada dan tidak boleh kosong")
	}

	var dayKeys []string
	for day := range byDay.Days {
		dayKeys = append(dayKeys, day)
	}
	sort.Strings(dayKeys)

	logs := make([]DailyLog, 0, len(dayKeys))
	for _, day := range dayKeys {
		entry := byDay.Days[day]
		logs = append(logs, DailyLog{
			Day:           day,
			ActivityLog:   entry.ActivityLog,
			LessonLearned: entry.LessonLearned,
			Obstacles:     entry.Obstacles,
		})
	}

	return logs, nil
}

// GetTodayLog returns a randomly selected log entry for today's day of week
func GetTodayLog(logs []DailyLog) *ResolvedDailyLog {
	today := time.Now().Weekday().String()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var candidates []DailyLog
	var fallback []DailyLog
	for i := range logs {
		if logs[i].Day == today {
			candidates = append(candidates, logs[i])
		}
		if logs[i].Day == "Fallback" {
			fallback = append(fallback, logs[i])
		}
	}

	if len(candidates) == 0 {
		candidates = fallback
	}

	if len(candidates) == 0 && len(logs) > 0 {
		candidates = logs
	}

	if len(candidates) == 0 {
		return nil
	}

	selected := candidates[rng.Intn(len(candidates))]

	return &ResolvedDailyLog{
		Day:           selected.Day,
		ActivityLog:   pickRandom(rng, selected.ActivityLog),
		LessonLearned: pickRandom(rng, selected.LessonLearned),
		Obstacles:     pickRandom(rng, selected.Obstacles),
	}
}

func pickRandom(rng *rand.Rand, values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[rng.Intn(len(values))]
}

// ValidateDailyLogs checks if all log fields have at least 100 characters
// Returns warnings for any fields that are too short
func ValidateDailyLogs(logs []DailyLog) []string {
	var warnings []string
	minLength := 100

	for _, log := range logs {
		if len(log.ActivityLog) == 0 {
			warnings = append(warnings, fmt.Sprintf("[%s] activity_log kosong", log.Day))
		}
		for i, value := range log.ActivityLog {
			if len(value) < minLength {
				warnings = append(warnings, fmt.Sprintf("[%s] activity_log[%d] terlalu pendek (%d karakter, minimal %d)", log.Day, i, len(value), minLength))
			}
		}

		if len(log.LessonLearned) == 0 {
			warnings = append(warnings, fmt.Sprintf("[%s] lesson_learned kosong", log.Day))
		}
		for i, value := range log.LessonLearned {
			if len(value) < minLength {
				warnings = append(warnings, fmt.Sprintf("[%s] lesson_learned[%d] terlalu pendek (%d karakter, minimal %d)", log.Day, i, len(value), minLength))
			}
		}

		if len(log.Obstacles) == 0 {
			warnings = append(warnings, fmt.Sprintf("[%s] obstacles kosong", log.Day))
		}
		for i, value := range log.Obstacles {
			if len(value) < minLength {
				warnings = append(warnings, fmt.Sprintf("[%s] obstacles[%d] terlalu pendek (%d karakter, minimal %d)", log.Day, i, len(value), minLength))
			}
		}
	}

	return warnings
}

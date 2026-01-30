package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"maganghub-autopresence/internal/api"
	"maganghub-autopresence/internal/browser"
	"maganghub-autopresence/internal/config"
	"maganghub-autopresence/internal/cookie_manager"
	"maganghub-autopresence/internal/schedule"
)

func main() {
	cfg := config.Load()

	// Load and validate daily logs at startup
	logs, err := schedule.LoadDailyLogs("daily_logs.json")
	if err != nil {
		log.Printf("⚠️  Warning: Failed to load daily_logs.json: %v", err)
	} else {
		warnings := schedule.ValidateDailyLogs(logs)
		if len(warnings) > 0 {
			log.Println("⚠️  WARNING: Beberapa field daily log terlalu pendek!")
			for _, w := range warnings {
				log.Printf("   - %s", w)
			}
		} else {
			log.Printf("✅ Daily logs validated: %d entries OK", len(logs))
		}
	}

	// Verify login at startup
	log.Println("🔐 Verifying login credentials...")
	client, err := browser.NewBrowserClient(cfg.Headless)
	if err != nil {
		log.Fatalf("❌ Failed to create browser: %v", err)
	}
	userName, err := client.Login(cfg.MaganghubConfig.Username, cfg.MaganghubConfig.Password)
	if err != nil {
		client.Close()
		log.Fatalf("❌ Login failed: %v", err)
	}
	log.Printf("✅ Login verified as: %s", userName)
	client.Close()

	// Run attendance check on startup
	log.Println("📋 Checking attendance for today...")
	runAttendance(cfg)

	// Create scheduler with cron from config
	scheduler := schedule.NewScheduler(cfg.CronSchedule)

	// Define the attendance job
	attendanceJob := func() {
		runAttendance(cfg)
	}

	// Start the scheduler
	if err := scheduler.Start(attendanceJob); err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}

	log.Printf("Next scheduled run: %s", scheduler.GetNextRun())

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down scheduler...")
	scheduler.Stop()
}

// Global cookie manager
var cookieManager = cookie_manager.NewCookieManager("cookies.json")

func runAttendance(cfg *config.Config) {
	// Load daily logs from JSON
	logs, err := schedule.LoadDailyLogs("daily_logs.json")
	if err != nil {
		log.Printf("Failed to load daily logs: %v", err)
		return
	}

	// Get today's log entry
	todayLog := schedule.GetTodayLog(logs)
	if todayLog == nil {
		log.Println("No log entry found for today")
		return
	}

	log.Printf("Using log for: %s", todayLog.Day)

	// Try to use cached cookies first
	var apiClient *api.Client
	var cachedUser *api.User
	if cookieManager.HasCookies() {
		log.Println("🔄 Trying with cached cookies...")
		apiClient = api.NewClient(cookieManager.Get())

		// Test if cookies are still valid by calling GetMe
		cachedUser, err = apiClient.GetMe()
		if err != nil || cachedUser.ID == "" {
			log.Println("⚠️  Cached cookies expired, re-logging in...")
			apiClient = nil
			cachedUser = nil
		} else {
			log.Println("✅ Cached cookies still valid")
		}
	}

	// If no valid cookies, login fresh
	if apiClient == nil {
		log.Println("🔐 Logging in...")
		client, err := browser.NewBrowserClient(cfg.Headless)
		if err != nil {
			log.Printf("Browser error: %v", err)
			return
		}

		userName, err := client.Login(cfg.MaganghubConfig.Username, cfg.MaganghubConfig.Password)
		if err != nil {
			client.Close()
			log.Printf("Login error: %v", err)
			return
		}
		log.Printf("✅ Logged in as: %s", userName)

		cookies, err := client.GetCookies()
		client.Close()
		if err != nil {
			log.Printf("Cookie error: %v", err)
			return
		}

		// Save cookies for future use
		if err := cookieManager.Save(cookies); err != nil {
			log.Printf("Warning: Failed to save cookies: %v", err)
		}

		apiClient = api.NewClient(cookies)
	}

	// Get user profile (reuse cached user if available)
	var user *api.User
	if cachedUser != nil {
		user = cachedUser
	} else {
		user, err = apiClient.GetMe()
		if err != nil {
			log.Printf("GetMe error: %v", err)
			cookieManager.Clear() // Clear invalid cookies
			return
		}
	}
	log.Printf("User: %s (%s)", user.Name, user.ID)

	// Check if already attended today
	hasAttended, err := apiClient.HasAttendedToday()
	if err != nil {
		log.Printf("HasAttendedToday error: %v", err)
		return
	}

	if hasAttended {
		log.Println("Already attended today, skipping submission")
		return
	}

	// Submit attendance with today's log
	response, err := apiClient.SubmitAttendance(api.DailyLogRequest{
		ActivityLog:   todayLog.ActivityLog,
		LessonLearned: todayLog.LessonLearned,
		Obstacles:     todayLog.Obstacles,
	})
	if err != nil {
		log.Printf("SubmitAttendance error: %v", err)
		return
	}

	log.Printf("Attendance submitted: %s", response)
}

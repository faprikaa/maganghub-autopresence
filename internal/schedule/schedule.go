package schedule

import (
	"log"

	"github.com/robfig/cron/v3"
)

// Scheduler wraps the cron scheduler
type Scheduler struct {
	cron           *cron.Cron
	cronExpression string
}

// NewScheduler creates a new scheduler with the given cron expression
func NewScheduler(cronExpression string) *Scheduler {
	return &Scheduler{
		cron:           cron.New(),
		cronExpression: cronExpression,
	}
}

// Start starts the scheduler with the given job function
func (s *Scheduler) Start(job func()) error {
	_, err := s.cron.AddFunc(s.cronExpression, func() {
		log.Println("Running scheduled job...")
		job()
		log.Println("Scheduled job completed")
	})
	if err != nil {
		return err
	}

	log.Printf("Scheduler started with cron expression: %s", s.cronExpression)
	s.cron.Start()
	return nil
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	s.cron.Stop()
	log.Println("Scheduler stopped")
}

// GetNextRun returns the next scheduled run time as a string
func (s *Scheduler) GetNextRun() string {
	entries := s.cron.Entries()
	if len(entries) > 0 {
		return entries[0].Next.Format("2006-01-02 15:04:05")
	}
	return "No scheduled job"
}

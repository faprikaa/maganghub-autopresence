package cookie_manager

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/playwright-community/playwright-go"
)

// CookieManager handles cookie persistence and caching
type CookieManager struct {
	mu       sync.RWMutex
	cookies  []playwright.Cookie
	filepath string
}

// NewCookieManager creates a new cookie manager
func NewCookieManager(filepath string) *CookieManager {
	cm := &CookieManager{filepath: filepath}
	cm.load() // Try to load existing cookies
	return cm
}

// Save saves cookies to file and memory
func (cm *CookieManager) Save(cookies []playwright.Cookie) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.cookies = cookies

	data, err := json.MarshalIndent(cookies, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cm.filepath, data, 0600)
}

// Get returns cached cookies
func (cm *CookieManager) Get() []playwright.Cookie {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.cookies
}

// HasCookies returns true if cookies are available
func (cm *CookieManager) HasCookies() bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return len(cm.cookies) > 0
}

// load loads cookies from file
func (cm *CookieManager) load() {
	data, err := os.ReadFile(cm.filepath)
	if err != nil {
		return
	}

	var cookies []playwright.Cookie
	if err := json.Unmarshal(data, &cookies); err != nil {
		return
	}

	cm.cookies = cookies
}

// Clear clears cached cookies
func (cm *CookieManager) Clear() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.cookies = nil
	os.Remove(cm.filepath)
}

package browser

import (
	"errors"
	"fmt"

	"github.com/playwright-community/playwright-go"
)

const (
	MaganghubAuthURL  = "https://account.kemnaker.go.id/auth/login"
	MonevDashboardURL = "https://monev.maganghub.kemnaker.go.id/dashboard"
)

// BrowserClient encapsulates the browser state
type BrowserClient struct {
	pw      *playwright.Playwright
	Browser playwright.Browser
	Context playwright.BrowserContext
	Page    playwright.Page
}

// NewBrowserClient initializes a new browser client
func NewBrowserClient(headless bool) (*BrowserClient, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("could not start playwright: %w", err)
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(headless),
	})
	if err != nil {
		return nil, fmt.Errorf("could not launch browser: %w", err)
	}

	return &BrowserClient{
		pw:      pw,
		Browser: browser,
	}, nil
}

// Login performs authentication and returns the logged-in user name
func (c *BrowserClient) Login(username, password string) (string, error) {
	page, err := c.Browser.NewPage()
	if err != nil {
		return "", fmt.Errorf("could not create page: %w", err)
	}
	c.Page = page
	c.Context = c.Browser.Contexts()[0]

	if _, err = page.Goto(MonevDashboardURL); err != nil {
		return "", fmt.Errorf("could not navigate: %w", err)
	}

	if err := page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateNetworkidle,
	}); err != nil {
		return "", fmt.Errorf("could not wait for load: %w", err)
	}

	// Check if redirected to auth page
	if page.URL() != MaganghubAuthURL {
		return "", errors.New("not redirected to auth page")
	}

	// Fill credentials
	if err := page.Locator("#username").Fill(username); err != nil {
		return "", fmt.Errorf("could not fill username: %w", err)
	}
	if err := page.Locator("#password").Fill(password); err != nil {
		return "", fmt.Errorf("could not fill password: %w", err)
	}
	if err := page.Locator("button[type='submit']").Click(); err != nil {
		return "", fmt.Errorf("could not click submit: %w", err)
	}

	if err := page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateNetworkidle,
	}); err != nil {
		return "", fmt.Errorf("could not wait for login: %w", err)
	}

	// Get user name using XPath selector
	nameText, err := page.Locator("//*[@id='__nuxt']/div/div/div/div/div/div/div[2]/div/div[1]/div[2]/div[1]/div/div[2]/div/div[1]/div[2]").First().TextContent()
	if err != nil {
		return "", fmt.Errorf("could not get name text: %w", err)
	}

	return nameText, nil
}

// GetCookies returns cookies from the current context
func (c *BrowserClient) GetCookies() ([]playwright.Cookie, error) {
	if c.Context == nil {
		return nil, errors.New("no browser context available")
	}
	return c.Context.Cookies()
}

// Close cleans up all browser resources
func (c *BrowserClient) Close() error {
	if c.Browser != nil {
		if err := c.Browser.Close(); err != nil {
			return err
		}
	}
	if c.pw != nil {
		if err := c.pw.Stop(); err != nil {
			return err
		}
	}
	return nil
}

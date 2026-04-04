package mailer

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"net/url"
	"strings"
	"time"

	"github.com/ima/diplom-backend/internal/pkg/logger"
)

type Mailer interface {
	SendOTP(ctx context.Context, toEmail, code string) error
}

type Config struct {
	APIKey    string
	APIURL    string
	FromEmail string
	FromName  string
	UploadDir string
}

type unisenderMailer struct {
	cfg    Config
	client *http.Client
}

func New(cfg Config) Mailer {
	return &unisenderMailer{
		cfg: cfg,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (m *unisenderMailer) SendOTP(ctx context.Context, toEmail, code string) error {
	log := logger.FromContext(ctx)

	// If API key is not set, just log it out (useful in development)
	if m.cfg.APIKey == "" {
		log.Info("UniSender API key not set, simulating email dispatch", "to", toEmail, "code", code)
		return nil
	}

	apiURL := m.cfg.APIURL
	if apiURL == "" {
		apiURL = "https://api.unisender.com/ru/api/sendEmail"
	}

	body := "Code: " + code + "\nExpires in 10 minutes."
	subject := "Your login code"

	params := url.Values{}
	params.Set("format", "json")
	params.Set("api_key", m.cfg.APIKey)
	params.Set("email", toEmail)
	params.Set("sender_name", m.cfg.FromName)
	params.Set("sender_email", m.cfg.FromEmail)
	params.Set("subject", subject)
	params.Set("body", body)
	params.Set("list_id", "1") // Assuming some list ID or not needed by API docs

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, strings.NewReader(params.Encode()))
	if err != nil {
		return errors.New("failed to create request: " + err.Error())
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := m.client.Do(req)
	if err != nil {
		return errors.New("failed to send email via UniSender: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return errors.New("unisender returned status " + strconv.Itoa(resp.StatusCode))
	}

	var jsonResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&jsonResp); err != nil {
		return errors.New("failed to decode unisender response: " + err.Error())
	}

	if msg, hasErr := jsonResp["error"]; hasErr {
		strMsg, _ := msg.(string)
		if strMsg == "" {
			strMsg = "unknown error"
		}
		return errors.New("unisender error: " + strMsg)
	}

	log.Info("OTP email sent successfully", "to", toEmail)
	return nil
}

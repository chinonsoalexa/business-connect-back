package email

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"strings"
	"time"
)

var disposableDomains = map[string]bool{}

// ------------------------------------------------------------
// 1. Load disposable email list (JSON array of strings)
// ------------------------------------------------------------
func LoadDisposableList(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var list []string
	decoder := json.NewDecoder(bufio.NewReader(file))

	err = decoder.Decode(&list)
	if err != nil {
		return err
	}

	for _, d := range list {
		disposableDomains[strings.ToLower(d)] = true
	}
	return nil
}

// ------------------------------------------------------------
// 2. Normalize Email (Gmail rules)
// ------------------------------------------------------------
func NormalizeEmail(email string) string {
	email = strings.TrimSpace(strings.ToLower(email))
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}

	local, domain := parts[0], parts[1]

	if domain == "gmail.com" || domain == "googlemail.com" {
		local = strings.Split(local, "+")[0]
		local = strings.ReplaceAll(local, ".", "")
	}

	return local + "@" + domain
}

// ------------------------------------------------------------
// 3. Check Disposable Domain
// ------------------------------------------------------------
func IsDisposable(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	domain := parts[1]
	return disposableDomains[domain]
}

// ------------------------------------------------------------
// 4. MX Check
// ------------------------------------------------------------
func HasMX(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	_, err := net.LookupMX(parts[1])
	return err == nil
}

// ------------------------------------------------------------
// 5. Catch-all Detection (Accept-all)
// ------------------------------------------------------------
func IsCatchAll(domain string) bool {
	testEmail := "does-not-exist-" + fmt.Sprint(time.Now().Unix()) + "@" + domain

	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		return false
	}

	// Try SMTP check
	mxHost := mxRecords[0].Host
	conn, err := net.DialTimeout("tcp", mxHost+":25", 3*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()

	// Fake SMTP session
	c, err := smtp.NewClient(conn, mxHost)
	if err != nil {
		return false
	}
	defer c.Close()

	c.Mail("validator@" + domain)
	err = c.Rcpt(testEmail)

	// If RCPT does NOT return error → catch-all domain
	return err == nil
}

// ------------------------------------------------------------
// 6. Username Pattern Risk Score
// ------------------------------------------------------------
func UsernameRiskScore(email string) int {
	local := strings.Split(email, "@")[0]
	localLower := strings.ToLower(local)

	score := 0

	// 1. Very short usernames → usually bots
	if len(local) < 3 {
		score += 40
	}

	// 2. Too many numbers → suspicious
	if percentNumbers(local) > 0.5 {
		score += 30
	}

	// 3. Repeated patterns (aaa333, qqq111)
	if hasRepetitivePattern(localLower) {
		score += 20
	}

	// 4. Too many random consonants in a row (e.g. "xqtrpld")
	if longConsonantRun(localLower) {
		score += 35
	}

	return score
}

func percentNumbers(s string) float64 {
	n := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n++
		}
	}
	return float64(n) / float64(len(s))
}

func hasRepetitivePattern(s string) bool {
	for i := 0; i < len(s)-2; i++ {
		if s[i] == s[i+1] && s[i] == s[i+2] {
			return true
		}
	}
	return false
}

func longConsonantRun(s string) bool {
	consonants := "bcdfghjklmnpqrstvwxyz"
	run := 0

	for _, c := range s {
		if strings.ContainsRune(consonants, c) {
			run++
			if run >= 5 {
				return true
			}
		} else {
			run = 0
		}
	}
	return false
}

// ------------------------------------------------------------
// Final Combined Result
// ------------------------------------------------------------

type EmailCheckResult struct {
	Normalized     string
	IsDisposable   bool
	HasMX          bool
	IsCatchAll     bool
	RiskScore      int
}

func ValidateEmail(email string) (EmailCheckResult, error) {
	if !strings.Contains(email, "@") {
		return EmailCheckResult{}, errors.New("invalid email format")
	}

	normalized := NormalizeEmail(email)
	parts := strings.Split(normalized, "@")
	domain := parts[1]

	result := EmailCheckResult{
		Normalized: normalized,
		IsDisposable: IsDisposable(normalized),
		HasMX:        HasMX(normalized),
		IsCatchAll:   IsCatchAll(domain),
	}

	// Combine scoring
	if result.IsDisposable {
		result.RiskScore += 80
	}
	if !result.HasMX {
		result.RiskScore += 80
	}
	if result.IsCatchAll {
		result.RiskScore += 40
	}

	// Add username scoring
	result.RiskScore += UsernameRiskScore(normalized)

	if result.RiskScore > 100 {
		result.RiskScore = 100
	}

	return result, nil
}

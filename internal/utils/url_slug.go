package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// GenerateEnterpriseTrialURLSlug generates a URL slug for enterprise free trial subscription
// using the provided domain name
func GenerateEnterpriseTrialURLSlug(domain string, companyName string, planType string) (string, error) {
	// Validate and clean domain
	domain = strings.TrimSpace(domain)
	if domain == "" {
		return "", fmt.Errorf("domain cannot be empty")
	}

	// Ensure domain has proper format
	if !strings.HasPrefix(domain, "http://") && !strings.HasPrefix(domain, "https://") {
		domain = "https://" + domain
	}

	// Parse domain URL
	parsedURL, err := url.Parse(domain)
	if err != nil {
		return "", fmt.Errorf("invalid domain format: %v", err)
	}

	// Clean company name and plan type
	companyName = strings.TrimSpace(companyName)
	planType = strings.TrimSpace(planType)

	// Create slug from company name
	slug := strings.ToLower(companyName)

	// Replace special characters with dashes using regex
	reg := regexp.MustCompile(`[^\w\s-]`)
	slug = reg.ReplaceAllString(slug, "")

	// Replace spaces and multiple dashes with single dash
	slug = strings.ReplaceAll(slug, " ", "-")
	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Remove leading and trailing dashes
	slug = strings.Trim(slug, "-")

	// Add plan type if provided
	if planType != "" {
		planSlug := strings.ToLower(planType)
		planSlug = reg.ReplaceAllString(planSlug, "")
		planSlug = strings.ReplaceAll(planSlug, " ", "-")
		planSlug = reg.ReplaceAllString(planSlug, "-")
		planSlug = strings.Trim(planSlug, "-")

		if planSlug != "" {
			slug = slug + "-" + planSlug + "-trial"
		}
	} else {
		slug = slug + "-enterprise-trial"
	}

	// Construct final URL
	finalURL := fmt.Sprintf("%s/enterprise/%s", parsedURL.String(), slug)

	return finalURL, nil
}

// GenerateCustomTrialURL generates a custom trial URL with specific parameters
func GenerateCustomTrialURL(domain string, params map[string]string) (string, error) {
	// Validate and clean domain
	domain = strings.TrimSpace(domain)
	if domain == "" {
		return "", fmt.Errorf("domain cannot be empty")
	}

	// Ensure domain has proper format
	if !strings.HasPrefix(domain, "http://") && !strings.HasPrefix(domain, "https://") {
		domain = "https://" + domain
	}

	// Parse domain URL
	parsedURL, err := url.Parse(domain)
	if err != nil {
		return "", fmt.Errorf("invalid domain format: %v", err)
	}

	// Build query parameters
	query := url.Values{}
	for key, value := range params {
		query.Add(key, value)
	}

	// Construct final URL with query parameters
	finalURL := fmt.Sprintf("%s/enterprise/trial?%s", parsedURL.String(), query.Encode())

	return finalURL, nil
}

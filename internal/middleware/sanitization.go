package middleware

import (
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"

	"holidayapi/internal/models"
)

// InputSanitizationMiddleware sanitizes request inputs
func InputSanitizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only sanitize for POST, PUT, PATCH requests with JSON content
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			contentType := c.GetHeader("Content-Type")
			if strings.Contains(contentType, "application/json") {
				// Read the body
				body, err := io.ReadAll(c.Request.Body)
				if err != nil {
					c.JSON(http.StatusBadRequest, models.ErrorResponse{
						Success: false,
						Message: "Failed to read request body",
						Error:   err.Error(),
					})
					c.Abort()
					return
				}

				// Sanitize the body
				sanitizedBody := sanitizeJSON(string(body))

				// Replace the body
				c.Request.Body = io.NopCloser(strings.NewReader(sanitizedBody))
				c.Request.ContentLength = int64(len(sanitizedBody))
			}
		}

		// Sanitize query parameters
		sanitizeQueryParams(c)

		c.Next()
	}
}

// sanitizeJSON sanitizes JSON string content
func sanitizeJSON(jsonStr string) string {
	// Only remove extremely dangerous patterns, don't break valid JSON
	dangerousPatterns := []string{
		`<script[^>]*>.*?</script>`,
		`javascript:`,
		`vbscript:`,
		`onload\s*=`,
		`onerror\s*=`,
		`onclick\s*=`,
	}

	for _, pattern := range dangerousPatterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		jsonStr = re.ReplaceAllString(jsonStr, "")
	}

	return jsonStr
}

// sanitizeQueryParams sanitizes query parameters
func sanitizeQueryParams(c *gin.Context) {
	query := c.Request.URL.Query()

	for key, values := range query {
		for i, value := range values {
			// Basic sanitization - only remove extremely dangerous content
			value = strings.TrimSpace(value)

			// Remove script tags and javascript
			scriptPattern := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
			value = scriptPattern.ReplaceAllString(value, "")

			jsPattern := regexp.MustCompile(`(?i)javascript:`)
			value = jsPattern.ReplaceAllString(value, "")

			values[i] = value
		}
		query[key] = values
	}

	c.Request.URL.RawQuery = query.Encode()
}

// SQLInjectionProtectionMiddleware provides additional SQL injection protection
func SQLInjectionProtectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for common SQL injection patterns in URL path
		path := strings.ToLower(c.Request.URL.Path)

		sqlInjectionPatterns := []string{
			"union", "select", "insert", "update", "delete", "drop",
			"create", "alter", "exec", "execute", "script", "--", "/*", "*/",
			"xp_", "sp_", "0x", "char(", "ascii(", "substring(",
		}

		for _, pattern := range sqlInjectionPatterns {
			if strings.Contains(path, pattern) {
				c.JSON(http.StatusBadRequest, models.ErrorResponse{
					Success: false,
					Message: "Invalid request",
					Error:   "Request contains potentially dangerous content",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// XSSProtectionMiddleware provides XSS protection
func XSSProtectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for XSS patterns in headers
		userAgent := strings.ToLower(c.GetHeader("User-Agent"))
		referer := strings.ToLower(c.GetHeader("Referer"))

		xssPatterns := []string{
			"<script", "javascript:", "vbscript:", "onload=", "onerror=",
			"onclick=", "onmouseover=", "onfocus=", "onblur=", "onchange=",
			"onsubmit=", "onreset=", "onselect=", "onunload=",
		}

		for _, pattern := range xssPatterns {
			if strings.Contains(userAgent, pattern) || strings.Contains(referer, pattern) {
				c.JSON(http.StatusBadRequest, models.ErrorResponse{
					Success: false,
					Message: "Invalid request",
					Error:   "Request headers contain potentially dangerous content",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

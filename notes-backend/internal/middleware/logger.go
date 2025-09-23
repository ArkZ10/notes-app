package middleware

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Clone request body
		var reqBody []byte
		if c.Request.Body != nil {
			reqBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}

		// Mask Authorization header
		headers := map[string]string{}
		for k, v := range c.Request.Header {
			if k == "Authorization" {
				headers[k] = "*****"
			} else {
				headers[k] = v[0]
			}
		}

		// Capture response body
		writer := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = writer

		c.Next() // process request

		// Calculate duration
		duration := time.Since(start)

		go func() {
			_, err := db.Exec(
				`INSERT INTO logs (method, endpoint, request_headers, request_body, response_body, status_code, created_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7)`,
				c.Request.Method,
				c.Request.URL.Path,
				headersToJSON(headers),
				bytesToJSON(reqBody),
				bytesToJSON(writer.body.Bytes()),
				c.Writer.Status(),
				time.Now(),
			)
			if err != nil {
				fmt.Printf("❌ Failed to insert log: %v\n", err)
			}
		}()

		// Console log for dev
		fmt.Printf("Request %s %s took %v\n", c.Request.Method, c.Request.URL.Path, duration)
	}
}

// Helper to convert map to JSON
func headersToJSON(h map[string]string) []byte {
	b, _ := json.Marshal(h)
	return b
}

func bytesToJSON(b []byte) []byte {
	if len(b) == 0 {
		return []byte("null")
	}

	// Ensure valid UTF-8
	if !json.Valid(b) {
		// If it's not valid JSON, wrap it as a string
		safe := string(b)
		sanitized, _ := json.Marshal(safe)
		return sanitized
	}

	return b
}

// bodyWriter captures response body
type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

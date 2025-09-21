package logs

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListLogsHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(`
			SELECT id, method, endpoint, request_headers, request_body, response_body, status_code, created_at
			FROM logs
			ORDER BY created_at DESC
			LIMIT 50
		`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch logs"})
			return
		}
		defer rows.Close()

		type Log struct {
			ID             int64       `json:"id"`
			Method         string      `json:"method"`
			Endpoint       string      `json:"endpoint"`
			RequestHeaders interface{} `json:"request_headers"`
			RequestBody    interface{} `json:"request_body"`
			ResponseBody   interface{} `json:"response_body"`
			StatusCode     int         `json:"status_code"`
			CreatedAt      string      `json:"created_at"`
		}

		var logs []Log
		for rows.Next() {
			var l Log
			var reqHeaders, reqBody, resBody []byte
			if err := rows.Scan(
				&l.ID, &l.Method, &l.Endpoint,
				&reqHeaders, &reqBody, &resBody,
				&l.StatusCode, &l.CreatedAt,
			); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse log"})
				return
			}
			l.RequestHeaders = jsonBytesToInterface(reqHeaders)
			l.RequestBody = jsonBytesToInterface(reqBody)
			l.ResponseBody = jsonBytesToInterface(resBody)
			logs = append(logs, l)
		}

		c.JSON(http.StatusOK, logs)
	}
}

func jsonBytesToInterface(b []byte) interface{} {
	if len(b) == 0 || string(b) == "null" {
		return nil
	}
	return string(b) // return as string for simplicity
}

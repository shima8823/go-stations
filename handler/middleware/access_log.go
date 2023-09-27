package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type accessLogInfo struct {
	Timestamp time.Time
	Latency   int64
	Path      string
	OS        string
}

func AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)

		osname, ok := r.Context().Value(OSNameKey).(string)
		if !ok {
			osname = "unknown"
		}
		accessLog := accessLogInfo{
			Timestamp: startTime,
			Latency:   time.Now().UnixNano() - startTime.UnixNano(),
			Path:      r.URL.Path,
			OS:        osname,
		}
		accessLogJSON, err := json.Marshal(accessLog)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(accessLogJSON))
	})
}

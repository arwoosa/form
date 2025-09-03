package ezgrpc

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func formToJSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// 把 form 轉成 map
			data := make(map[string]string)
			for k, v := range r.Form {
				if len(v) > 0 {
					data[k] = v[0]
				}
			}

			// encode 成 JSON，重寫 Body
			buf, _ := json.Marshal(data)
			r.Body = io.NopCloser(bytes.NewReader(buf))
			r.ContentLength = int64(len(buf))
			r.Header.Set("Content-Type", "application/json")
		}
		next.ServeHTTP(w, r)
	})
}

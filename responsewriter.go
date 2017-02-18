package looli

import (
    "net/http"
)

type ResponseWriter struct {
    http.ResponseWriter
}

func (rw *ResponseWriter) WriteString(statusCode int, response string) {

}

func (rw *ResponseWriter) CloseNotify() <-chan bool {
    return rw.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

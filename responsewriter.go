package looli

import (
    "net/http"
)

type ResponseWriter struct {
    http.ResponseWriter
}

func (rw *ResponseWriter) String(statusCode int, response string) {
	
}

func (rw *ResponseWriter) JSON(statusCode int, response interface{}) {

}
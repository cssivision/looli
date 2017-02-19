package looli

import (
    "fmt"
    "io"
    "time"
    "os"
)

var defaultWriter = os.Stdout

func Logger() HandlerFunc {
    return LoggerWithWriter(defaultWriter)
}

func LoggerWithWriter(out io.Writer) HandlerFunc {
    return func(c *Context) {
        start := time.Now()
        path := c.Path
        method := c.Request.Method
        c.Next()

        end := time.Now()
        latency := end.Sub(start)
        clientIP := c.ClientIP()
        statusCode := c.StatusCode

        comment := c.ErrorMessage
        fmt.Fprintf(out, "[LOOLI] %v | %3d | %13v | %s | %-7s %s\n%s",
                end.Format("2017/02/19 - 20:23:15"),
                statusCode,
                latency,
                clientIP,
                method,
                path,
                comment,
            )
    }
}
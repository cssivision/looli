package looli

import (
    "io"
    "os"
)

var defaultErrorWriter = os.Stderr

func Recover() HandlerFunc {
    return RecoverWithWriter(defaultErrorWriter)
}

func RecoverWithWriter(out io.Writer) HandlerFunc {
    return func(c *Context) {
        defer func() {
            if err := recover(); err != nil {
                c.AbortWithStatus(500)
            }
        }()
        c.Next()
    }
}
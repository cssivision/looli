package looli

import (
	"fmt"
	"io"
	"os"
	"runtime"
)

var defaultErrorWriter = os.Stderr

func Recover() HandlerFunc {
	return RecoverWithWriter(defaultErrorWriter)
}

func RecoverWithWriter(out io.Writer) HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 2048)
				buf = buf[:runtime.Stack(buf, false)]
				fmt.Printf("[Recover] panic recovered:\n%s\n%s\n", string(buf), err)

				c.AbortWithStatus(500)
				return
			}
		}()
		c.Next()
	}
}

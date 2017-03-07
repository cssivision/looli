package looli

import (
	"io"
	"log"
	"net/http/httputil"
	"os"
)

var defaultErrorWriter = os.Stderr

func Recover() HandlerFunc {
	return RecoverWithWriter(defaultErrorWriter)
}

func RecoverWithWriter(out io.Writer) HandlerFunc {
	var logger *log.Logger
	if out != nil {
		logger = log.New(out, "\n\n\x1b[31m", log.LstdFlags)
	}

	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				if logger != nil {
					httprequest, _ := httputil.DumpRequest(c.Request, false)
					logger.Printf("[Recovery] panic recovered:\n%s\n%s\n", string(httprequest), err)
				}

				if !c.written {
					c.AbortWithStatus(500)
				}
			}
		}()
		c.Next()
	}
}

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
		logger = log.New(out, "", log.LstdFlags)
	}

	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				if logger != nil {
					httprequest, _ := httputil.DumpRequest(c.Request, false)
					logger.Printf("[Recover] panic recovered:\n%s\n%s\n", string(httprequest), err)
				}

				c.AbortWithStatus(500)
				return
			}
		}()
		c.Next()
	}
}

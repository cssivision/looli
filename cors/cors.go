package cors

import (
	"fmt"
	"github.com/cssivision/looli"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	defaultAllowOrigins = []string{"*"}
	defaultAllowMethods = []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPut,
		http.MethodPost,
		http.MethodDelete,
		http.MethodPatch,
	}
)

type CorsOption struct {
	// AllowedOrigins is a list of origins a cross-domain request can be executed from
	// If the special "*" value is present in the list, all origins will be allowed.
	AllowOrigins []string

	// AllowOriginFunc is a custom function to validate the origin. It take the origin
	// as argument and returns true if allowed or false otherwise. If this option is
	// set, the content of AllowedOrigins is ignored.
	AllowOriginsFunc func(string) bool

	// specifies the method or methods allowed when accessing the resource. This is used in response
	// to a preflight request, default is []string{"GET", "HEAD", "PUT", "POST", "DELETE", "PATCH"}.
	AllowMethods []string

	// header is used in response to a preflight request to indicate which HTTP headers can be used when
	// making the actual request
	AllowHeaders []string

	// The Access-Control-Allow-Credentials header Indicates whether or not the response to the request
	// can be exposed when the credentials flag is true.  When used as part of a response to a preflight
	// request, this indicates whether or not the actual request can be made using credentials. Note that
	// simple GET requests are not preflighted, and so if a request is made for a resource with credentials,
	// if this header is not returned with the resource, the response is ignored by the browser and not
	// returned to web content.
	AllowCredentials bool

	// ExposeHeaders header lets a server whitelist headers that browsers are allowed to access.
	// For example: Access-Control-Expose-Headers: "X-My-Custom-Header, X-Another-Custom-Header"
	// This allows the X-My-Custom-Header and X-Another-Custom-Header headers to be exposed to the browser.
	ExposeHeaders []string

	// MaxAge indicates how long the results of a preflight request can be cached.
	MaxAge time.Duration
}

func Cors(option CorsOption) looli.HandlerFunc {
	if option.AllowOrigins == nil {
		option.AllowOrigins = defaultAllowOrigins
	}

	if option.AllowMethods == nil {
		option.AllowMethods = defaultAllowMethods
	}

	if option.AllowOriginsFunc == nil {
		option.AllowOriginsFunc = func(origin string) bool {
			origin = strings.ToLower(origin)
			for _, o := range option.AllowOrigins {
				if o == origin || o == "*" {
					return true
				}
			}

			return false
		}
	}

	return func(c *looli.Context) {
		origin := c.Header("Origin")
		if origin == "" {
			return
		}

		if !option.AllowOriginsFunc(origin) {
			c.AbortWithStatus(http.StatusForbidden)
			c.String(fmt.Sprintf("Origin: %v is not allowed", origin))
			return
		}

		// Always set Vary headers
		c.ResponseWriter.Header().Add("Vary", "Origin")

		// handle preflight request
		if c.Method == http.MethodOptions {
			// Always set Vary headers
			c.ResponseWriter.Header().Add("Vary", "Access-Control-Request-Method")
			c.ResponseWriter.Header().Add("Vary", "Access-Control-Request-Headers")

			requestMethod := c.Header("Access-Control-Request-Method")

			// invalid preflighted request, missing Access-Control-Request-Method header
			if requestMethod == "" {
				c.AbortWithStatus(http.StatusForbidden)
				c.String("invalid preflighted request, missing Access-Control-Request-Method header")
				return
			}

			// methods allowed when accessing the resource with actual request
			if len(option.AllowMethods) > 0 {
				c.SetHeader("Access-Control-Allow-Methods", strings.Join(option.AllowMethods, ", "))
			}

			allowHeaders := ""
			if len(option.AllowHeaders) > 0 {
				allowHeaders = strings.Join(option.AllowHeaders, ", ")
			} else {
				// when AllowHeaders is nil, defualt allow all header
				allowHeaders = c.Request.Header.Get("Access-Control-Request-Headers")
			}

			if allowHeaders != "" {
				c.SetHeader("Access-Control-Allow-Headers", allowHeaders)
			}

			// if MaxAge is set means preflight request can be cached for MaxAge time.
			if option.MaxAge > 0 {
				c.SetHeader("Access-Control-Max-Age", strconv.Itoa(int(option.MaxAge.Seconds())))
			}

			// origin is alway a specific domain, not "*".
			// because when AllowCredentials == true, we must specify a domain, and cannot use wildcard.
			// origin is a specific domain, make it.
			c.SetHeader("Access-Control-Allow-Origin", origin)

			if option.AllowCredentials {
				c.SetHeader("Access-Control-Allow-Credentials", "true")
			}
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// origin is alway a specific domain, not "*".
		// because when AllowCredentials == true, we must specify a domain, and cannot use wildcard.
		// origin is a specific domain, make it.
		c.SetHeader("Access-Control-Allow-Origin", origin)

		if option.AllowCredentials {
			c.SetHeader("Access-Control-Allow-Credentials", "true")
		}

		// handle normal cors request
		if len(option.ExposeHeaders) > 0 {
			c.SetHeader("Access-Control-Expose-Headers", strings.Join(option.ExposeHeaders, ", "))
		}
	}
}

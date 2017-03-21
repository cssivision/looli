package csrf

import (
	"github.com/cssivision/looli"
	"net/http"
)

var (
	maxAge     = 12 * 3600
	cookieName = "_csrf"
)

type Options struct {
	FormKey   string
	HeaderKey string
	Skip      func(*looli.Context) bool
	MaxAge    int
	Domain    string
	Path      string
	HttpOnly  bool
	Secure    bool
}

func Default() looli.HandlerFunc {
	return New(Options{})
}

func New(options Options) looli.HandlerFunc {
	if options.FormKey == "" {
		options.FormKey = "csrf_token"
	}

	if options.HeaderKey == "" {
		options.HeaderKey = "X-CSRF-Token"
	}

	if options.MaxAge == 0 {
		options.MaxAge = maxAge
	}

	return func(c *looli.Context) {
		if options.Skip != nil && options.Skip(c) {
			return
		}

		// Set the Vary: Cookie header to protect clients from caching the response.
		c.ResponseWriter.Header().Add("Vary", "Cookie")

		csrfToken := c.PostForm(options.FormKey)
		if csrfToken == "" {
			csrfToken = c.Header(options.HeaderKey)
		}

		if csrfToken == "" || !verify(getSecret(c), csrfToken) {
			c.AbortWithStatus(http.StatusForbidden)
			c.String("invalid csrf token")
			return
		}
	}
}

func getSecret(c *looli.Context) string {
	value, err := c.Cookie(cookieName)
	var secretCookie *http.Cookie

	if err != nil {
		secretCookie = &http.Cookie{}
		secretCookie.Name = cookieName
		secretCookie.Value = newSecret(c)
		secretCookie.MaxAge = maxAge
		value = secretCookie.Value
		c.SetCookie(secretCookie)
	}
	return value
}

func NewToken(c *looli.Context) string {
	return ""
}

func newSecret(c *looli.Context) string {
	return ""
}

func verify(secret, token string) bool {
	return true
}

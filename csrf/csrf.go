package csrf

import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base64"
	"github.com/cssivision/looli"
	"net/http"
	"strings"
)

var (
	maxAge        = 12 * 3600
	cookieName    = "_csrf"
	formKey       = "csrf_token"
	headerKey     = "X-CSRF-Token"
	cookieOptions = &http.Cookie{}
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
		options.FormKey = formKey
	}

	if options.HeaderKey == "" {
		options.HeaderKey = headerKey
	}

	if options.MaxAge == 0 {
		options.MaxAge = maxAge
	}

	cookieOptions.MaxAge = options.MaxAge
	cookieOptions.Path = options.Path
	cookieOptions.HttpOnly = options.HttpOnly
	cookieOptions.Secure = options.Secure
	cookieOptions.Domain = options.Domain

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
		*secretCookie = *cookieOptions
		secretCookie.Name = cookieName
		secretCookie.Value = randomKey(16)
		value = secretCookie.Value
		c.SetCookie(secretCookie)
	}
	return value
}

func NewToken(c *looli.Context) string {
	secret := getSecret(c)
	salt := randomKey(16)
	return generateWithSalt(secret, salt)
}

func randomKey(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func verify(secret, token string) bool {
	i := strings.Index(token, ".")

	if i == -1 {
		return false
	}

	return subtle.ConstantTimeCompare([]byte(token), []byte(generateWithSalt(secret, token[0:i]))) == 1
}

func generateWithSalt(secret, salt string) string {
	h := sha1.New()
	h.Write([]byte(salt + "." + secret))

	return salt + "." + base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

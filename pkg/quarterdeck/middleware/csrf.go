package middleware

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Parameters and headers for double-cookie submit CSRF protection
const (
	CSRFCookie          = "csrf_token"
	CSRFReferenceCookie = "csrf_reference_token"
	CSRFHeader          = "X-CSRF-TOKEN"
)

// DoubleCookie is a Cross-Site Request Forgery (CSR/XSRF) protection middleware that
// checks the presence of an X-CSRF-TOKEN header containing a cryptographically random
// token that matches a token contained in the CSRF-TOKEN cookie in the request.
// Because of the same-origin poicy, an attacker cannot access the cookies or scripts
// of the safe site, therefore the X-CSRF-TOKEN header cannot be forged, and if it is
// omitted because it is being re-posted by an attacker site then the request will be
// rejected with a 403 error. Note that this protection requires TLS to prevent MITM.
func DoubleCookie() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(CSRFReferenceCookie)
		if err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusForbidden, api.ErrorResponse(ErrCSRFVerification))
			return
		}

		header := c.GetHeader(CSRFHeader)
		if header, err = url.QueryUnescape(header); err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, api.ErrorResponse(err))
			return
		}

		if cookie == "" || header == "" {
			log.Warn().Bool("header_exists", header != "").Bool("cookie_exists", cookie != "").Msg("missing either csrf token header or reference cookie")
			c.AbortWithStatusJSON(http.StatusForbidden, api.ErrorResponse(ErrCSRFVerification))
			return
		}

		if cookie != header {
			log.Warn().Bool("header_exists", header != "").Bool("cookie_exists", cookie != "").Msg("csrf token cookie/header mismatch")
			c.AbortWithStatusJSON(http.StatusForbidden, api.ErrorResponse(ErrCSRFVerification))
			return
		}

		c.Next()
	}
}

// SetDoubleCookieToken is a helper function to set cookies on a gin request.
func SetDoubleCookieToken(c *gin.Context, domain string, expires time.Time) error {
	// Generate a secure token
	token, err := GenerateCSRFToken()
	if err != nil {
		return err
	}

	// Set the reference cookie as http only but allow access to the csrf cookie so that the
	// front-end can fetch it to add it to the X-CSRF-TOKEN header in the request.
	maxAge := int((time.Until(expires)).Seconds()) + 60
	c.SetCookie(CSRFReferenceCookie, token, maxAge, "/", domain, true, true)
	c.SetCookie(CSRFCookie, token, maxAge, "/", domain, true, false)
	return nil
}

// Random seed is an additional barrier to cryptanalysis and is unique to each quarterdeck process.
var (
	seed     []byte
	initseed sync.Once
)

func GenerateCSRFToken() (_ string, err error) {
	// Ensure the seed is generated the first time this method is called
	initseed.Do(func() {
		seed = make([]byte, 16)
		if _, err = rand.Read(seed); err != nil {
			log.WithLevel(zerolog.FatalLevel).Err(err).Msg("could not set csrf token generator random seed")
		}
	})

	if err != nil {
		return "", err
	}

	nonce := make([]byte, 32)
	if _, err = rand.Read(nonce); err != nil {
		return "", err
	}

	sig := sha256.New()
	sig.Write(seed)
	sig.Write(nonce)

	return base64.URLEncoding.EncodeToString(sig.Sum(nil)), nil
}

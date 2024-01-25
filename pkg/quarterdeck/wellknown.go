package quarterdeck

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
)

// JWKS returns the JSON web key set for the public RSA keys that are currently being
// used by Quarterdeck to sign JWT acccess and refresh tokens. External callers can use
// these keys to verify that a JWT token was in fact issued by the Quarterdeck API.
// TODO: add Cache-Control or Expires header to the response
// TODO: move the jwks construction to the token manager for easier key management.
func (s *Server) JWKS(c *gin.Context) {
	keys, err := s.tokens.Keys()
	if err != nil {
		sentry.Error(c).Err(err).Msg("could not get token keys")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("an internal error occurred"))
		return
	}
	c.JSON(http.StatusOK, keys)
}

// Returns a JSON document with the OpenID configuration as defined by the OpenID
// Connect standard: https://connect2id.com/learn/openid-connect. This document helps
// clients understand how to authenticate with Quarterdeck.
// TODO: once OpenID endpoints have been configured add them to this JSON response
func (s *Server) OpenIDConfiguration(c *gin.Context) {
	// Parse the token issuer for the OpenID configuration
	base, err := url.Parse(s.conf.Token.Issuer)
	if err != nil {
		sentry.Error(c).Err(err).Msg("could not parse issuer")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("openid is not configured correctly"))
		return
	}

	openid := &api.OpenIDConfiguration{
		Issuer:                        base.ResolveReference(&url.URL{Path: "/"}).String(),
		JWKSURI:                       base.ResolveReference(&url.URL{Path: "/.well-known/jwks.json"}).String(),
		ScopesSupported:               []string{"openid", "profile", "email"},
		ResponseTypesSupported:        []string{"token", "id_token"},
		CodeChallengeMethodsSupported: []string{"S256", "plain"},
		ResponseModesSupported:        []string{"query", "fragment", "form_post"},
		SubjectTypesSupported:         []string{"public"},
		IDTokenSigningAlgValues:       []string{"HS256", "RS256"},
		TokenEndpointAuthMethods:      []string{"client_secret_basic", "client_secret_post"},
		ClaimsSupported:               []string{"aud", "email", "exp", "iat", "iss", "sub"},
		RequestURIParameterSupported:  false,
	}

	c.JSON(http.StatusOK, openid)
}

// Writes the security.txt file generated from https://securitytxt.org/ and digitally
// signed with the info@rotational.io PGP keys to alert security researchers to our
// security policies and allow them to contact us with any security flaws.
func (s *Server) SecurityTxt(c *gin.Context) {
	c.String(http.StatusOK, securityTxt)
}

const securityTxt = `-----BEGIN PGP SIGNED MESSAGE-----
Hash: SHA256

Contact: mailto:info@rotational.io
Expires: 2025-01-22T00:00:00.000Z
Encryption: https://keys.openpgp.org/vks/v1/by-fingerprint/8976758EC922C9362520F83E953A902549826C95
Preferred-Languages: en
Canonical: https://auth.rotational.app/.well-known/security.txt
-----BEGIN PGP SIGNATURE-----

iQIzBAEBCAAdFiEEiXZ1jskiyTYlIPg+lTqQJUmCbJUFAmWy38YACgkQlTqQJUmC
bJWi6g/8DWaMK3688/Xz8upaD5lrSdo+ENChYTdNZ1E7s7XGpuBdftNMc0+aPHKI
rp3HwPOnOzlT2Servup9H8zspFE2P/BSQuTJg13/AgVA53z8T/3wLXPjWRUf4slU
+6WlEi/r1lGGWwl2QyHdBJbbKz0uFH/xl9aAshzjuf2wLrDiOh1QWawtVfeRUkiH
fspHm8lzkP2uqVIH5OKNK2lT4zoy/zUjo38dyIPmNjtJbDe5xgmc+qcVaqBmy2Se
iMNf0XihYuBGfFVeHoAHseJvLJjseojteQ4ayWyWff5AbRKuJvwIPk3rDfJBWhrO
8y7qCimpO7hb0kxRT2oJifq2ihMBsdzuOJYRJIQLyUbtU3IULYiaYaDdOrsQ53tL
9KxkUMuduUHM69Yb72jh/3mcTmZjcTeGnRXaZi3l/fV9JjIQmTl+sku59MBAIJTR
GjZ3FZmxotezLseKrZmhXlwToYQSQ4YzxqsxSJgS2qFiuzKHY5k/lRAknv+FdrD7
7sKLSSKGjL1kEq0A47UuseNbrO43XwQRFnuBQeMwUN/XPF71TMw2UicS3iARgJAj
udl6gQVMyHtJCHNW6lwRjma63IB5MuUAyMcqLtG1pynbA/6FAE/e8PScMxY05u7/
wxWfJ8C0RRVFGSWri0TZa2S2O/hqbbOcBujOsJrH2UKHr8ex44s=
=OmQE
-----END PGP SIGNATURE-----
`

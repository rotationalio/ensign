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
Expires: 2026-04-21T04:00:00.000Z
Encryption: https://keys.openpgp.org/vks/v1/by-fingerprint/04AD032D06627D133451713BD169CA44AA700A0C
Preferred-Languages: en
Canonical: https://auth.rotational.app/.well-known/security.txt
-----BEGIN PGP SIGNATURE-----

iQIzBAEBCAAdFiEEBK0DLQZifRM0UXE70WnKRKpwCgwFAmgGawEACgkQ0WnKRKpw
CgwN/w/+JNeNuMY5Tp8IT+EoD7YstbWoCrIYoy/Cga8PsyRcFAlsROeziv8VWJmj
HdjYMAZD/xMyjFBEDlvYLtOU3OL02ncnuUUWewSgmkaF8DBa4ie1yiv0WfnLgK+C
lZgbkY5LOiXUA6HzCDHONqHR/tRS2bmn2LxjZehyKdY5A+fZ9Mh1ap6U9y1jm4PH
eg3YxI6tygv0m9uQDpaT7kX0F/O3j7BAyrsU5dMx3/RMaBJMU+RB1Wo3TAP2PJCC
5MQpR6hNiTs2VbgvHebFZaJp2bbA9JqeQVA5d3AKevU7s26d/Ufq1J4WfKCPCwFn
gqFHsTynx/wPkmacJRH23o8IHS2dilVzfjDC8JhkL/ZyLi9r1UFSZjMwd58cn/fD
F+DbwexIEcrhfVAIU/YjCw7gSnJ+iPffW+xYdEjEl4fH96aUUk+mk8HDB255A8bs
dBOd7uf1yAlUCXl4zP8claZfevUYVqmb09aIDrBMXgNB/YtfeBE97jpXQ1hlaa1t
YrTKQShnKHUj0vswMNrIridgFYTzqT77hWiBkZgGR/4wvxkQZylB/5ivf30dG87x
pU8qOM3hCHoCsPPc0pcc7/YeYB+7EHnO8/mkebeHIpkKwBRYM4pkvgVoXkHEhZrc
RtqJJn0GDFOyJJjj+kWjQuHdM5zhZBxQN4VuZLh7ymFnO2mo8eE=
=lhER
-----END PGP SIGNATURE-----
`

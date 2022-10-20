package quarterdeck

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rs/zerolog/log"
)

// JWKS returns the JSON web key set for the public RSA keys that are currently being
// used by Quarterdeck to sign JWT acccess and refresh tokens. External callers can use
// these keys to verify that the a JWT token was in fact issued by the Quarterdeck API.
func (s *Server) JWKS(c *gin.Context) {
	jwks := jwk.NewSet()
	for keyid, pubkey := range s.tokens.Keys() {
		key, err := jwk.FromRaw(pubkey)
		if err != nil {
			log.Error().Err(err).Str("kid", keyid.String()).Msg("could not parse tokens public key")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("an internal error occurred"))
			return
		}

		if err = key.Set(jwk.KeyIDKey, keyid.String()); err != nil {
			log.Error().Err(err).Str("kid", keyid.String()).Msg("could not set tokens public key id")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("an internal error occurred"))
			return
		}

		if err = key.Set(jwk.KeyUsageKey, jwk.ForSignature); err != nil {
			log.Error().Err(err).Str("kid", keyid.String()).Msg("could not set tokens public key use")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("an internal error occurred"))
			return
		}

		// NOTE: the algorithm should match the signing method in tokens.go
		if err = key.Set(jwk.AlgorithmKey, jwa.RS256); err != nil {
			log.Error().Err(err).Str("kid", keyid.String()).Msg("could not set tokens public key algorithm")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("an internal error occurred"))
			return
		}

		if err = jwks.AddKey(key); err != nil {
			log.Error().Err(err).Str("kid", keyid.String()).Msg("could not add key to jwks")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("an internal error occurred"))
			return
		}
	}
	c.JSON(http.StatusOK, jwks)
}

// Returns a JSON document with the OpenID configuration as defined by the OpenID
// Connect standard: https://connect2id.com/learn/openid-connect. This document helps
// clients understand how to authenticate with Quarterdeck.
// TODO: once OpenID endpoints have been configured add them to this JSON response
func (s *Server) OpenIDConfiguration(c *gin.Context) {
	// Parse the token issuer for the OpenID configuration
	base, err := url.Parse(s.conf.Token.Issuer)
	if err != nil {
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
		RequestURIPArameterSupported:  false,
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
Contact: https://rotational.io/contact/
Expires: 2024-01-21T18:00:00.000Z
Encryption: https://keys.openpgp.org/vks/v1/by-fingerprint/8976758EC922C9362520F83E953A902549826C95
Preferred-Languages: en
Canonical: https://auth.rotational.app/.well-known/security.txt
-----BEGIN PGP SIGNATURE-----

iQIzBAEBCAAdFiEEiXZ1jskiyTYlIPg+lTqQJUmCbJUFAmNRQ20ACgkQlTqQJUmC
bJXlAg/+IlL16JQgTLkF0QhNI6pIn42KIynVBZXbzj8qlOvwWZxqFOZvOZHSSK7A
lzjSlfJ+6CDeE4jBUgXmmY9P413CoVINKvkzlGehUFQEy7kHc/WKTmm5nuMDsY8I
uK+pp4Yl3mZ6+ey5Lx8WavlIuACh3YD41OuwXQTw4MfOO7a6dCkQ6WVRcUObw4sS
kYC/cF6VFF6Uv0QhlMs/ofRbpQAdeahgJyvADB4HOqvE3C8OkZlGwiUIm1/rePYd
QxkykSn72ry/Y879kasPbhAYzrEJBuM2ejQ9L62XslzIzFBpWc75RI5a3jAPyFzw
DXzByA4NP4jQC64Vrh8v08eZcevour4PRi4d/pfb2oD7znmSdz9dgebqqcu6WiAZ
Z4CCjTcgSHBnZD5dbBWFly2Xlz0DWPoZxXjs3KQvGEs9Yyw7D9cMSrqI2rqOvY4h
lfRhq35ANePMpVQOG6U7EwoWruGM26SLhYHeD6ylUVLd/E9Yd8XqECtSKXEYSfQG
0BItbSWul9b1Ou7Q/iBSBVFys0b4LeHZrW2lwoA8KQDYOLBydQ14Ul78Ybc0zsRl
9PR5mp7hcNZ7E1oCbJUgHLbdGx/vdU+J1ckmRzimiRslnQSd56P5QzE+Bb6v4GQU
xO0TPc6AN1sdXOi+d0NOcQ0edk3Bt1TsxUKqvhhmytXVtLlDslY=
=o6ob
-----END PGP SIGNATURE-----
`

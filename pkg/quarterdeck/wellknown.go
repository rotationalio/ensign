package quarterdeck

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
)

// JWKS returns the JSON web key set for the public RSA keys that are currently being
// used by Quarterdeck to sign JWT acccess and refresh tokens. External callers can use
// these keys to verify that the a JWT token was in fact issued by the Quarterdeck API.
func (s *Server) JWKS(c *gin.Context) {
	jwks := jwk.NewSet()
	c.JSON(http.StatusOK, jwks)
}

// Writes the security.txt file generated from https://securitytxt.org/ and digitally
// signed with the info@rotational.io PGP keys to alert security researchers to our
// security policies and allow them to contact us with any security flaws.
func (s *Server) SecurityTxt(c *gin.Context) {
	c.String(http.StatusOK, securityTxt)
}

func (s *Server) OpenIDConfiguration(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, api.ErrorResponse("not implemented yet"))
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

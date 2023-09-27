package passwd

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"golang.org/x/crypto/argon2"
)

//===========================================================================
// Derived Key Algorithm
//===========================================================================

// Argon2 constants for the derived key (dk) algorithm
// See: https://cryptobook.nakov.com/mac-and-key-derivation/argon2
const (
	dkAlg  = "argon2id"        // the derived key algorithm
	dkTime = uint32(1)         // draft RFC recommends time = 1
	dkMem  = uint32(64 * 1024) // draft RFC recommends memory as ~64MB (or as much as possible)
	dkProc = uint8(2)          // can be set to the number of available CPUs
	dkSLen = 16                // the length of the salt to generate per user
	dkKLen = uint32(32)        // the length of the derived key (32 bytes is the required key size for AES-256)
)

// Argon2 variables for the derived key (dk) algorithm
var (
	dkParse = regexp.MustCompile(`^\$(?P<alg>[\w\d]+)\$v=(?P<ver>\d+)\$m=(?P<mem>\d+),t=(?P<time>\d+),p=(?P<procs>\d+)\$(?P<salt>[\+\/\=a-zA-Z0-9]+)\$(?P<key>[\+\/\=a-zA-Z0-9]+)$`)
)

// CreateDerivedKey creates an encoded derived key with a random hash for the password.
func CreateDerivedKey(password string) (_ string, err error) {
	if password == "" {
		return "", errors.New("cannot create derived key for empty password")
	}

	salt := make([]byte, dkSLen)
	if _, err = rand.Read(salt); err != nil {
		return "", fmt.Errorf("could not generate %d length salt: %s", dkSLen, err)
	}

	dk := argon2.IDKey([]byte(password), salt, dkTime, dkMem, dkProc, dkKLen)
	b64salt := base64.StdEncoding.EncodeToString(salt)
	b64dk := base64.StdEncoding.EncodeToString(dk)
	return fmt.Sprintf("$%s$v=%d$m=%d,t=%d,p=%d$%s$%s", dkAlg, argon2.Version, dkMem, dkTime, dkProc, b64salt, b64dk), nil
}

// VerifyDerivedKey checks that the submitted password matches the derived key.
func VerifyDerivedKey(dk, password string) (_ bool, err error) {
	if dk == "" || password == "" {
		return false, errors.New("cannot verify empty derived key or password")
	}

	dkb, salt, t, m, p, err := ParseDerivedKey(dk)
	if err != nil {
		return false, err
	}

	vdk := argon2.IDKey([]byte(password), salt, t, m, p, uint32(len(dkb)))
	return bytes.Equal(dkb, vdk), nil
}

// ParseDerivedKey returns the parts of the encoded derived key string.
func ParseDerivedKey(encoded string) (dk, salt []byte, time, memory uint32, threads uint8, err error) {
	if !dkParse.MatchString(encoded) {
		return nil, nil, 0, 0, 0, errors.New("cannot parse encoded derived key, does not match regular expression")
	}
	parts := dkParse.FindStringSubmatch(encoded)

	if len(parts) != 8 {
		return nil, nil, 0, 0, 0, errors.New("cannot parse encoded derived key, matched expression does not contain enough subgroups")
	}

	// check the algorithm
	if parts[1] != dkAlg {
		return nil, nil, 0, 0, 0, fmt.Errorf("current code only works with the the dk protcol %q not %q", dkAlg, parts[1])
	}

	// check the version
	if version, err := strconv.Atoi(parts[2]); err != nil || version != argon2.Version {
		return nil, nil, 0, 0, 0, fmt.Errorf("expected %s version %d got %q", dkAlg, argon2.Version, parts[2])
	}

	var (
		time64    uint64
		memory64  uint64
		threads64 uint64
	)

	if memory64, err = strconv.ParseUint(parts[3], 10, 32); err != nil {
		return nil, nil, 0, 0, 0, fmt.Errorf("could not parse memory %q: %s", parts[3], err)
	}
	memory = uint32(memory64)

	if time64, err = strconv.ParseUint(parts[4], 10, 32); err != nil {
		return nil, nil, 0, 0, 0, fmt.Errorf("could not parse time %q: %s", parts[4], err)
	}
	time = uint32(time64)

	if threads64, err = strconv.ParseUint(parts[5], 10, 8); err != nil {
		return nil, nil, 0, 0, 0, fmt.Errorf("could not parse threads %q: %s", parts[5], err)
	}
	threads = uint8(threads64)

	if salt, err = base64.StdEncoding.DecodeString(parts[6]); err != nil {
		return nil, nil, 0, 0, 0, fmt.Errorf("could not parse salt: %s", err)
	}

	if dk, err = base64.StdEncoding.DecodeString(parts[7]); err != nil {
		return nil, nil, 0, 0, 0, fmt.Errorf("could not parse derived key: %s", err)
	}

	return dk, salt, time, memory, threads, nil
}

func IsDerivedKey(s string) bool {
	return dkParse.MatchString(s)
}

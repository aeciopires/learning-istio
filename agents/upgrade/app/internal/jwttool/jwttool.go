// Package jwttool generates RSA signing keys and creates, decodes and
// verifies JSON Web Tokens, so the walkthroughs in this repository can
// mint one token per application and exercise Istio's JWT-based
// policies (RequestAuthentication + AuthorizationPolicy) end to end.
//
// It is a thin layer over github.com/golang-jwt/jwt/v5: that library
// signs and verifies, this package adds what Istio needs around it -
// a JWKS document for RequestAuthentication's "jwks" field and custom
// claims shaped the way AuthorizationPolicy can match them (a string or
// a list of strings).
//
// Everything here is offline. Keys and tokens produced by this package
// are for study clusters only.
package jwttool

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// MinKeyBits is the smallest RSA modulus GenerateKey accepts.
const MinKeyBits = 2048

// GenerateKey creates a new RSA private key of the given size.
func GenerateKey(bits int) (*rsa.PrivateKey, error) {
	if bits < MinKeyBits {
		return nil, fmt.Errorf("RSA key size %d is too small (minimum %d)", bits, MinKeyBits)
	}
	return rsa.GenerateKey(rand.Reader, bits)
}

// EncodePrivateKeyPEM serializes key as a PKCS#8 "PRIVATE KEY" PEM block.
func EncodePrivateKeyPEM(key *rsa.PrivateKey) ([]byte, error) {
	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("marshaling private key: %w", err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}), nil
}

// EncodePublicKeyPEM serializes key as a PKIX "PUBLIC KEY" PEM block,
// the format jwt.io expects when verifying an RS256 signature.
func EncodePublicKeyPEM(key *rsa.PublicKey) ([]byte, error) {
	der, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return nil, fmt.Errorf("marshaling public key: %w", err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}), nil
}

// ParsePrivateKeyPEM reads an RSA private key in PKCS#8 or PKCS#1 PEM form.
func ParsePrivateKeyPEM(data []byte) (*rsa.PrivateKey, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(data)
	if err != nil {
		return nil, fmt.Errorf("parsing RSA private key: %w", err)
	}
	return key, nil
}

// ParsePublicKeyPEM reads an RSA public key in PKIX or PKCS#1 PEM form.
func ParsePublicKeyPEM(data []byte) (*rsa.PublicKey, error) {
	key, err := jwt.ParseRSAPublicKeyFromPEM(data)
	if err != nil {
		return nil, fmt.Errorf("parsing RSA public key: %w", err)
	}
	return key, nil
}

// JWK is a single RSA public key in JSON Web Key form (RFC 7517).
type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// JWKS is a JSON Web Key Set - the document Istio's
// RequestAuthentication reads from "jwks" (inline) or "jwksUri".
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// BuildJWKS wraps a public key in a single-key JWKS. kid must match the
// "kid" header of the tokens signed with the corresponding private key.
func BuildJWKS(key *rsa.PublicKey, kid string) (JWKS, error) {
	if key == nil {
		return JWKS{}, errors.New("public key is required")
	}
	if kid == "" {
		return JWKS{}, errors.New("kid is required")
	}
	return JWKS{Keys: []JWK{{
		Kty: "RSA",
		Use: "sig",
		Alg: jwt.SigningMethodRS256.Alg(),
		Kid: kid,
		N:   base64.RawURLEncoding.EncodeToString(key.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(key.E)).Bytes()),
	}}}, nil
}

// SignOptions describes the token Sign creates.
type SignOptions struct {
	// Kid is written to the token's "kid" header so a verifier can pick
	// the matching key out of a JWKS.
	Kid string
	// Issuer and Subject become the "iss" and "sub" claims. Istio builds
	// the request principal from them as "<iss>/<sub>".
	Issuer  string
	Subject string
	// Audiences becomes the "aud" claim; omitted when empty.
	Audiences []string
	// ExpiresIn is how long the token stays valid, counted from Now.
	ExpiresIn time.Duration
	// Now is the issue time; the zero value means time.Now().
	Now time.Time
	// Claims are extra, application-defined claims. Registered claims
	// (iss, sub, aud, exp, iat, nbf) cannot be overridden here.
	Claims map[string]any
}

var reservedClaims = map[string]bool{"iss": true, "sub": true, "aud": true, "exp": true, "iat": true, "nbf": true}

// Sign creates an RS256-signed token.
func Sign(key *rsa.PrivateKey, opts SignOptions) (string, error) {
	if key == nil {
		return "", errors.New("private key is required")
	}
	if opts.Issuer == "" {
		return "", errors.New("issuer is required")
	}
	if opts.Subject == "" {
		return "", errors.New("subject is required")
	}
	if opts.ExpiresIn <= 0 {
		return "", errors.New("expiration must be positive")
	}
	now := opts.Now
	if now.IsZero() {
		now = time.Now()
	}

	claims := jwt.MapClaims{}
	for name, value := range opts.Claims {
		if reservedClaims[name] {
			return "", fmt.Errorf("claim %q is a registered claim; set it through its own option", name)
		}
		claims[name] = value
	}
	claims["iss"] = opts.Issuer
	claims["sub"] = opts.Subject
	claims["iat"] = now.Unix()
	claims["exp"] = now.Add(opts.ExpiresIn).Unix()
	if len(opts.Audiences) > 0 {
		claims["aud"] = opts.Audiences
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	if opts.Kid != "" {
		token.Header["kid"] = opts.Kid
	}
	signed, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}
	return signed, nil
}

// Decoded is a token split into its JSON parts.
type Decoded struct {
	Header map[string]any `json:"header"`
	Claims map[string]any `json:"claims"`
	// Verified reports whether the signature and time-based claims were
	// checked against a public key. A decoded but unverified token must
	// never be trusted.
	Verified bool `json:"verified"`
}

// Decode reads a token WITHOUT checking its signature or expiry - the
// same thing pasting it into a JWT debugger does.
func Decode(token string) (Decoded, error) {
	claims := jwt.MapClaims{}
	parsed, _, err := jwt.NewParser().ParseUnverified(strings.TrimSpace(token), claims)
	if err != nil {
		return Decoded{}, fmt.Errorf("decoding token: %w", err)
	}
	return Decoded{Header: parsed.Header, Claims: claims}, nil
}

// Verify decodes a token and checks its RS256 signature and its
// time-based claims against key.
func Verify(token string, key *rsa.PublicKey) (Decoded, error) {
	if key == nil {
		return Decoded{}, errors.New("public key is required")
	}
	claims := jwt.MapClaims{}
	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}))
	parsed, err := parser.ParseWithClaims(strings.TrimSpace(token), claims, func(*jwt.Token) (any, error) {
		return key, nil
	})
	if err != nil {
		return Decoded{}, fmt.Errorf("verifying token: %w", err)
	}
	return Decoded{Header: parsed.Header, Claims: claims, Verified: true}, nil
}

// ParseClaimFlags turns repeated "name=value" flags into custom claims.
// scalar entries become string claims; list entries are split on commas
// into a list of strings. Those are the two claim shapes Istio's
// AuthorizationPolicy can match with request.auth.claims[...].
func ParseClaimFlags(scalar, list []string) (map[string]any, error) {
	claims := make(map[string]any, len(scalar)+len(list))
	add := func(raw string, asList bool) error {
		name, value, found := strings.Cut(raw, "=")
		name = strings.TrimSpace(name)
		if !found || name == "" {
			return fmt.Errorf("invalid claim %q: expected name=value", raw)
		}
		if _, dup := claims[name]; dup {
			return fmt.Errorf("claim %q given more than once", name)
		}
		if !asList {
			claims[name] = value
			return nil
		}
		var items []string
		for _, item := range strings.Split(value, ",") {
			if item = strings.TrimSpace(item); item != "" {
				items = append(items, item)
			}
		}
		if len(items) == 0 {
			return fmt.Errorf("invalid claim %q: list has no values", raw)
		}
		claims[name] = items
		return nil
	}
	for _, raw := range scalar {
		if err := add(raw, false); err != nil {
			return nil, err
		}
	}
	for _, raw := range list {
		if err := add(raw, true); err != nil {
			return nil, err
		}
	}
	return claims, nil
}

// MarshalIndent renders v as indented JSON followed by a newline.
func MarshalIndent(v any) ([]byte, error) {
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(out, '\n'), nil
}

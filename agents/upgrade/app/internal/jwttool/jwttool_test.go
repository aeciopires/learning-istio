package jwttool

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"strings"
	"testing"
	"time"
)

func testKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	key, err := GenerateKey(MinKeyBits)
	if err != nil {
		t.Fatalf("GenerateKey: %v", err)
	}
	return key
}

func TestGenerateKeyRejectsSmallKeys(t *testing.T) {
	if _, err := GenerateKey(1024); err == nil {
		t.Fatal("expected an error for a 1024-bit key")
	}
}

func TestPEMRoundTrip(t *testing.T) {
	key := testKey(t)

	privPEM, err := EncodePrivateKeyPEM(key)
	if err != nil {
		t.Fatalf("EncodePrivateKeyPEM: %v", err)
	}
	gotPriv, err := ParsePrivateKeyPEM(privPEM)
	if err != nil {
		t.Fatalf("ParsePrivateKeyPEM: %v", err)
	}
	if !gotPriv.Equal(key) {
		t.Error("private key changed across a PEM round trip")
	}

	pubPEM, err := EncodePublicKeyPEM(&key.PublicKey)
	if err != nil {
		t.Fatalf("EncodePublicKeyPEM: %v", err)
	}
	gotPub, err := ParsePublicKeyPEM(pubPEM)
	if err != nil {
		t.Fatalf("ParsePublicKeyPEM: %v", err)
	}
	if !gotPub.Equal(&key.PublicKey) {
		t.Error("public key changed across a PEM round trip")
	}
}

func TestBuildJWKS(t *testing.T) {
	key := testKey(t)
	jwks, err := BuildJWKS(&key.PublicKey, "learning-istio")
	if err != nil {
		t.Fatalf("BuildJWKS: %v", err)
	}
	if len(jwks.Keys) != 1 {
		t.Fatalf("expected 1 key, got %d", len(jwks.Keys))
	}
	jwk := jwks.Keys[0]
	if jwk.Kty != "RSA" || jwk.Alg != "RS256" || jwk.Use != "sig" || jwk.Kid != "learning-istio" {
		t.Errorf("unexpected JWK metadata: %+v", jwk)
	}
	n, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		t.Fatalf("decoding n: %v", err)
	}
	if new(big.Int).SetBytes(n).Cmp(key.N) != 0 {
		t.Error("JWK modulus does not match the key")
	}
	e, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		t.Fatalf("decoding e: %v", err)
	}
	if int(new(big.Int).SetBytes(e).Int64()) != key.E {
		t.Error("JWK exponent does not match the key")
	}

	if _, err := BuildJWKS(&key.PublicKey, ""); err == nil {
		t.Error("expected an error for an empty kid")
	}
}

func TestSignAndVerify(t *testing.T) {
	key := testKey(t)
	token, err := Sign(key, SignOptions{
		Kid:       "learning-istio",
		Issuer:    "learning-istio@example.com",
		Subject:   "curl",
		Audiences: []string{"helloworld"},
		ExpiresIn: time.Hour,
		Claims:    map[string]any{"app": "curl", "groups": []string{"readers", "team-a"}},
	})
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}

	got, err := Verify(token, &key.PublicKey)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if !got.Verified {
		t.Error("expected Verified to be true")
	}
	if got.Header["kid"] != "learning-istio" || got.Header["alg"] != "RS256" {
		t.Errorf("unexpected header: %v", got.Header)
	}
	if got.Claims["iss"] != "learning-istio@example.com" || got.Claims["sub"] != "curl" || got.Claims["app"] != "curl" {
		t.Errorf("unexpected claims: %v", got.Claims)
	}
	groups, ok := got.Claims["groups"].([]any)
	if !ok || len(groups) != 2 || groups[0] != "readers" {
		t.Errorf("expected groups to be a 2-item list, got %#v", got.Claims["groups"])
	}

	// The payload must be plain JSON any JWT debugger can read.
	payload, err := base64.RawURLEncoding.DecodeString(strings.Split(token, ".")[1])
	if err != nil {
		t.Fatalf("decoding payload: %v", err)
	}
	if !json.Valid(payload) {
		t.Errorf("payload is not valid JSON: %s", payload)
	}
}

func TestVerifyRejectsWrongKeyAndExpiredTokens(t *testing.T) {
	key := testKey(t)
	other := testKey(t)

	token, err := Sign(key, SignOptions{Issuer: "iss", Subject: "sub", ExpiresIn: time.Hour})
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	if _, err := Verify(token, &other.PublicKey); err == nil {
		t.Error("expected verification with a different key to fail")
	}

	expired, err := Sign(key, SignOptions{
		Issuer: "iss", Subject: "sub", ExpiresIn: time.Minute, Now: time.Now().Add(-time.Hour),
	})
	if err != nil {
		t.Fatalf("Sign: %v", err)
	}
	if _, err := Verify(expired, &key.PublicKey); err == nil {
		t.Error("expected an expired token to fail verification")
	}
	// Decode must still read it: it never checks expiry or signature.
	decoded, err := Decode(expired)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if decoded.Verified {
		t.Error("Decode must never report a token as verified")
	}
}

func TestSignValidation(t *testing.T) {
	key := testKey(t)
	cases := map[string]SignOptions{
		"missing issuer":      {Subject: "s", ExpiresIn: time.Hour},
		"missing subject":     {Issuer: "i", ExpiresIn: time.Hour},
		"non-positive expiry": {Issuer: "i", Subject: "s"},
		"reserved claim":      {Issuer: "i", Subject: "s", ExpiresIn: time.Hour, Claims: map[string]any{"exp": "1"}},
	}
	for name, opts := range cases {
		if _, err := Sign(key, opts); err == nil {
			t.Errorf("%s: expected an error", name)
		}
	}
}

func TestParseClaimFlags(t *testing.T) {
	claims, err := ParseClaimFlags([]string{"app=curl", "note=a=b"}, []string{"groups=readers, team-a"})
	if err != nil {
		t.Fatalf("ParseClaimFlags: %v", err)
	}
	if claims["app"] != "curl" {
		t.Errorf("app = %v", claims["app"])
	}
	if claims["note"] != "a=b" {
		t.Errorf("only the first '=' should split name from value, got %v", claims["note"])
	}
	groups, ok := claims["groups"].([]string)
	if !ok || len(groups) != 2 || groups[1] != "team-a" {
		t.Errorf("groups = %#v", claims["groups"])
	}

	for _, bad := range [][2][]string{
		{{"novalue"}, nil},
		{{"=x"}, nil},
		{{"a=1", "a=2"}, nil},
		{nil, {"groups= , "}},
	} {
		if _, err := ParseClaimFlags(bad[0], bad[1]); err == nil {
			t.Errorf("expected an error for %v / %v", bad[0], bad[1])
		}
	}
}

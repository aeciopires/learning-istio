package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// runCLI executes the root command with args and returns its stdout.
func runCLI(t *testing.T, stdin string, args ...string) (string, error) {
	t.Helper()
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	rootCmd.SetErr(&out)
	rootCmd.SetIn(strings.NewReader(stdin))
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	return out.String(), err
}

func TestJWTCommands_EndToEnd(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "keys")

	if _, err := runCLI(t, "", "jwt", "keygen", "--out-dir", dir, "--kid", "test-kid"); err != nil {
		t.Fatalf("keygen: %v", err)
	}
	info, err := os.Stat(filepath.Join(dir, "private.pem"))
	if err != nil {
		t.Fatalf("private.pem missing: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("private.pem mode = %o, want 600", info.Mode().Perm())
	}

	// A second keygen must refuse to replace the key without --force.
	if _, err := runCLI(t, "", "jwt", "keygen", "--out-dir", dir, "--kid", "test-kid"); err == nil {
		t.Fatal("expected keygen to refuse overwriting an existing key")
	}

	jwksOut, err := runCLI(t, "", "jwt", "jwks", "--public-key", filepath.Join(dir, "public.pem"), "--kid", "test-kid")
	if err != nil {
		t.Fatalf("jwks: %v", err)
	}
	onDisk, err := os.ReadFile(filepath.Join(dir, "jwks.json"))
	if err != nil {
		t.Fatal(err)
	}
	if jwksOut != string(onDisk) {
		t.Errorf("jwks output differs from the jwks.json written by keygen")
	}

	token, err := runCLI(t, "", "jwt", "sign",
		"--private-key", filepath.Join(dir, "private.pem"), "--kid", "test-kid",
		"--issuer", "learning-istio@example.com", "--subject", "curl",
		"--claim", "app=curl", "--claim-list", "groups=readers,team-a")
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	token = strings.TrimSpace(token)
	if strings.Count(token, ".") != 2 {
		t.Fatalf("expected a compact JWT, got %q", token)
	}

	// Verified decode, token piped on stdin.
	decodedOut, err := runCLI(t, token+"\n", "jwt", "decode", "--token", "", "--public-key", filepath.Join(dir, "public.pem"))
	if err != nil {
		t.Fatalf("decode --public-key: %v", err)
	}
	var decoded struct {
		Header   map[string]any `json:"header"`
		Claims   map[string]any `json:"claims"`
		Verified bool           `json:"verified"`
	}
	if err := json.Unmarshal([]byte(decodedOut), &decoded); err != nil {
		t.Fatalf("decode output is not JSON: %v\n%s", err, decodedOut)
	}
	if !decoded.Verified {
		t.Error("expected verified=true when --public-key is given")
	}
	if decoded.Header["kid"] != "test-kid" || decoded.Claims["app"] != "curl" || decoded.Claims["sub"] != "curl" {
		t.Errorf("unexpected decode output: %s", decodedOut)
	}

	// Unverified decode must say so.
	unverifiedOut, err := runCLI(t, "", "jwt", "decode", "--token", token, "--public-key", "")
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !strings.Contains(unverifiedOut, `"verified": false`) {
		t.Errorf("expected verified=false without --public-key, got: %s", unverifiedOut)
	}

	// A tampered token must fail verification.
	tampered := token[:len(token)-4] + "AAAA"
	if _, err := runCLI(t, "", "jwt", "decode", "--token", tampered, "--public-key", filepath.Join(dir, "public.pem")); err == nil {
		t.Error("expected a tampered token to fail verification")
	}
}

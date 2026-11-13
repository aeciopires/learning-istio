package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/aeciopires/learning-istio/agents/upgrade/internal/jwttool"
)

var jwtCmd = &cobra.Command{
	Use:   "jwt",
	Short: "Generate keys and create, decode and verify JWTs for the Istio JWT walkthrough",
	Long: `jwt mints and inspects JSON Web Tokens for study clusters, using
github.com/golang-jwt/jwt. It covers what Istio's JWT policies need:

  keygen  an RSA key pair plus the JWKS for RequestAuthentication
  jwks    the JWKS of an existing public key
  sign    an RS256 token with iss/sub/aud and custom claims
  decode  a token's header and claims, optionally verifying the signature

Everything runs offline; nothing here talks to a cluster. The keys and
tokens are for learning environments only - never reuse them elsewhere.`,
}

var (
	jwtKeygenOutDir string
	jwtKeygenKid    string
	jwtKeygenBits   int
	jwtKeygenForce  bool
)

var jwtKeygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Generate an RSA key pair and its JWKS",
	Long: `keygen writes three files to --out-dir:

  private.pem  PKCS#8 private key, used by "jwt sign" (keep it out of Git)
  public.pem   PKIX public key, used by "jwt decode --public-key" and jwt.io
  jwks.json    JSON Web Key Set, pasted into RequestAuthentication's "jwks"

Existing files are never overwritten unless --force is given: replacing the
key invalidates every token already signed with it.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		paths := []string{
			filepath.Join(jwtKeygenOutDir, "private.pem"),
			filepath.Join(jwtKeygenOutDir, "public.pem"),
			filepath.Join(jwtKeygenOutDir, "jwks.json"),
		}
		if !jwtKeygenForce {
			for _, p := range paths {
				if _, err := os.Stat(p); err == nil {
					return fmt.Errorf("%s already exists; pass --force to replace the key (tokens signed with the old one stop validating)", p)
				} else if !errors.Is(err, os.ErrNotExist) {
					return err
				}
			}
		}

		key, err := jwttool.GenerateKey(jwtKeygenBits)
		if err != nil {
			return err
		}
		privPEM, err := jwttool.EncodePrivateKeyPEM(key)
		if err != nil {
			return err
		}
		pubPEM, err := jwttool.EncodePublicKeyPEM(&key.PublicKey)
		if err != nil {
			return err
		}
		jwks, err := jwttool.BuildJWKS(&key.PublicKey, jwtKeygenKid)
		if err != nil {
			return err
		}
		jwksJSON, err := jwttool.MarshalIndent(jwks)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(jwtKeygenOutDir, 0o700); err != nil {
			return fmt.Errorf("creating %s: %w", jwtKeygenOutDir, err)
		}
		files := []struct {
			path string
			data []byte
			mode os.FileMode
		}{
			{paths[0], privPEM, 0o600},
			{paths[1], pubPEM, 0o644},
			{paths[2], jwksJSON, 0o644},
		}
		for _, f := range files {
			if err := os.WriteFile(f.path, f.data, f.mode); err != nil {
				return fmt.Errorf("writing %s: %w", f.path, err)
			}
			// WriteFile keeps the mode of a file that already existed.
			if err := os.Chmod(f.path, f.mode); err != nil {
				return fmt.Errorf("setting mode on %s: %w", f.path, err)
			}
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "wrote %s\n", f.path); err != nil {
				return err
			}
		}
		return nil
	},
}

var (
	jwtJWKSPublicKey string
	jwtJWKSKid       string
)

var jwtJWKSCmd = &cobra.Command{
	Use:   "jwks",
	Short: "Print the JWKS for an RSA public key",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile(jwtJWKSPublicKey)
		if err != nil {
			return fmt.Errorf("reading %s: %w", jwtJWKSPublicKey, err)
		}
		key, err := jwttool.ParsePublicKeyPEM(data)
		if err != nil {
			return err
		}
		jwks, err := jwttool.BuildJWKS(key, jwtJWKSKid)
		if err != nil {
			return err
		}
		out, err := jwttool.MarshalIndent(jwks)
		if err != nil {
			return err
		}
		_, err = cmd.OutOrStdout().Write(out)
		return err
	},
}

var (
	jwtSignPrivateKey string
	jwtSignKid        string
	jwtSignIssuer     string
	jwtSignSubject    string
	jwtSignAudiences  []string
	jwtSignExpiresIn  time.Duration
	jwtSignClaims     []string
	jwtSignClaimLists []string
)

var jwtSignCmd = &cobra.Command{
	Use:   "sign",
	Short: "Create an RS256-signed JWT with custom claims",
	Long: `sign prints a compact JWT signed with --private-key.

Istio builds the request principal as "<iss>/<sub>", so --issuer and
--subject decide what AuthorizationPolicy's requestPrincipals matches.
Custom claims are matched with request.auth.claims[<name>]:

  --claim app=curl                  a string claim
  --claim-list groups=readers,ops   a list-of-strings claim

--kid must be the kid of the JWKS the verifier uses.`,
	Example: `  istio-upgrade jwt sign --private-key .jwt/private.pem \
    --issuer learning-istio@example.com --subject curl \
    --claim app=curl --claim-list groups=readers`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile(jwtSignPrivateKey)
		if err != nil {
			return fmt.Errorf("reading %s: %w", jwtSignPrivateKey, err)
		}
		key, err := jwttool.ParsePrivateKeyPEM(data)
		if err != nil {
			return err
		}
		claims, err := jwttool.ParseClaimFlags(jwtSignClaims, jwtSignClaimLists)
		if err != nil {
			return err
		}
		token, err := jwttool.Sign(key, jwttool.SignOptions{
			Kid:       jwtSignKid,
			Issuer:    jwtSignIssuer,
			Subject:   jwtSignSubject,
			Audiences: jwtSignAudiences,
			ExpiresIn: jwtSignExpiresIn,
			Claims:    claims,
		})
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(cmd.OutOrStdout(), token)
		return err
	},
}

var (
	jwtDecodeToken     string
	jwtDecodePublicKey string
)

var jwtDecodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "Print a JWT's header and claims, optionally verifying it",
	Long: `decode prints the token's header and claims as JSON. The token comes
from --token or, when that is empty, from standard input.

Without --public-key the token is only decoded: its signature and expiry
are NOT checked ("verified": false), exactly like pasting it into a JWT
debugger. With --public-key the RS256 signature and the time-based claims
are verified, and the command fails if either check does not pass.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		token := strings.TrimSpace(jwtDecodeToken)
		if token == "" {
			raw, err := io.ReadAll(cmd.InOrStdin())
			if err != nil {
				return fmt.Errorf("reading token from stdin: %w", err)
			}
			token = strings.TrimSpace(string(raw))
		}
		if token == "" {
			return errors.New("no token given: pass --token or pipe it on stdin")
		}

		var decoded jwttool.Decoded
		if jwtDecodePublicKey == "" {
			var err error
			if decoded, err = jwttool.Decode(token); err != nil {
				return err
			}
		} else {
			data, err := os.ReadFile(jwtDecodePublicKey)
			if err != nil {
				return fmt.Errorf("reading %s: %w", jwtDecodePublicKey, err)
			}
			key, err := jwttool.ParsePublicKeyPEM(data)
			if err != nil {
				return err
			}
			if decoded, err = jwttool.Verify(token, key); err != nil {
				return err
			}
		}
		out, err := jwttool.MarshalIndent(decoded)
		if err != nil {
			return err
		}
		_, err = cmd.OutOrStdout().Write(out)
		return err
	},
}

func init() {
	jwtKeygenCmd.Flags().StringVar(&jwtKeygenOutDir, "out-dir", ".jwt", "directory to write private.pem, public.pem and jwks.json to")
	jwtKeygenCmd.Flags().StringVar(&jwtKeygenKid, "kid", "learning-istio", "key ID recorded in the JWKS")
	jwtKeygenCmd.Flags().IntVar(&jwtKeygenBits, "bits", jwttool.MinKeyBits, "RSA key size in bits")
	jwtKeygenCmd.Flags().BoolVar(&jwtKeygenForce, "force", false, "replace an existing key pair")

	jwtJWKSCmd.Flags().StringVar(&jwtJWKSPublicKey, "public-key", ".jwt/public.pem", "RSA public key (PEM)")
	jwtJWKSCmd.Flags().StringVar(&jwtJWKSKid, "kid", "learning-istio", "key ID recorded in the JWKS")

	jwtSignCmd.Flags().StringVar(&jwtSignPrivateKey, "private-key", ".jwt/private.pem", "RSA private key (PEM) to sign with")
	jwtSignCmd.Flags().StringVar(&jwtSignKid, "kid", "learning-istio", "key ID written to the token header")
	jwtSignCmd.Flags().StringVar(&jwtSignIssuer, "issuer", "", "iss claim (required)")
	jwtSignCmd.Flags().StringVar(&jwtSignSubject, "subject", "", "sub claim (required)")
	jwtSignCmd.Flags().StringSliceVar(&jwtSignAudiences, "audience", nil, "aud claim; repeat or comma-separate for several")
	jwtSignCmd.Flags().DurationVar(&jwtSignExpiresIn, "expires-in", time.Hour, "token lifetime, e.g. 30m, 24h")
	jwtSignCmd.Flags().StringArrayVar(&jwtSignClaims, "claim", nil, "custom string claim as name=value; repeatable")
	jwtSignCmd.Flags().StringArrayVar(&jwtSignClaimLists, "claim-list", nil, "custom list claim as name=a,b,c; repeatable")
	_ = jwtSignCmd.MarkFlagRequired("issuer")
	_ = jwtSignCmd.MarkFlagRequired("subject")

	jwtDecodeCmd.Flags().StringVar(&jwtDecodeToken, "token", "", "token to decode (default: read from stdin)")
	jwtDecodeCmd.Flags().StringVar(&jwtDecodePublicKey, "public-key", "", "RSA public key (PEM); when set, the signature and expiry are verified")

	jwtCmd.AddCommand(jwtKeygenCmd, jwtJWKSCmd, jwtSignCmd, jwtDecodeCmd)
	rootCmd.AddCommand(jwtCmd)
}

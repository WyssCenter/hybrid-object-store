package auth

import (
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	errors "github.com/gigantum/hoss-error"
)

func JSONWebKeyFromRSA(key *rsa.PrivateKey) JSONWebKey {
	// The modulus and exponent arg represented in the JWK as a URL safe b64 encoded strings
	nEncoded := base64.RawURLEncoding.EncodeToString(key.N.Bytes())

	eBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(eBytes, uint64(key.E))
	eEncoded := base64.RawURLEncoding.EncodeToString(eBytes[5:8])

	// Kid is a unique case-sensitive string representing this key. We'll simply use
	// the SHA-1 hash of the encoded modulus string
	h := sha1.New()
	h.Write([]byte(nEncoded))
	bs := h.Sum(nil)
	keyID := fmt.Sprintf("%x", bs)

	return JSONWebKey{
		Alg: "RS256",
		Kty: "RSA",
		Kid: keyID,
		Use: "sig",
		E:   eEncoded,
		N:   nEncoded,
	}
}

// JSONWebKey represents public key info in a jwks.json file.
type JSONWebKey struct {
	Alg string `json:"alg"` // algorithm
	Kty string `json:"kty"` // algorithm family
	Kid string `json:"kid"` // unique key identifier
	Use string `json:"use"` // purpose of key
	N   string `json:"n"`   // RSA public key modulus
	E   string `json:"e"`   // RSA public key expoonent
	//X5c []string `json:"x5c"` // x509 cert chain; first entry is used for token verification
	//X5t string   `json:"X5t"` // x509 cert thumbprint

	public interface{}
}

// GetKey converts the JSONWebKey into the appropriate key object (ex rsa.PublicKey)
func (key JSONWebKey) GetKey() (interface{}, error) {
	if key.public != nil {
		return key.public, nil
	}

	switch key.Kty {
	case "RSA":
		eDecoded, err := base64.RawURLEncoding.DecodeString(key.E)
		if err != nil {
			return nil, errors.Wrap(err, "Invalid E value")
		}
		e := int(new(big.Int).SetBytes(eDecoded).Int64())

		nDecoded, err := base64.RawURLEncoding.DecodeString(key.N)
		if err != nil {
			return nil, errors.Wrap(err, "Invalid N value")
		}
		n := new(big.Int).SetBytes(nDecoded)

		key.public = &rsa.PublicKey{N: n, E: e}
	default:
		return nil, errors.New("Key algorithm not implemented")
	}

	return key.public, nil
}

// Jwks is a collection of public keys.
type Jwks struct {
	Keys []JSONWebKey `json:"keys"`
}

// GetSigningKey returns the first signing key that matches the Key ID
func (jwks Jwks) GetSigningKey(kid string) *JSONWebKey {
	for _, key := range jwks.Keys {
		if key.Kid == kid && key.Use == "sig" {
			return &key
		}
	}

	return nil
}

// GetJWKS downloads and parses the JWKS information from an OAuth / OIDC endpoint
func GetJWKS(url string) (Jwks, error) {
	jwks := Jwks{}
	resp, err := http.Get(url)

	if err != nil {
		return jwks, err
	}

	defer resp.Body.Close()
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return jwks, err
	}

	if resp.StatusCode != 200 {
		return jwks, errors.New("Could not retrieve JWKS information: " + resp.Status)
	}

	err = json.Unmarshal(responseBytes, &jwks)
	if err != nil {
		return jwks, err
	}

	return jwks, nil
}

package main

// This is a sample of numerous network wrappers useful for supporting
// common functionality, including CORS, BasicAuth, JWT Refresh and Access tokens
// without adding redundant checks into every function, and supporting restrictions
// at route registration.
//
// Both BasicAuth and requests for an AccessToken using a RefreshToken must use HTTPS or they
// will be insecure.  However, requests that use AccessAuth which simply checks the token signature
// are secure with or without HTTPS overhead.
//
// This assumes you are handling an OAuth style authentication in-house, otherwise if connecting to
// an external source then the routing may differ slightly in order to support them.
//
// Personally I favor ECC due to size and performance gains, but the NIST algorithms that are the basis of
// the ECDSA implementations are speculated to have been compromised by the NSA, while the ED25519
// is said to be more secure, and so I have included a reference implementation of the ED25519 with JWT.
//
// references:
//
// @link: https://stackoverflow.com/questions/46086746/how-do-you-specify-the-http-referrer-in-golang
// @link: https://medium.com/@matryer/the-http-handler-wrapper-technique-in-golang-updated-bc7fbcffa702
// @link: https://web.dev/sign-in-form-best-practices/
// @link: https://itnext.io/getting-started-with-oauth2-in-go-1c692420e03
// @link: https://blainsmith.com/articles/signing-jwts-with-gos-crypto-ed25519/
// @link: https://www.sohamkamani.com/blog/2018/03/25/golang-session-authentication/

import (
	"os"
	"log"
	"time"
	"errors"
	"strings"
	"context"
	"net/http"
	"crypto/subtle"
	"crypto/ed25519"
	"encoding/pem"
	crand "crypto/rand"


	"github.com/julienschmidt/httprouter"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var publicKey ed25519.PublicKey
var privateKey ed25519.PrivateKey
var pemPublicKey []byte
var pemPrivateKey []byte


// @note: the refresh_token is best left as a hash id which links to a database
// and not a full JWT.
const refreshToken = `{"refresh_token": "a3713565-1a5b-45ea-8ac0-5f49cc4c7017"}`

func PublicKey(w http.ResponseWriter, r *http.Request) {
	log.Printf("Public Key Requested: %s", r.Context().Value("uuid"))
	w.Header().Set("Content-Type", "application/x-pem-file")
	w.WriteHeader(http.StatusOK)
	w.Write(pemPublicKey)
}

// func DefaultOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
func DefaultOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "GET,PUT,POST,PATCH,DELETE,OPTIONS")
	w.WriteHeader(http.StatusOK)
}

// func CORSOptions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
func CORSOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "GET,PUT,POST,PATCH,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,PATCH,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Vary", "Origin")
	w.WriteHeader(http.StatusOK)
}

// Automatically print a log message for any request wrapped in the log method.
//
// This example only adds a uuid to help track logs across the system, but you
// could create a custom struct to track a bunch of common information and use
// that instead of a single uuid.
func Log(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		log.Printf("%s request to %s (%s)\n", r.Method, r.URL, id)
		h(w, r.WithContext(context.WithValue(r.Context(), "uuid", id)))
	})
}

// BasicAuth is a rudimentary example of login using username and password such as
// from a web form, and is expected to only ever be used over HTTPS.
//
// Normally you would encrypt and compare the username and password against
// stored encrypted fields, and return a JWT refresh token for subsequent requests.
//
// A slightly older but still common option is to set a cookie for a session that
// can be sent from all subsequent requests allowing them to be authenticated using
// the basic auth provided.  This requires at least one extra database call per operation
// to verify a session identifier, check expiration, and check permissions.
//
// In this particular example we are using hard-coded credentials and securing against
// timing based attacks by leveraging the subtle package ontop of string comparison.
//
// Finally, according to the specification, you can set the WWW-Authenticate
// header to send back a `Basic realm=`, which allows you to define a scope that
// would be used for restricting access.
func BasicAuth(w http.ResponseWriter, r *http.Request) {
    user, pass, ok := r.BasicAuth()
    if !ok || subtle.ConstantTimeCompare([]byte(user), []byte("exampleuser")) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte("examplepassword")) != 1 {
    	log.Printf("failed basic authentication: %s", r.Context().Value("uuid"))
        w.WriteHeader(http.StatusUnauthorized) // 401
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(refreshToken))
}

// Extract Refresh Token (id) from header and then generate an Access Token JWT-style
// using the ED25519 private key to sign it.
func AccessToken(w http.ResponseWriter, r *http.Request) {
	var rToken string
	if t := strings.Split(r.Header.Get("Authorization"), " "); len(t) == 2 {
		rToken = t[1]
	} else {
		log.Printf("No token found: %s", r.Context().Value("uuid"))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if rToken != "a3713565-1a5b-45ea-8ac0-5f49cc4c7017" {
		log.Printf("invalid token: %s", r.Context().Value("uuid"))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// define custom claims with 5 minute expiration
	claims := Claims{
		Permissions: []string{"user", "admin"},
		StandardClaims: jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
			Issuer: "system-name",
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
			Subject: "user-id",
		},
	}

	// generate and sign token
	token := jwt.NewWithClaims(&edDSASigningMethod, claims)

	log.Printf("Token: %#v\n", token)

	jwtstring, err := token.SignedString(privateKey)
	if err != nil {
		log.Printf("failed to create token: %s (%s)\n", err, r.Context().Value("uuid"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jwtstring))
}

// Extract the JWT from the headers, and then use the ED25519 public key
// to check the signature.
//
// Afterwards, check the metadata for perms (eg. permissions) with at least
// one match to what is provided in the strings.
func JWTAuth(h http.HandlerFunc, perms ...string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var jwtstring string
		if t := strings.Split(r.Header.Get("Authorization"), " "); len(t) == 2 {
			jwtstring = t[1]
		} else {
			log.Printf("No token found: %s", r.Context().Value("uuid"))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// @note: we can define the claims up front to save ourselves some trouble later
		// @note: while it may be possible to differentiate between validation errors the response
		// should always be a 401 so there is little point except for debugging.
		claims := &Claims{}
		if _, err := jwt.ParseWithClaims(jwtstring, claims, func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		}); err != nil {
			log.Printf("%#v", err)
			log.Printf("Failed to parse token: %s (%s)", err, r.Context().Value("uuid"))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// If they have valid credentials but not permissions then a 403 is expected
		if !claims.Can("admin") {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		h(w, r)
	})
}

// This is an example function secured by JWT access authentication
func Example(w http.ResponseWriter, r *http.Request) {
	log.Printf("This operation will only be reached if JWT authentication was successful: %s", r.Context().Value("uuid"))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success!"))
}

// Generate a keypair
//
// If you provide a nil reader it will default to crypto-rand but this
// is explicit for the sake of clarity.
func KeyGen() (err error) {
	publicKey, privateKey, err = ed25519.GenerateKey(crand.Reader)
	pemPublicKey = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: publicKey})
	pemPrivateKey = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privateKey})
	return
}

var ErrEdDSAVerification = errors.New("eddsa: verification error")

// Custom claims with permissions
type Claims struct {
	Permissions []string `json:"perms"`
	jwt.StandardClaims
}

// A basic abstraction that checks for a single match to return true.
func (c *Claims) Can(perms ...string) bool {
	for _, p := range c.Permissions {
		for _, v := range perms {
			if v == p {
				return true
			}
		}
	}
	return false
}

// @link: https://blainsmith.com/articles/signing-jwts-with-gos-crypto-ed25519/
type SigningMethodEdDSA struct{}

func (m *SigningMethodEdDSA) Alg() string {
	return "EdDSA"
}

func (m *SigningMethodEdDSA) Verify(signingString string, signature string, key interface{}) error {
	var err error

	var sig []byte
	if sig, err = jwt.DecodeSegment(signature); err != nil {
		return err
	}

	var ed25519Key ed25519.PublicKey
	var ok bool
	if ed25519Key, ok = key.(ed25519.PublicKey); !ok {
		return jwt.ErrInvalidKeyType
	}

	if len(ed25519Key) != ed25519.PublicKeySize {
		return jwt.ErrInvalidKey
	}

	if ok := ed25519.Verify(ed25519Key, []byte(signingString), sig); !ok {
		return ErrEdDSAVerification
	}

	return nil
}

func (m *SigningMethodEdDSA) Sign(signingString string, key interface{}) (str string, err error) {
	var ed25519Key ed25519.PrivateKey
	var ok bool
	if ed25519Key, ok = key.(ed25519.PrivateKey); !ok {
		return "", jwt.ErrInvalidKeyType
	}

	if len(ed25519Key) != ed25519.PrivateKeySize {
		return "", jwt.ErrInvalidKey
	}

	sig := ed25519.Sign(ed25519Key, []byte(signingString))
	return jwt.EncodeSegment(sig), nil
}

var edDSASigningMethod SigningMethodEdDSA

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lshortfile)
	jwt.RegisterSigningMethod(edDSASigningMethod.Alg(), func() jwt.SigningMethod { return &edDSASigningMethod })
}

// Demonstration of all behaviors with pre-defined routes for select examples
//
// Generate an ED25519 keypair for JWT processing.
//
// All routes carry logging, which adds a UUID to the context that can be used to track
// an operation as deep into the system as the context is passed.
//
// Add CORS dynamic support to all API routes.
//
// POST Basic Authentication handler.
//
// A route to request an access token (eg. JWT) using the refresh token.
//
// Finally an example of a secured route that demonstrates JWT Access Token validation.
func main() {
	if err := KeyGen(); err != nil {
		log.Printf("failed to generate keypair: %s\n", err)
		os.Exit(1)
	}

	router := httprouter.New()
	router.HandlerFunc(http.MethodOptions, "/", Log(DefaultOptions))
	router.HandlerFunc(http.MethodOptions, "/api/*all", Log(CORSOptions))
	router.HandlerFunc(http.MethodGet, "/key.pub", Log(PublicKey))
	router.HandlerFunc(http.MethodPost, "/api/login", Log(BasicAuth))
	router.HandlerFunc(http.MethodGet, "/api/access", Log(AccessToken))
	router.HandlerFunc(http.MethodGet, "/api/secure", Log(JWTAuth(Example, "admin")))
	http.ListenAndServe(":3000", router)
}

package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nazibul7/go-grpc-microservices/user-service/internal/model"
)

func VerifyToken(
	tokenString string,
	secret string,
) (*model.Claims, error) {

	// Empty struct where JWT library will store parsed claims.
	//
	// After successful verification:
	// claim.UserID
	// claim.Role
	// claim.ExpiresAt
	// etc...
	//
	// will be automatically populated.
	claim := model.Claims{}

	// ParseWithClaims does MANY things internally:
	//
	// 1. Split token into:
	//    header.payload.signature
	//
	// 2. Decode header + payload
	//
	// 3. Read algorithm from token header
	//
	// 4. Call our callback function below
	//
	// 5. Get secret key from callback
	//
	// 6. Recompute expected signature using:
	//    HMAC(header.payload, secret)
	//
	// 7. Compare expected signature
	//    with incoming signature
	//
	// 8. Validate standard claims:
	//    exp, iat, nbf, etc...
	//
	// 9. Populate `claim` struct
	token, err := jwt.ParseWithClaims(

		// Incoming JWT token from client
		tokenString,

		// Struct where decoded claims will be stored
		&claim,

		// Callback function.
		//
		// JWT library pauses parsing here and asks:
		//
		// "What key should I use to verify this token?"
		func(t *jwt.Token) (any, error) {

			// SECURITY CHECK
			//
			// Check whether incoming token claims
			// to use expected algorithm.
			//
			// Example valid:
			// HS256
			//
			// Example invalid:
			// RS256
			// none
			//
			// This prevents algorithm confusion attacks.
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}

			// Return server secret key.
			//
			// JWT library will NOW use this secret
			// internally to verify token signature.
			//
			// Conceptually:
			//
			// expectedSignature =
			//    HMAC(header.payload, secret)
			//
			// compare:
			// expectedSignature == incomingSignature
			// here we are returning secret so that invisible verify function can verify.
			return []byte(secret), nil
		},
	)

	// Verification failed:
	//
	// - invalid signature
	// - expired token
	// - malformed token
	// - wrong secret
	// - invalid claims
	if err != nil {
		return nil, err
	}

	// Extra safety check.
	//
	// token.Valid becomes true ONLY IF:
	//
	// - signature valid
	// - token not expired
	// - claims valid
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// At this point:
	//
	// Token is cryptographically verified.
	//
	// claim struct now contains:
	//
	// claim.UserID
	// claim.Role
	// claim.ExpiresAt
	// etc...
	return &claim, nil
}
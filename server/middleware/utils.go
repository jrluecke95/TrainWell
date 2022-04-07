package middleware

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/sessions"
)

var (
	store             = sessions.NewCookieStore([]byte(GoDotEnvVariable("SESSION_KEY")))
	CoachSessionName  = "coach-token"
	ClientSessionName = "client-token"
)

type Claims struct {
	jwt.StandardClaims
	Email string `json:"email"`
}

var jwtKey = []byte([]byte(GoDotEnvVariable("SESSION_KEY")))

func TokenCheck(res http.ResponseWriter, req *http.Request) error {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := req.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			res.WriteHeader(http.StatusUnauthorized)
			return errors.New("no cookie")
		}
		// For any other type of error, return a bad request status
		res.WriteHeader(http.StatusBadRequest)
		return errors.New("generic error with cookie")
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			res.WriteHeader(http.StatusUnauthorized)
			return nil
		}
		res.WriteHeader(http.StatusBadRequest)
		return errors.New("another error")
	}
	if !tkn.Valid {
		res.WriteHeader(http.StatusUnauthorized)
		return errors.New("token invalid")
	}
	return nil
}

func ConvertStringToInt(value string) int16 {
	initInt, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		panic(err)
	}

	finalInt := int16(initInt)
	return finalInt
}

//TODO fix the jwt util - cookie isn't stored or not being fetched correctly

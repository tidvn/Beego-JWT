package models

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/golang-jwt/jwt/v5"
)

func init() {
}

/**
 * @todo authenticate username & password using OAuth Grant Type Password
 * @todo add scopes and permissions to JWt
 * @todo open public apis only to unauthorized clients using Grant Type Client
 */
func AddToken(u User, d string) string {

	// current timestamp
	currentTimestamp := time.Now().UTC().Unix()
	var ttl int64 = 3600
	// md5 of sub & iat
	h := md5.New()
	io.WriteString(h, (u.ID.Hex()))
	io.WriteString(h, strconv.FormatInt(int64(currentTimestamp), 10))
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   u.ID.Hex(),
		"admin": u.IsAdmin,
		"iat":   currentTimestamp,
		"exp":   currentTimestamp + ttl,
		"nbf":   currentTimestamp,
		"iss":   d,
		"jti":   h.Sum(nil),
	})

	// Sign and get the complete encoded token as a string using the secret
	HMACKEY, _ := beego.AppConfig.String("HMACKEY")
	tokenString, err := token.SignedString([]byte(HMACKEY))

	if err != nil {
		log.Fatal(err)
	}

	return (tokenString)
}

func ParseToken(tokenString string) map[string]string {

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		HMACKEY, _ := beego.AppConfig.String("HMACKEY")
		return []byte(HMACKEY), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && claims != nil {
		uid := claims["sub"].(string)
		admin := claims["admin"].(string)
		return map[string]string{
			"uid":   uid,
			"admin": admin,
		}
	}
	return map[string]string{}
}

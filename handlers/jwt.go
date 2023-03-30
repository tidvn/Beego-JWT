package handlers

import (
	"DPay/models"
	"encoding/json"
	"fmt"

	"strings"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/golang-jwt/jwt/v5"
)

func init() {
}

func Jwt(ctx *context.Context) {
	ctx.Output.Header("Content-Type", "application/json")
	var uri string = ctx.Input.URI()
	fmt.Print(strings.HasPrefix(uri, "/v1/user/login"))
	if strings.HasPrefix(uri, "/v1/user/login") {
		return
	}
	if strings.HasPrefix(uri, "/v1/user/register") {
		return
	}
	if strings.HasPrefix(uri, "/v1/user/logout") {
		return
	}

	if ctx.Input.Header("Authorization") == "" {
		ctx.Output.SetStatus(403)
		resBody, err := json.Marshal(models.APIResponse{Status: 403, Message: "Failed to process request", Errors: "Undefined authorization"})
		ctx.Output.Body(resBody)
		if err != nil {
			panic(err)
		}
	}

	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	var tokenString string = ctx.Input.Header("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		HMACKEY, _ := beego.AppConfig.String("HMACKEY")
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(HMACKEY), nil
	})

	if err != nil {
		ctx.Output.SetStatus(403)
		var responseBody models.APIResponse = models.APIResponse{Status: 403, Message: "Failed to process request", Errors: err.Error()}
		resBytes, err := json.Marshal(responseBody)
		ctx.Output.Body(resBytes)
		if err != nil {
			panic(err)
		}
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && claims != nil {
		return
	} else {
		ctx.Output.SetStatus(403)
		resBody, err := json.Marshal(models.APIResponse{Status: 403, Message: ctx.Input.Header("Authorization")})
		ctx.Output.Body(resBody)
		if err != nil {
			panic(err)
		}
	}
}

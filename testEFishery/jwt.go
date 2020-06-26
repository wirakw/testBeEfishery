package main

import (
	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
)

// generate token to use.

func myAuthenticatedHandler(ir iris.Context) {
	jwtToken := ir.Values().Get("jwt").(*jwt.Token)
	foobar := jwtToken.Claims.(jwt.MapClaims)
	ir.Values().Set("role", foobar["role"])
	ir.Next()
}

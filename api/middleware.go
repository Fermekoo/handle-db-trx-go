package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Fermekoo/handle-db-tx-go/token"
	"github.com/gin-gonic/gin"
)

const (
	authorization_header_key  = "Authorization"
	authorization_header_type = "bearer"
	authorization_payload_key = "auth_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorization_header := ctx.GetHeader(authorization_header_key)
		if len(authorization_header) < 1 {
			err := errors.New("acces token is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorization_header)
		if len(fields) < 2 {
			err := errors.New("invalid access token format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		if authorization_header_type != strings.ToLower(fields[0]) {
			err := fmt.Errorf("unsupported authorization type %s", authorization_header_key)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		access_token := fields[1]
		payload, err := tokenMaker.VerifyToken(access_token)
		if err != nil {

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorization_payload_key, payload)
		ctx.Next()
	}

}

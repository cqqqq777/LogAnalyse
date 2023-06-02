package middleware

import (
	"context"
	"errors"
	"net/http"

	"LogAnalyse/app/shared/consts"
	"LogAnalyse/app/shared/errz"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/golang-jwt/jwt"
)

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
	TokenNotFound    = errors.New("no token")
)

func JwtAuth(secretKey string) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// get token
		token := c.Query("token")

		// check token
		if token == "" {
			c.JSON(http.StatusOK, utils.H{
				"code": errz.CodeInvalidParam,
				"msg":  TokenNotFound.Error(),
			})
		}
		j := NewJWT(secretKey)
		claim, err := j.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, utils.H{
				"code": errz.CodeTokenInvalid,
				"msg":  err.Error(),
			})
			c.Abort()
			return
		}

		//set context
		c.Set(consts.AccountID, claim.ID)
		c.Set(consts.AccountIdentity, claim.Identity)
		c.Next(ctx)
	}
}

type JWT struct {
	SigningKey []byte
}

type CustomClaims struct {
	ID       int64
	Identity string
	jwt.StandardClaims
}

func NewJWT(secretKey string) *JWT {
	return &JWT{
		SigningKey: []byte(secretKey),
	}
}

// CreateToken to create a token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParseToken to parse a token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	} else {
		return nil, TokenInvalid
	}
}

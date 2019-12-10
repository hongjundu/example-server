package server

import (
	"example-server/internal/app/example-server/storage/db"
	"example-server/internal/pkg/consts"
	"example-server/pkg/utils"
	"fmt"
	"github.com/casbin/casbin"
	"github.com/casbin/xorm-adapter"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/hongjundu/go-level-logger"
	"github.com/hongjundu/go-rest-api-helper.v1"
	"net/http"
	"strings"
)

func (server *Server) loadJwtKeys() (err error) {
	if server.jwtPublicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(consts.JWTPubKeyString)); err != nil {
		logger.Errorf("[Server] ReadJWTPublicKey failed, %v", err)
		return
	} else {
		logger.Infof("[Server] ReadJWTPublicKey successfully")
	}

	if server.jwtPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(consts.JWTPrivateKeyString)); err != nil {
		logger.Errorf("[Server] ReadJWTPrivateKey failed, %v", err)
		return
	} else {
		logger.Infof("[Server] ReadJWTPrivateKey successfully")
	}

	return
}

func (server *Server) createEnforcer() error {
	m := casbin.NewModel(consts.CosbinConf)

	a := xormadapter.NewAdapter("mysql", db.GetConnString(), true)

	server.enforcer = casbin.NewEnforcer()
	server.enforcer.InitWithModelAndAdapter(m, a)

	err := server.enforcer.LoadPolicy()

	return err
}

func (server *Server) jwtTokenRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		var shortToken string
		var jwtTokenString string
		var jwtToken *jwt.Token
		var err error

		var customClaims struct {
			jwt.StandardClaims
			User string `json:"user"`
		}

		authHeader := c.GetHeader("Authorization")
		if len(authHeader) == 0 {
			c.JSON(http.StatusUnauthorized, apihelper.NewErrorResponse(apihelper.NewError(http.StatusUnauthorized, "No Authorization header")))
			goto ABORT
		}

		if strings.HasPrefix(authHeader, "bearer") || strings.HasPrefix(authHeader, "Bearer") {
			splits := strings.Split(authHeader, " ")
			if len(splits) == 2 {
				shortToken = splits[1]
				if jwtTokenString, err = server.getJwtToken(shortToken); err != nil || len(jwtTokenString) == 0 {
					c.JSON(http.StatusUnauthorized, apihelper.NewErrorResponse(apihelper.NewError(http.StatusUnauthorized, "Wrong token or token expires")))
					goto ABORT
				}
			} else {
				c.JSON(http.StatusUnauthorized, apihelper.NewErrorResponse(apihelper.NewError(http.StatusUnauthorized, "Invalid Authorization header")))
				goto ABORT
			}

		} else {
			c.JSON(http.StatusUnauthorized, apihelper.NewErrorResponse(apihelper.NewError(http.StatusUnauthorized, "Invalid Authorization header")))
			goto ABORT
		}

		jwtToken, err = jwt.ParseWithClaims(jwtTokenString, &customClaims,
			func(token *jwt.Token) (interface{}, error) {
				return server.jwtPublicKey, nil
			})

		if err == nil {
			if jwtToken.Valid {
				c.Set("token", shortToken)
				c.Set("user", customClaims.User)
				c.Next()
				return
			} else {
				c.JSON(http.StatusUnauthorized, apihelper.NewErrorResponse(apihelper.NewError(http.StatusUnauthorized, "Invalid Token")))
				goto ABORT
			}

		} else {
			c.JSON(http.StatusUnauthorized, apihelper.NewErrorResponse(apihelper.NewError(http.StatusUnauthorized, "Unauthorized")))
			goto ABORT
		}

	ABORT:
		c.Abort()
	}
}

func (server *Server) acl(obj string, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		sub := c.GetString("user")
		if len(sub) == 0 {
			c.JSON(http.StatusUnauthorized, apihelper.NewErrorResponse(apihelper.NewError(http.StatusUnauthorized, "Unauthorized: No user")))
			c.Abort()
			return
		}

		ok := server.enforcer.Enforce(sub, obj, act)

		if !ok {
			c.JSON(http.StatusForbidden, apihelper.NewErrorResponse(apihelper.NewError(http.StatusForbidden, "StatusForbidden")))
			c.Abort()
			return
		}

		c.Next()
	}
}

func (server *Server) genShortToken(jwtToken string) (shortToken string, err error) {
	shortToken = utils.Sha256Encode(jwtToken)
	err = server.bigCache.Set(fmt.Sprintf("token-%s", shortToken), []byte(jwtToken))
	return
}

func (server *Server) getJwtToken(shortToken string) (jwtToken string, err error) {
	if bytesJwtToken, e := server.bigCache.Get(fmt.Sprintf("token-%s", shortToken)); e == nil {
		jwtToken = string(bytesJwtToken)
	} else {
		err = e
	}
	return
}

func (server *Server) clearShortToken(shortToken string) error {
	return server.bigCache.Delete(fmt.Sprintf("token-%s", shortToken))
}

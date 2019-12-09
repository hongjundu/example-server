package server

import (
	"example-server/internal/app/example-server/storage/db"
	"example-server/internal/pkg/consts"
	"github.com/casbin/casbin"
	//"github.com/casbin/casbin/model"
	"github.com/casbin/xorm-adapter"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"github.com/hongjundu/go-level-logger"
	"github.com/hongjundu/go-rest-api-helper.v1"
	"net/http"
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

func (server *Server) tokenRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		var customClaims struct {
			jwt.StandardClaims
			User string `json:"user"`
		}

		token, err := request.ParseFromRequestWithClaims(c.Request, request.AuthorizationHeaderExtractor, &customClaims,
			func(token *jwt.Token) (interface{}, error) {
				return server.jwtPublicKey, nil
			})

		if err == nil {
			if token.Valid {
				c.Set("user", customClaims.User)
				c.Next()
			} else {
				c.JSON(http.StatusUnauthorized, apihelper.NewErrorResponse(apihelper.NewError(http.StatusUnauthorized, "Invalid Token")))
				c.Abort()
			}

		} else {
			c.JSON(http.StatusUnauthorized, apihelper.NewErrorResponse(apihelper.NewError(http.StatusUnauthorized, "Unauthorized")))
			c.Abort()
		}
	}
}

func (server *Server) acl(obj string, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		sub := c.GetString("user")
		if len(sub) == 0 {
			c.JSON(http.StatusUnauthorized, apihelper.NewErrorResponse(apihelper.NewError(http.StatusUnauthorized, "Unauthorized: No user")))
			c.Abort()
		}

		ok := server.enforcer.Enforce(sub, obj, act)

		if !ok {
			c.JSON(http.StatusForbidden, apihelper.NewErrorResponse(apihelper.NewError(http.StatusForbidden, "StatusForbidden")))
			c.Abort()
		}

		c.Next()
	}
}

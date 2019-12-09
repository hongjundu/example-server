package server

import (
	"example-server/pkg/version"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/hongjundu/go-rest-api-helper.v1"
	"net/http"
	"time"
)

func ginHandlerFunc(f func(c *gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {

		rsp, err := f(c)

		if err == nil {
			c.JSON(http.StatusOK, rsp)
		} else {
			code := http.StatusInternalServerError

			if apiErr, ok := err.(apihelper.ApiError); ok {
				code = apiErr.Code()
			}

			c.JSON(code, apihelper.NewErrorResponse(err))
		}
	}
}

func (server *Server) notFoundHandler(c *gin.Context) (response interface{}, err error) {
	err = apihelper.NewError(http.StatusNotImplemented, "not found")
	return
}

func (server *Server) helloHandler(c *gin.Context) (response interface{}, err error) {

	response = apihelper.NewOKResponse("hello " + c.GetString("user"))
	return
}

func (server *Server) loginHandler(c *gin.Context) (response interface{}, err error) {
	var param struct {
		User     string `form:"user" json:"user" binding:"required"`
		Password string `form:"password" json:"password" binding:"required"`
	}

	if e := c.ShouldBindJSON(&param); e != nil {
		err = apihelper.NewError(http.StatusBadRequest, e.Error())
		return
	}

	claims := jwt.MapClaims{}

	now := time.Now()
	exp := now.AddDate(0, 0, 1)

	claims["iat"] = now.Unix()
	claims["exp"] = exp.Unix()
	claims["user"] = param.User

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	if tokens, e := token.SignedString(server.jwtPrivateKey); e == nil {
		response = apihelper.NewOKResponse(gin.H{"token": tokens, "user": param.User})
	} else {
		err = apihelper.NewError(http.StatusInternalServerError, e.Error())
	}

	return
}

func (server *Server) versionHandler(c *gin.Context) (response interface{}, err error) {
	var res struct {
		Version   string `json:"version"`
		GOVersion string `json:"goVersion"`
		BuildTime string `json:"buildTime"`
		BuildHost string `json:"buildHost"`
	}

	res.Version, res.GOVersion, res.BuildTime, res.BuildHost = version.Version, version.GOVersion, version.BuildTime, version.BuildHost

	response = apihelper.NewOKResponse(res)

	return
}

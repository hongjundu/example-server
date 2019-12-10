package server

import (
	"example-server/internal/app/example-server/server/apimodel"
	"example-server/pkg/version"
	"fmt"
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

func (server *Server) readHandler(c *gin.Context) (response interface{}, err error) {
	response = apihelper.NewOKResponse("read by: " + c.GetString("user"))
	return
}

func (server *Server) writeHandler(c *gin.Context) (response interface{}, err error) {
	response = apihelper.NewOKResponse("write by:  " + c.GetString("user"))
	return
}

// @Tags login
// @Summary user login
// @Param loginParam body apimodel.LoginParam true "login param"
// @accept application/json
// @Produce application/json
// @Success 200
// @Success 401
// @Router /login [post]
func (server *Server) loginHandler(c *gin.Context) (response interface{}, err error) {
	var param apimodel.LoginParam

	if e := c.ShouldBindJSON(&param); e != nil {
		err = apihelper.NewError(http.StatusBadRequest, e.Error())
		return
	}

	// TODO: verify user & password from DB ...

	claims := jwt.MapClaims{}

	now := time.Now()
	exp := now.AddDate(0, 0, 1)

	claims["iat"] = now.Unix()
	claims["exp"] = exp.Unix()
	claims["user"] = param.User

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	if tokens, e := token.SignedString(server.jwtPrivateKey); e == nil {
		if shortToken, e2 := server.genShortToken(tokens); e2 == nil {
			response = apihelper.NewOKResponse(gin.H{"token": shortToken})
		} else {
			err = e2
		}

	} else {
		err = apihelper.NewError(http.StatusInternalServerError, e.Error())
	}

	return
}

// @Tags logout
// @Summary user logout
// @Param Authorization header string true "bearer token"
// @accept application/json
// @Produce application/json
// @Success 200
// @Success 401
// @Router /logout [post]
func (server *Server) logoutHandler(c *gin.Context) (response interface{}, err error) {

	shortToken := c.GetString("token")
	if len(shortToken) == 0 {
		err = apihelper.NewError(http.StatusUnauthorized, "No token")
		return
	}

	user := c.GetString("user")

	err = server.clearShortToken(shortToken)
	response = apihelper.NewOKResponse(fmt.Sprintf("%s: logged out", user))

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

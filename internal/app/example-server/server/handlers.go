package server

import (
	"example-server/pkg/version"
	"github.com/gin-gonic/gin"
	"github.com/hongjundu/go-rest-api-helper.v1"
	"net/http"
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

func (server *Server) createTaskHandler(c *gin.Context) (response interface{}, err error) {
	response = apihelper.NewOKResponse(nil)
	return
}

func (server *Server) loginHandler(c *gin.Context) (response interface{}, err error) {
	response = apihelper.NewOKResponse(nil)
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

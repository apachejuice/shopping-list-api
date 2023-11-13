package api

import (
	"net/http"
	"strings"
	"time"

	"apachejuice.dev/apachejuice/shopping-list-api/internal/apispec"
	"apachejuice.dev/apachejuice/shopping-list-api/internal/logging"
	"github.com/gin-gonic/gin"
	"github.com/palantir/stacktrace"
)

// Implements the ServerInterface
type ApiImpl struct {
	auth     Authenticator
	delegate ApiDelegate
}

func NewApiImpl(auth Authenticator, delegate ApiDelegate) ApiImpl {
	return ApiImpl{
		auth:     auth,
		delegate: delegate,
	}
}

func (a *ApiImpl) Run(addr string) {
	r := gin.Default()

	// Set up swagger UI
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/openapi")
	})
	r.Static("/openapi", "./swaggerui/dist")

	apispec.RegisterHandlers(r, a)
	r.Run(addr)
}

func userError(msg string) apispec.UserError {
	return apispec.UserError{
		Time:        time.Now().UTC(),
		UserMessage: msg,
	}
}

func serverError(err error, msg string) apispec.ServerError {
	return apispec.ServerError{
		ErrorId:      logging.GenErrCode(),
		ErrorMessage: err.Error(),
		Time:         time.Now().UTC(),
		UserMessage:  msg,
	}
}

func (a *ApiImpl) guard(c *gin.Context) (bool, int, any) {
	token := strings.Trim(strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")[1], " ")
	if token == "" {
		return false, http.StatusUnauthorized, userError("No token provided")
	}

	ok, err := a.auth.CheckToken(token)
	if err != nil && err.IsServerError() {
		serr := serverError(stacktrace.RootCause(err), "Failed to validate token")
		logging.Error(err, serr.ErrorId)

		return false, http.StatusInternalServerError, serr
	} else if !ok {
		return false, http.StatusUnauthorized, userError("Expired or invalid token")
	}

	return true, -1, nil // endpoint defines status
}

var _ apispec.ServerInterface = (*ApiImpl)(nil)

func (a *ApiImpl) GetLists(c *gin.Context) {
	if ok, status, errObj := a.guard(c); !ok {
		c.JSON(status, errObj)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (a *ApiImpl) GetListsId(c *gin.Context, id int) {}

func (a *ApiImpl) GetMe(c *gin.Context) {}

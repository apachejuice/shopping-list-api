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
	delegate apiDelegate
}

func NewApiImpl(auth Authenticator) ApiImpl {
	return ApiImpl{
		auth:     auth,
		delegate: apiDelegate{},
	}
}

func (a *ApiImpl) Run(addr string, trustedProxies []string) {
	r := gin.Default()
	r.SetTrustedProxies(trustedProxies)

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

func (a *ApiImpl) handleApiError(c *gin.Context, aerr *ApiError) {
	if aerr.IsServerError() {
		serr := serverError(aerr, "Internal server error")
		logging.Error(aerr, serr.ErrorId)

		c.JSON(http.StatusInternalServerError, serr)
	} else if aerr.GetCode() != -1 {
		uerr := userError(aerr.Error())
		c.JSON(aerr.GetCode(), uerr)
	} else {
		// Return as 400
		uerr := userError("Malformed request")
		c.JSON(http.StatusBadRequest, uerr)
	}
}

func (a *ApiImpl) guard(c *gin.Context, setUserId *string) (bool, int, any) {
	token := strings.Trim(strings.Split(c.Request.Header.Get("Authorization"), "Bearer ")[1], " ")
	if token == "" {
		return false, http.StatusUnauthorized, userError("No token provided")
	}

	ok, userId, err := a.auth.CheckToken(c, token)
	if err != nil && err.IsServerError() {
		serr := serverError(stacktrace.RootCause(err), "Failed to validate token")
		logging.Error(err, serr.ErrorId)

		return false, http.StatusInternalServerError, serr
	} else if !ok {
		return false, http.StatusUnauthorized, userError("Expired or invalid token")
	}

	if setUserId != nil {
		*setUserId = userId
	}

	return true, -1, nil // endpoint defines status
}

var _ apispec.ServerInterface = (*ApiImpl)(nil)

func (a *ApiImpl) GetLists(c *gin.Context) {
	if ok, status, errObj := a.guard(c, nil); !ok {
		c.JSON(status, errObj)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (a *ApiImpl) GetListsId(c *gin.Context, id int) {}

func (a *ApiImpl) GetMe(c *gin.Context) {
	var userId string
	if ok, status, errObj := a.guard(c, &userId); !ok {
		c.JSON(status, errObj)
		return
	}

	me, aerr := a.delegate.getMe(c, userId)
	if aerr != nil {
		a.handleApiError(c, aerr)
		return
	}

	c.JSON(http.StatusOK, me)
}

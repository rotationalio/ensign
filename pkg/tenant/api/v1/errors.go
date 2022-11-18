package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	unsuccessful         = Reply{Success: false}
	notFound             = Reply{Success: false, Error: "resource not found"}
	notAllowed           = Reply{Success: false, Error: "method not allowed"}
	ErrAPIKeyIDRequired  = errors.New("apikey id is required for this endpoint")
	ErrMemberIDRequired  = errors.New("member id is required for this endpoint")
	ErrProjectIDRequired = errors.New("project id is required for this endpoint")
	ErrTenantIDRequired  = errors.New("tenant id is required for this endpoint")
	ErrTopicIDRequired   = errors.New("topic id is required for this endpoint")
)

// Constructs a new response for an error or returns unsuccessful.
func ErrorResponse(err interface{}) Reply {
	if err == nil {
		return unsuccessful
	}

	rep := Reply{Success: false}
	switch err := err.(type) {
	case error:
		rep.Error = err.Error()
	case string:
		rep.Error = err
	case fmt.Stringer:
		rep.Error = err.String()
	case json.Marshaler:
		data, e := err.MarshalJSON()
		if e != nil {
			panic(err)
		}
		rep.Error = string(data)
	default:
		rep.Error = "unhandled error response"
	}

	return rep
}

// NotFound returns a JSON response for the API.
// NOTE: we know it's weird to put server-side handlers like NotFound and NotAllowed
// here in the client/api side package but it unifies where we keep our error handling
// mechanisms.
func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, notFound)
}

// NotAllowed returns a JSON 405 response for the API.
func NotAllowed(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, notAllowed)
}

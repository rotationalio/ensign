package tenant

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
)

func (s *Server) UserList(c *gin.Context) {
	user := &api.UserPage{}
	c.JSON(http.StatusOK, user)
}

func (s *Server) UserCreate(c *gin.Context) {
	var newUser *api.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		return
	}

	c.JSON(http.StatusOK, newUser)
}

func (s *Server) UserDetail(c *gin.Context) {
	// Authorization and Authentication happen in middleware but may add data to the
	// context, for example the user, permissions, organization, etc.
	// Step 0: Perform any final checks or fetch middleware data from the context

	// Step 1: Load the request (either params from GET or body from POST) using c.Bind
	// c.BindJSON, etc. and validate it, returning BadRequest if it's invalid.
	user := &api.User{}
	if err := c.ShouldBindUri(user); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Step 2: Perform the work on the request, e.g. database interactions, sending
	// notifications, accessing other services, etc.

	// Step 3: Prepare the response

	// Step 4: Return the response with the correct status code
	c.JSON(http.StatusOK, user)
}

func (s *Server) UserUpdate(c *gin.Context) {
	var updateUser *api.User

	if err := c.ShouldBindUri(&updateUser); err != nil {
		return
	}

	c.JSON(http.StatusOK, updateUser)
}

func (s *Server) UserDelete(c *gin.Context) {
	deleteUser := &api.User{}
	if err := c.ShouldBindUri(&deleteUser); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, deleteUser)
}

func (s *Server) AppList(c *gin.Context) {
	user := &api.AppPage{}
	c.JSON(http.StatusOK, user)
}

func (s *Server) AppCreate(c *gin.Context) {
	var newApp *api.App

	if err := c.ShouldBindUri(&newApp); err != nil {
		return
	}

	c.JSON(http.StatusOK, newApp)
}

func (s *Server) AppDetail(c *gin.Context) {
	// Authorization and Authentication happen in middleware but may add data to the
	// context, for example the user, permissions, organization, etc.
	// Step 0: Perform any final checks or fetch middleware data from the context

	// Step 1: Load the request (either params from GET or body from POST) using c.Bind
	// c.BindJSON, etc. and validate it, returning BadRequest if it's invalid.
	app := &api.App{}
	if err := c.ShouldBindUri(app); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Step 2: Perform the work on the request, e.g. database interactions, sending
	// notifications, accessing other services, etc.

	// Step 3: Prepare the response

	// Step 4: Return the response with the correct status code
	c.JSON(http.StatusOK, app)
}

func (s *Server) TopicList(c *gin.Context) {
	user := &api.TopicPage{}
	c.JSON(http.StatusOK, user)
}

func (s *Server) AppDelete(c *gin.Context) {
	deleteApp := &api.App{}
	if err := c.ShouldBindUri(&deleteApp); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, deleteApp)
}

func (s *Server) TopicCreate(c *gin.Context) {
	var newTopic *api.Topic

	if err := c.ShouldBindUri(&newTopic); err != nil {
		return
	}

	c.JSON(http.StatusOK, newTopic)
}

func (s *Server) TopicDetail(c *gin.Context) {
	// Authorization and Authentication happen in middleware but may add data to the
	// context, for example the user, permissions, organization, etc.
	// Step 0: Perform any final checks or fetch middleware data from the context

	// Step 1: Load the request (either params from GET or body from POST) using c.Bind
	// c.BindJSON, etc. and validate it, returning BadRequest if it's invalid.
	topic := &api.Topic{}
	if err := c.ShouldBindUri(topic); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Step 2: Perform the work on the request, e.g. database interactions, sending
	// notifications, accessing other services, etc.

	// Step 3: Prepare the response

	// Step 4: Return the response with the correct status code
	c.JSON(http.StatusOK, topic)
}

func (s *Server) TopicDelete(c *gin.Context) {
	deleteTopic := &api.Topic{}
	if err := c.ShouldBindUri(&deleteTopic); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, deleteTopic)
}

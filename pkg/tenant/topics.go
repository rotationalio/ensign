package tenant

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rs/zerolog/log"
)

func (s *Server) ProjectTopicList(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func (s *Server) ProjectTopicCreate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func (s *Server) TopicList(c *gin.Context) {
	// The following TODO task items will need to be
	// implemented for each endpoint.

	// TODO: Add authentication and authorization middleware
	// TODO: Identify top-level info
	// TODO: Parse and validate user input
	// TODO: Perform work on the request, e.g. database interactions,
	// sending notifications, accessing other services, etc.

	// Return response with the correct status code

	// TODO: Replace StatusNotImplemented with StatusOk and
	// replace "not yet implemented" message.
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

func (s *Server) TopicCreate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

// TopicDetail retrieves a summary detail of a topic with a given ID
// and returns a 200 OK response.
func (s *Server) TopicDetail(c *gin.Context) {
	var (
		err   error
		reply *api.Topic
	)

	// Get the topic ID from the URL and return a 400 response
	// if the topic does not exist.
	//
	// Route: /topic/:topicID
	var topicID ulid.ULID
	if topicID, err = ulid.Parse(c.Param("topicID")); err != nil {
		log.Error().Err(err).Msg("could not parse topic ulid")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse topic ulid"))
		return
	}

	// Get the specified topic from the database and return a 404 response
	// if it cannot be retrieved.
	var topic *db.Topic
	if topic, err = db.RetrieveTopic(c.Request.Context(), topicID); err != nil {
		log.Error().Err(err).Str("topicID", topicID.String()).Msg("could not retrieve topic")
		c.JSON(http.StatusNotFound, api.ErrorResponse("could not retrieve topic"))
		return
	}

	reply = &api.Topic{
		ID:   topic.ID.String(),
		Name: topic.Name,
	}

	c.JSON(http.StatusOK, reply)
}

// TopicUpdate updates the record of a topic with a given ID and
// returns a 200 OK response.
//
// Route: /topic/:topicID
func (s *Server) TopicUpdate(c *gin.Context) {
	var (
		err   error
		topic *api.Topic
	)

	// Get the topic ID from the URL and return a 400 response if
	// the topic ID is not a ULID.
	var topicID ulid.ULID
	if topicID, err = ulid.Parse(c.Param("topicID")); err != nil {
		log.Error().Err(err).Msg("could not parse topic ulid")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse topic ulid"))
		return
	}

	// Bind the user request with JSON and return a 400 response if
	// binding is not successful.
	if err = c.BindJSON(&topic); err != nil {
		log.Warn().Err(err).Msg("could not parse topic update request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not bind user request"))
		return
	}

	// Verify the topic name exists and return a 400 response if it doesn't.
	if topic.Name == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("topic name is required"))
		return
	}

	// Get the specified topic from the database and return a 404 response if
	// it cannot be retrieved.
	var t *db.Topic
	if t, err = db.RetrieveTopic(c.Request.Context(), topicID); err != nil {
		log.Error().Err(err).Str("topicID", topicID.String()).Msg("could not retrieve topic")
		c.JSON(http.StatusNotFound, api.ErrorResponse("topic not found"))
		return
	}

	// Update topic in the database and return a 500 response if the topic
	// record cannot be updated.
	if err = db.UpdateTopic(c.Request.Context(), t); err != nil {
		log.Error().Err(err).Str("topicID", topicID.String()).Msg("could not save topic")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update topic"))
		return
	}

	c.JSON(http.StatusOK, topic)
}

// TopicDelete deletes a topic from a user's request with a given ID
// and returns a 200 OK response instead of an error response.
//
// Route: /topic/:topicID
func (s *Server) TopicDelete(c *gin.Context) {
	var (
		err error
	)

	// Get the topic ID from the URL and return a 400 response
	// if the topic does not exist.
	var topicID ulid.ULID
	if topicID, err = ulid.Parse(c.Param("topicID")); err != nil {
		log.Error().Err(err).Msg("could not parse topic ulid")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse topic ulid"))
		return
	}

	// Delete the topic and return a 404 response if it cannot be removed.
	if err = db.DeleteTopic(c.Request.Context(), topicID); err != nil {
		log.Error().Err(err).Str("topicID", topicID.String()).Msg("could not delete topic")
		c.JSON(http.StatusNotFound, api.ErrorResponse("could not delete topic"))
		return
	}

	c.Status(http.StatusOK)
}

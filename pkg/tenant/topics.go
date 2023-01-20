package tenant

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	middleware "github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rs/zerolog/log"
)

// ProjectTopicList retrieves all topics assigned to a project
// and returns a 200 OK response.
//
// Route: /projects/:projectID/topics
func (s *Server) ProjectTopicList(c *gin.Context) {
	var (
		err error
	)

	// Get the topic's project ID from the URL and return a 400 response
	// if the project ID is not a ULID.
	var projectID ulid.ULID
	if projectID, err = ulid.Parse(c.Param("projectID")); err != nil {
		log.Error().Err(err).Msg("could not parse project ulid")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse project ulid"))
		return
	}

	// Get topics from the database and return a 500 response
	// if not successful.
	var topics []*db.Topic
	if topics, err = db.ListTopics(c.Request.Context(), projectID); err != nil {
		log.Error().Err(err).Msg("could not fetch topics from the database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not fetch topics from the database"))
		return
	}

	// Build the response.
	out := &api.ProjectTopicPage{
		ProjectID: projectID.String(),
		Topics:    make([]*api.Topic, 0),
	}

	// Loop over topics. For each db.Topic inside the array, create a topic
	// which will be an api.Topic{} and assign the ID and Name fetched from db.Topic
	// to that struct and then append to the out.Topics array.
	for _, dbTopic := range topics {
		topic := &api.Topic{
			ID:   dbTopic.ID.String(),
			Name: dbTopic.Name,
		}
		out.Topics = append(out.Topics, topic)
	}

	c.JSON(http.StatusOK, out)
}

// ProjectTopicCreate adds a topic to a project in the database
// and returns a 201 StatusCreated response.
//
// Route: /projects/:projectID/topics
func (s *Server) ProjectTopicCreate(c *gin.Context) {
	var (
		err   error
		topic *api.Topic
		out   *api.Topic
	)

	var projectID ulid.ULID
	if projectID, err = ulid.Parse(c.Param("projectID")); err != nil {
		log.Error().Err(err).Msg("could not parse project ulid")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse project id"))
		return
	}

	// Bind the user request with JSON and return a 400 response
	// if binding is not successful.
	if err = c.BindJSON(&topic); err != nil {
		log.Warn().Err(err).Msg("could not bind project topic create request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not bind request"))
		return
	}

	// Verify that a topic ID does not exist and return a 400 response
	// if the topic ID exists.
	if topic.ID != "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("topic id cannot be specified on create"))
		return
	}

	// Verify that a topic name has been provided and return a 400 response
	// if the topic name does not exist.
	if topic.Name == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("topic name is required"))
		return
	}

	t := &db.Topic{
		ProjectID: projectID,
		Name:      topic.Name,
	}

	// Add topic to the database and return a 500 response if not successful.
	if err = db.CreateTopic(c.Request.Context(), t); err != nil {
		log.Error().Err(err).Msg("could not create project topic in the database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not add project topic"))
		return
	}

	out = &api.Topic{
		ID:   t.ID.String(),
		Name: topic.Name,
	}

	c.JSON(http.StatusCreated, out)
}

// Route: /topics
func (s *Server) TopicCreate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

// TopicList retrieves all topics assigned to an organization
// and returns a 200 OK response.
//
// Route: /topics
func (s *Server) TopicList(c *gin.Context) {
	var (
		err   error
		topic *tokens.Claims
	)

	// Fetch topic from the context.
	if topic, err = middleware.GetClaims(c); err != nil {
		log.Error().Err(err).Msg("could not fetch topic from context")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not fetch topic from context"))
		return
	}

	// Get topic's project id and return a 400 response if it is not a ULID.
	var projectID ulid.ULID
	if projectID, err = ulid.Parse(topic.ProjectID); err != nil {
		log.Error().Err(err).Msg("could not parse project id")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse project id"))
		return
	}

	// Get topics from the database and return a 500 response if not successful.
	var topics []*db.Topic
	if topics, err = db.ListTopics(c.Request.Context(), projectID); err != nil {
		log.Error().Err(err).Msg("could not fetch topics from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not fetch topics from database"))
		return
	}
	// Build the response.
	out := &api.TopicPage{Topics: make([]*api.Topic, 0)}

	// Loop over db.Topic and retrieve each topic.
	for _, dbTopic := range topics {
		topic := &api.Topic{
			ID:   dbTopic.ID.String(),
			Name: dbTopic.Name,
		}
		out.Topics = append(out.Topics, topic)
	}

	c.JSON(http.StatusOK, out)
}

// TopicDetail retrieves a summary detail of a topic with a given ID
// and returns a 200 OK response.
//
// Route: /topic/:topicID
func (s *Server) TopicDetail(c *gin.Context) {
	var (
		err   error
		reply *api.Topic
	)

	// Get the topic ID from the URL and return a 400 response
	// if the topic does not exist.
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

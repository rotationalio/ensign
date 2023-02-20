package tenant

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	pb "github.com/rotationalio/ensign/pkg/api/v1beta1"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	middleware "github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		out.Topics = append(out.Topics, dbTopic.ToAPI())
	}

	c.JSON(http.StatusOK, out)
}

// ProjectTopicCreate adds a topic to a project in the database
// and returns a 201 StatusCreated response.
//
// Route: /projects/:projectID/topics
func (s *Server) ProjectTopicCreate(c *gin.Context) {
	var (
		err    error
		claims *tokens.Claims
		topic  *api.Topic
	)

	// Fetch member claims from the context.
	if claims, err = middleware.GetClaims(c); err != nil {
		log.Error().Err(err).Msg("could not fetch member from context")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(err))
		return
	}

	// Get the member's organization ID and return a 500 response if it is not a ULID.
	var orgID ulid.ULID
	if orgID, err = ulid.Parse(claims.OrgID); err != nil {
		log.Error().Err(err).Msg("could not parse org id")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not parse org id"))
		return
	}

	// Get project ID from the URL and return a 400 response if it is missing.
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
		OrgID:     orgID,
		ProjectID: projectID,
		Name:      topic.Name,
	}

	// Add topic to the database and return a 500 response if not successful.
	if err = db.CreateTopic(c.Request.Context(), t); err != nil {
		log.Error().Err(err).Msg("could not create project topic in the database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not add project topic"))
		return
	}

	c.JSON(http.StatusCreated, t.ToAPI())
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
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not fetch topic from context"))
		return
	}

	// Get topic's organization id and return a 500 response if it is not a ULID.
	var orgID ulid.ULID
	if orgID, err = ulid.Parse(topic.OrgID); err != nil {
		log.Error().Err(err).Msg("could not parse org id")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not parse org id"))
		return
	}

	// Get topics from the database and return a 500 response if not successful.
	var topics []*db.Topic
	if topics, err = db.ListTopics(c.Request.Context(), orgID); err != nil {
		log.Error().Err(err).Msg("could not fetch topics from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not fetch topics from database"))
		return
	}

	// Build the response.
	out := &api.TopicPage{Topics: make([]*api.Topic, 0)}

	// Loop over db.Topic and retrieve each topic.
	for _, dbTopic := range topics {
		out.Topics = append(out.Topics, dbTopic.ToAPI())
	}

	c.JSON(http.StatusOK, out)
}

// TopicDetail retrieves a summary detail of a topic with a given ID
// and returns a 200 OK response.
//
// Route: /topic/:topicID
func (s *Server) TopicDetail(c *gin.Context) {
	var err error

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

	c.JSON(http.StatusOK, topic.ToAPI())
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

	c.JSON(http.StatusOK, t.ToAPI())
}

// TopicDelete completely destroys a topic, removing the metadata in Trtl and as well
// as all of the data in Ensign. Because this is irreversible, the first call returns
// a confirmation token to the user. The user must provide this token in a subsequent
// request in order to confirm the deletion. Because this operation is asynchronous,
// the endpoint returns a 202 Accepted response.
//
// Route: /topic/:topicID
func (s *Server) TopicDelete(c *gin.Context) {
	var (
		err    error
		ctx    context.Context
		claims *tokens.Claims
	)

	// User credentials are required for Quarterdeck requests
	if ctx, err = middleware.ContextFromRequest(c); err != nil {
		log.Error().Err(err).Msg("could not create user context from request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not fetch credentials for authenticated user"))
		return
	}

	// User claims are required to verify that the user owns the topic
	if claims, err = middleware.GetClaims(c); err != nil {
		log.Error().Err(err).Msg("could not fetch claims from context")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not fetch claims from context"))
		return
	}

	// Get the topic ID from the URL and return a 400 response
	// if the ID is not parseable
	var topicID ulid.ULID
	if topicID, err = ulid.Parse(c.Param("topicID")); err != nil {
		log.Warn().Err(err).Msg("could not parse topic id")
		c.JSON(http.StatusNotFound, api.ErrorResponse("topic not found"))
		return
	}

	// Parse the request body for the confirmation token
	confirm := &api.Confirmation{}
	if err = c.BindJSON(confirm); err != nil {
		log.Warn().Err(err).Msg("could not bind topic delete request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not bind user request"))
		return
	}

	// Sanity check that the ID in the request body matches the ID in the URL
	if confirm.ID != topicID.String() {
		log.Warn().Err(err).Msg("topic id in request body does not match topic id in URL")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("id in request body does not match id in URL"))
		return
	}

	// Fetch the topic metadata from the database
	var topic *db.Topic
	if topic, err = db.RetrieveTopic(ctx, topicID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			log.Warn().Err(err).Str("topicID", topicID.String()).Msg("topic not found")
			c.JSON(http.StatusNotFound, api.ErrorResponse("topic not found"))
			return
		}
		log.Error().Err(err).Str("topicID", topicID.String()).Msg("could not retrieve topic")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not delete topic"))
		return
	}

	// Verify that the user owns the topic
	if claims.OrgID != topic.OrgID.String() {
		log.Warn().Err(err).Str("user_org", claims.OrgID).Str("topic_org", topic.OrgID.String()).Msg("topic OrgID does not match user OrgID")
		c.JSON(http.StatusNotFound, api.ErrorResponse("topic not found"))
		return
	}

	// Send confirmation token if not provided
	if confirm.Token == "" {
		// Create a short-lived confirmation token in the database
		if topic.ConfirmDeleteToken, err = db.NewResourceToken(topic.ID); err != nil {
			log.Error().Err(err).Str("topicID", topicID.String()).Msg("could not generate confirmation token")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not generate confirmation token"))
			return
		}
		if err = db.UpdateTopic(ctx, topic); err != nil {
			log.Error().Err(err).Str("topicID", topicID.String()).Msg("could not save topic")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not generate confirmation token"))
			return
		}

		confirm.Name = topic.Name
		confirm.Token = topic.ConfirmDeleteToken
		c.JSON(http.StatusOK, confirm)
		return
	}

	// Check that the token is valid and has not expired
	token := &db.ResourceToken{}
	if err = token.Decode(confirm.Token); err != nil {
		log.Warn().Err(err).Msg("could not decode confirmation token")
		c.JSON(http.StatusPreconditionFailed, api.ErrorResponse("invalid confirmation token"))
		return
	}

	if token.IsExpired() {
		log.Warn().Msg("confirmation token has expired")
		c.JSON(http.StatusPreconditionFailed, api.ErrorResponse("invalid confirmation token"))
		return
	}

	// Verify that the right token was provided
	if confirm.Token != topic.ConfirmDeleteToken {
		log.Warn().Msg("confirmation tokens do not match")
		c.JSON(http.StatusPreconditionFailed, api.ErrorResponse("invalid confirmation token"))
		return
	}

	// Request access to the project from Quarterdeck
	req := &qd.Project{
		ProjectID: topic.ProjectID,
	}
	var rep *qd.LoginReply
	if rep, err = s.quarterdeck.ProjectAccess(ctx, req); err != nil {
		log.Error().Err(err).Msg("could not request one-time claims")
		c.JSON(qd.ErrorStatus(err), api.ErrorResponse("could not delete topic"))
		return
	}

	// Create the Ensign context from the one-time claims
	ensignContext := qd.ContextWithToken(ctx, rep.AccessToken)

	// Send the delete topic request to Ensign
	deleteRequest := &pb.TopicMod{
		Id:        topic.ID.String(),
		Operation: pb.TopicMod_DESTROY,
	}
	var tombstone *pb.TopicTombstone
	if tombstone, err = s.ensign.DeleteTopic(ensignContext, deleteRequest); err != nil {
		log.Error().Err(err).Msg("could not delete topic in ensign")
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, api.ErrorResponse("topic not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not delete topic"))
		return
	}

	// The delete request is asynchronous so just update the state in the database
	topic.State = tombstone.State
	if err = db.UpdateTopic(ctx, topic); err != nil {
		log.Error().Err(err).Str("topicID", topicID.String()).Msg("could not update topic state")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not delete topic"))
		return
	}

	// Set 202 for the response so the frontend knows the delete is in progress
	confirm.Name = topic.Name
	confirm.Status = tombstone.State.String()
	c.JSON(http.StatusAccepted, confirm)
}

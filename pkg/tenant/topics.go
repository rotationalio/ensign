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
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
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
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not list topics"))
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
		err   error
		ctx   context.Context
		topic *api.Topic
	)

	// Get user credentials to make request to Quarterdeck.
	if ctx, err = middleware.ContextFromRequest(c); err != nil {
		log.Error().Err(err).Msg("could not create user context from request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not fetch user credentials"))
		return
	}

	// orgID is required to check ownership of the project.
	var orgID ulid.ULID
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Bind the user request with JSON and return a 400 response if binding is not successful.
	if err = c.BindJSON(&topic); err != nil {
		log.Warn().Err(err).Msg("could not bind project topic create request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not bind request"))
		return
	}

	// Verify that a topic ID does not exist and return a 400 response if the topic ID exists.
	if topic.ID != "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("topic id cannot be specified on create"))
		return
	}

	// Verify that a topic name has been provided and return a 400 response if the topic name does not exist.
	if topic.Name == "" {
		c.JSON(http.StatusBadRequest, api.ErrorResponse("topic name is required"))
		return
	}

	// Get project ID from the URL and return a 404 response if it is missing.
	var projectID ulid.ULID
	if projectID, err = ulid.Parse(c.Param("projectID")); err != nil {
		log.Error().Err(err).Str("projectID", projectID.String()).Msg("could not parse project ulid")
		c.JSON(http.StatusNotFound, api.ErrorResponse("project not found"))
		return
	}

	// Retrieve project from the database.
	var project *db.Project
	if project, err = db.RetrieveProject(ctx, projectID); err != nil {
		log.Error().Err(err).Str("projectID", projectID.String()).Msg("could not retrieve project from database")
		c.JSON(http.StatusNotFound, api.ErrorResponse("project not found"))
		return
	}

	// Ensure project belongs to the organization of the requesting user.
	if orgID.Compare(project.OrgID) != 0 {
		log.Error().Err(err).Str("user_orgID", orgID.String()).Str("project_orgID", project.OrgID.String()).Msg("project org ID does not match user org ID")
		c.JSON(http.StatusNotFound, api.ErrorResponse("project not found"))
		return
	}

	// Get access to the project from Quarterdeck.
	req := &qd.Project{
		ProjectID: project.ID,
	}

	var rep *qd.LoginReply
	if rep, err = s.quarterdeck.ProjectAccess(ctx, req); err != nil {
		log.Error().Err(err).Msg("could not get access to project claims")
		c.JSON(qd.ErrorStatus(err), api.ErrorResponse("could not create topic"))
		return
	}

	// Create Ensign context.
	// TODO: ensure the context has PerRPCCredentials for gRPC authentication
	enCtx := qd.ContextWithToken(ctx, rep.AccessToken)

	// Send create project topic request to Ensign.
	create := &pb.Topic{
		ProjectId: project.ID[:],
	}

	var enTopic *pb.Topic
	if enTopic, err = s.ensign.CreateTopic(enCtx, create); err != nil {
		log.Error().Err(err).Msg("could not create topic in ensign")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create topic"))
		return
	}

	// Add topic to the database and return a 500 response if not successful.
	t := &db.Topic{
		OrgID:     project.OrgID,
		ProjectID: projectID,
		Name:      enTopic.Name,
	}

	if err = db.CreateTopic(ctx, t); err != nil {
		log.Error().Err(err).Msg("could not create project topic")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create project topic"))
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
		orgID ulid.ULID
	)

	// orgID is required to retrieve the topic
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Get topics from the database and return a 500 response if not successful.
	var topics []*db.Topic
	if topics, err = db.ListTopics(c.Request.Context(), orgID); err != nil {
		log.Error().Err(err).Msg("could not fetch topics from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not list topics"))
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
		c.JSON(http.StatusNotFound, api.ErrorResponse("topic not found"))
		return
	}

	c.JSON(http.StatusOK, topic.ToAPI())
}

// TopicUpdate updates the record of a topic with a given ID and returns a 200 OK
// response. The editable fields are the topic name and state, although the topic state
// can only be set to READONLY which archives the topic.
//
// Route: /topic/:topicID
func (s *Server) TopicUpdate(c *gin.Context) {
	var (
		err   error
		ctx   context.Context
		topic *api.Topic
	)

	// User credentials are required for Quarterdeck requests
	if ctx, err = middleware.ContextFromRequest(c); err != nil {
		log.Error().Err(err).Msg("could not create user context from request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not fetch credentials for authenticated user"))
		return
	}

	// orgID is required to check ownership of the topic
	var orgID ulid.ULID
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

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

	// Sanity check that the ID in the URL matches the ID in the request body.
	if topic.ID != topicID.String() {
		log.Warn().Msg("topic id in request body does not match topic id in URL")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("id in request body does not match id in URL"))
		return
	}

	// Fetch the topic metadata from the database.
	var t *db.Topic
	if t, err = db.RetrieveTopic(ctx, topicID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			log.Warn().Err(err).Str("topicID", topicID.String()).Msg("topic not found")
			c.JSON(http.StatusNotFound, api.ErrorResponse("topic not found"))
			return
		}
		log.Error().Err(err).Str("topicID", topicID.String()).Msg("could not retrieve topic")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update topic"))
		return
	}

	// Ensure the new name is valid
	t.Name = topic.Name
	if err = t.Validate(); err != nil {
		log.Warn().Err(err).Msg("could not validate topic update")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Verify that the user owns the topic
	if orgID.Compare(t.OrgID) != 0 {
		log.Warn().Str("user_org", orgID.String()).Str("topic_org", t.OrgID.String()).Msg("user does not own topic")
		c.JSON(http.StatusNotFound, api.ErrorResponse("topic not found"))
		return
	}

	// Check if we have to update the topic state
	if topic.State != t.State.String() {
		// Topic state can only be set to READONLY
		if topic.State != pb.TopicTombstone_READONLY.String() {
			c.JSON(http.StatusBadRequest, api.ErrorResponse("topic state can only be set to READONLY"))
			return
		}

		// Don't proceed if the topic is already being deleted
		if t.State == pb.TopicTombstone_DELETING {
			c.JSON(http.StatusBadRequest, api.ErrorResponse("topic is already being deleted"))
			return
		}

		// Request one-time claims for the topic update request
		req := &qd.Project{
			ProjectID: t.ProjectID,
		}
		var rep *qd.LoginReply
		if rep, err = s.quarterdeck.ProjectAccess(ctx, req); err != nil {
			log.Error().Err(err).Msg("could not request one-time claims")
			c.JSON(qd.ErrorStatus(err), api.ErrorResponse("could not update topic"))
			return
		}

		// Create the Ensign context with the one-time claims
		ensignContext := qd.ContextWithToken(ctx, rep.AccessToken)

		// Archive means a "soft" delete in Ensign (no data is destroyed)
		updateRequest := &pb.TopicMod{
			Id:        t.ID.String(),
			Operation: pb.TopicMod_ARCHIVE,
		}
		var tombstone *pb.TopicTombstone
		if tombstone, err = s.ensign.DeleteTopic(ensignContext, updateRequest); err != nil {
			if status.Code(err) == codes.NotFound {
				log.Warn().Err(err).Str("topicID", updateRequest.Id).Msg("topic not found in ensign even though it is in tenant")
				c.JSON(http.StatusNotFound, api.ErrorResponse("topic not found"))
				return
			}
			log.Error().Err(err).Msg("could not update topic in ensign")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update topic"))
			return
		}
		t.State = tombstone.State
	}

	// Update topic in the database and return a 500 response if the topic
	// record cannot be updated.
	if err = db.UpdateTopic(ctx, t); err != nil {
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
		err error
		ctx context.Context
	)

	// User credentials are required for Quarterdeck requests
	if ctx, err = middleware.ContextFromRequest(c); err != nil {
		log.Error().Err(err).Msg("could not create user context from request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not fetch credentials for authenticated user"))
		return
	}

	// orgID is required to verify that the user owns the topic
	var orgID ulid.ULID
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
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
		log.Warn().Msg("topic id in request body does not match topic id in URL")
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
	if orgID.Compare(topic.OrgID) != 0 {
		log.Warn().Str("user_org", orgID.String()).Str("topic_org", topic.OrgID.String()).Msg("topic OrgID does not match user OrgID")
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
		if status.Code(err) == codes.NotFound {
			log.Error().Err(err).Str("topicID", deleteRequest.Id).Msg("topic not found in ensign even though it is in tenant")
			c.JSON(http.StatusNotFound, api.ErrorResponse("topic not found"))
			return
		}
		log.Error().Err(err).Msg("could not delete topic in ensign")
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

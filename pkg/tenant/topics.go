package tenant

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	middleware "github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	tk "github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	pg "github.com/rotationalio/ensign/pkg/utils/pagination"
	responses "github.com/rotationalio/ensign/pkg/utils/responses"
	"github.com/rotationalio/ensign/pkg/utils/sentry"
	"github.com/rotationalio/ensign/pkg/utils/tasks"
	"github.com/rotationalio/ensign/pkg/utils/tokens"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/rotationalio/ensign/pkg/utils/units"
	pb "github.com/rotationalio/go-ensign/api/v1beta1"
	"github.com/rs/zerolog/log"
)

// ProjectTopicList retrieves topics assigned to a specified
// project and returns a 200 OK response.
//
// Route: /projects/:projectID/topics
func (s *Server) ProjectTopicList(c *gin.Context) {
	var (
		err        error
		next, prev *pg.Cursor
	)

	query := &api.PageQuery{}
	if err = c.BindQuery(query); err != nil {
		log.Error().Err(err).Msg("could not parse query")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse query"))
		return
	}

	if query.NextPageToken != "" {
		if prev, err = pg.Parse(query.NextPageToken); err != nil {
			c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse next page token"))
			return
		}
	} else {
		prev = pg.New("", "", int32(query.PageSize))
	}

	// orgID is required to check project ownership.
	var orgID ulid.ULID
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Get the topic's project ID from the URL and return a 400 response
	// if the project ID is not a ULID.
	var projectID ulid.ULID
	if projectID, err = ulid.Parse(c.Param("projectID")); err != nil {
		sentry.Warn(c).Err(err).Str("projectID", c.Param("projectID")).Msg("could not parse project id")
		c.JSON(http.StatusNotFound, api.ErrorResponse("project not found"))
		return
	}

	// Verify project exists in the organization.
	if err = db.VerifyOrg(c, orgID, projectID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("project not found"))
			return
		}
		sentry.Warn(c).Err(err).Msg("could not check verification")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not verify organization"))
		return
	}

	// Get topics from the database and return a 500 response
	// if not successful.
	var topics []*db.Topic
	if topics, next, err = db.ListTopics(c.Request.Context(), projectID, prev); err != nil {
		sentry.Error(c).Err(err).Msg("could not list topics in database")
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

	if next != nil {
		if out.NextPageToken, err = next.NextPageToken(); err != nil {
			log.Error().Err(err).Msg("could not set next page token")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not list topics"))
			return
		}
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
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	// orgID is required to create the topic.
	var orgID ulid.ULID
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Get the user ID from the context
	var userID ulid.ULID
	if userID = userIDFromContext(c); ulids.IsZero(userID) {
		return
	}

	// Bind the user request with JSON and return a 400 response if binding is not successful.
	if err = c.BindJSON(&topic); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse topic create request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
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

	// Get project ID from the URL.
	var projectID ulid.ULID
	if projectID, err = ulid.Parse(c.Param("projectID")); err != nil {
		sentry.Warn(c).Err(err).Str("projectID", c.Param("projectID")).Msg("could not parse project id")
		c.JSON(http.StatusNotFound, api.ErrorResponse("project not found"))
		return
	}

	// Verify project exists in the organization.
	if err = db.VerifyOrg(c, orgID, projectID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("project not found"))
			return
		}
		sentry.Warn(c).Err(err).Msg("could not check verification")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not verify organization"))
		return
	}

	// Get the access token for the Ensign request. This method handles logging and
	// error responses.
	var accessToken string
	if accessToken, err = s.EnsignProjectToken(ctx, userID, projectID); err != nil {
		sentry.Warn(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	// Create the topic in Ensign.
	var topicID string
	if topicID, err = s.ensign.InvokeOnce(accessToken).CreateTopic(ctx, topic.Name); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing ensign error in tenant")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create topic"))
		return
	}

	// Add topic to the database and return a 500 response if not successful.
	t := &db.Topic{
		OrgID:     orgID,
		ProjectID: projectID,
		Name:      topic.Name,
	}

	if t.ID, err = ulid.Parse(topicID); err != nil {
		sentry.Error(c).Err(err).Str("topicID", topicID).Msg("could not parse topic id created by ensign")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create topic"))
		return
	}

	if err = db.CreateTopic(ctx, t); err != nil {
		sentry.Error(c).Err(err).Msg("could not create topic in database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create topic"))
		return
	}

	// Update project stats in the background
	s.tasks.QueueContext(middleware.TaskContext(c), tasks.TaskFunc(func(ctx context.Context) error {
		return s.UpdateProjectStats(ctx, userID, t.ProjectID)
	}), tasks.WithError(fmt.Errorf("could not update stats for project %s", t.ProjectID.String())))

	c.JSON(http.StatusCreated, t.ToAPI())
}

// Route: /topics
func (s *Server) TopicCreate(c *gin.Context) {
	sentry.Warn(c).Msg("topic create not implemented yet")
	c.JSON(http.StatusNotImplemented, "not implemented yet")
}

// TopicList retrieves topics assigned to a specified organization
// and returns a 200 OK response.
//
// Route: /topics
func (s *Server) TopicList(c *gin.Context) {
	var (
		err        error
		orgID      ulid.ULID
		next, prev *pg.Cursor
	)

	// orgID is required to retrieve the topic
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	query := &api.PageQuery{}
	if err = c.BindQuery(query); err != nil {
		log.Error().Err(err).Msg("could not parse query")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse query"))
		return
	}

	if query.NextPageToken != "" {
		if prev, err = pg.Parse(query.NextPageToken); err != nil {
			c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse next page token"))
			return
		}
	} else {
		prev = pg.New("", "", int32(query.PageSize))
	}

	// Get topics from the database.
	var topics []*db.Topic
	if topics, next, err = db.ListTopics(c.Request.Context(), orgID, prev); err != nil {
		sentry.Error(c).Err(err).Msg("could not list topics in database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not list topics"))
		return
	}

	// Build the response.
	out := &api.TopicPage{Topics: make([]*api.Topic, 0)}

	// Loop over db.Topic and retrieve each topic.
	for _, dbTopic := range topics {
		out.Topics = append(out.Topics, dbTopic.ToAPI())
	}

	if next != nil {
		if out.NextPageToken, err = next.NextPageToken(); err != nil {
			log.Error().Err(err).Msg("could not set next page token")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not list topics"))
			return
		}
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
		orgID ulid.ULID
	)

	// orgID is required to check ownership of the topic
	// TODO: Ensure that the topic exists in the organization.
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Get the topic ID from the URL and return a 400 response
	// if the topic does not exist.
	var topicID ulid.ULID
	if topicID, err = ulid.Parse(c.Param("topicID")); err != nil {
		sentry.Warn(c).Err(err).Str("topicID", c.Param("topicID")).Msg("could not parse topic id")
		c.JSON(http.StatusNotFound, api.ErrorResponse("topic not found"))
		return
	}

	// Verify topic exists in the organization.
	if err = db.VerifyOrg(c, orgID, topicID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("topic not found"))
			return
		}
		sentry.Warn(c).Err(err).Msg("could not check verification")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not verify organization"))
		return
	}

	// Get the specified topic from the database and return a 404 response
	// if it cannot be retrieved.
	var topic *db.Topic
	if topic, err = db.RetrieveTopic(c.Request.Context(), topicID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("topic not found"))
			return
		}

		sentry.Error(c).Err(err).Msg("could not retrieve topic from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve topic"))
		return
	}

	c.JSON(http.StatusOK, topic.ToAPI())
}

// TopicEvents returns an event info "breakdown" for the topic, which includes info
// about all of the schema Types in the topic and the number of events and storage each
// schema Type contributes to.
//
// Route: /topic/:topicID/events
func (s *Server) TopicEvents(c *gin.Context) {
	var (
		err   error
		orgID ulid.ULID
		ctx   context.Context
	)

	// Get user credentials to make request to Quarterdeck.
	if ctx, err = middleware.ContextFromRequest(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	// orgID is required to check ownership of the topic.
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Parse the topicID from the URL.
	var topicID ulid.ULID
	if topicID, err = ulid.Parse(c.Param("topicID")); err != nil {
		sentry.Warn(c).Err(err).Str("topicID", c.Param("topicID")).Msg("could not parse topic id")
		c.JSON(http.StatusNotFound, api.ErrorResponse(responses.ErrTopicNotFound))
		return
	}

	// Verify topic exists in the organization.
	if err = db.VerifyOrg(c, orgID, topicID); err != nil {
		if errors.Is(err, db.ErrNotFound) || errors.Is(err, db.ErrOrgNotVerified) {
			c.JSON(http.StatusNotFound, api.ErrorResponse(responses.ErrTopicNotFound))
			return
		}

		sentry.Warn(c).Err(err).Msg("could not check verification")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Retrieve the topic from the database to get the project ID.
	var topic *db.Topic
	if topic, err = db.RetrieveTopic(c.Request.Context(), topicID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse(responses.ErrTopicNotFound))
			return
		}

		sentry.Error(c).Err(err).Msg("could not retrieve topic from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Get the access token for the Ensign request. This method handles logging and
	// error responses.
	var accessToken string
	if accessToken, err = s.EnsignProjectToken(ctx, orgID, topic.ProjectID); err != nil {
		sentry.Warn(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	// Get info for this specific topic from Ensign.
	var info *pb.ProjectInfo
	if info, err = s.ensign.InvokeOnce(accessToken).Info(c, topic.ID.String()); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing ensign error in tenant")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse(responses.ErrSomethingWentWrong))
		return
	}

	// Build the response.
	out := make([]*api.EventTypeInfo, 0)
	for _, topic := range info.Topics {
		// Sanity check: this is the topic we are looking for
		if !bytes.Equal(topic.TopicId, topicID[:]) {
			continue
		}

		for _, typeInfo := range topic.Types {
			info := &api.EventTypeInfo{
				Type:     typeInfo.Type.Name,
				Version:  typeInfo.Type.Semver(),
				Mimetype: typeInfo.Mimetype.MimeType(),
				Events: &api.StatValue{
					Name:  "events",
					Value: float64(typeInfo.Events),
				},
				Duplicates: &api.StatValue{
					Name:  "duplicates",
					Value: float64(typeInfo.Duplicates),
				},
				Storage: &api.StatValue{
					Name: "storage",
				},
			}
			if topic.Events > 0 {
				info.Events.Percent = (float64(typeInfo.Events) / float64(topic.Events)) * 100
			}
			if topic.Duplicates > 0 {
				info.Duplicates.Percent = (float64(typeInfo.Duplicates) / float64(topic.Duplicates)) * 100
			}
			info.Storage.Units, info.Storage.Value = units.FromBytes(typeInfo.DataSizeBytes)
			if topic.DataSizeBytes > 0 {
				info.Storage.Percent = (float64(typeInfo.DataSizeBytes) / float64(topic.DataSizeBytes)) * 100
			}
			out = append(out, info)
		}

		// Ensure only one topic is counted
		break
	}

	if len(out) == 0 {
		sentry.Warn(c).ULID("topicID", topicID).Msg("topic exists in tenant but was not returned by ensign")
		c.JSON(http.StatusNotFound, api.ErrorResponse(responses.ErrTopicNotFound))
		return
	}

	c.JSON(http.StatusOK, out)
}

// TopicStats returns a snapshot of statistics for a topic with a given ID in a 200 OK
// response.
//
// Route: /topic/:topicID/stats
func (s *Server) TopicStats(c *gin.Context) {
	var (
		err   error
		orgID ulid.ULID
	)

	// orgID is required to check ownership of the topic.
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Parse the topicID from the URL.
	var topicID ulid.ULID
	if topicID, err = ulid.Parse(c.Param("topicID")); err != nil {
		sentry.Warn(c).Err(err).Str("topicID", c.Param("topicID")).Msg("could not parse topic id")
		c.JSON(http.StatusNotFound, api.ErrorResponse("topic not found"))
		return
	}

	// Verify topic exists in the organization.
	if err = db.VerifyOrg(c, orgID, topicID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("topic not found"))
			return
		}
		sentry.Warn(c).Err(err).Msg("could not check verification")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not verify organization"))
		return
	}

	// Retrieve the topic from the database.
	var topic *db.Topic
	if topic, err = db.RetrieveTopic(c.Request.Context(), topicID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("topic not found"))
			return
		}

		sentry.Error(c).Err(err).Msg("could not retrieve topic from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not retrieve topic stats"))
		return
	}

	// Construct the stats reply.
	// TODO: Data storage percentage and units are currently hardcoded.
	out := []*api.StatValue{
		{
			Name:  "Online Publishers",
			Value: float64(topic.Publishers.Active),
		},
		{
			Name:  "Online Subscribers",
			Value: float64(topic.Subscribers.Active),
		},
		{
			Name:  "Total Events",
			Value: float64(topic.Events),
		},
		{
			Name:    "Data Storage",
			Value:   float64(topic.Storage),
			Units:   "GB",
			Percent: 0.0,
		},
	}

	c.JSON(http.StatusOK, out)
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
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	// orgID is required to check ownership of the topic
	var orgID ulid.ULID
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Get the user ID from the context
	var userID ulid.ULID
	if userID = userIDFromContext(c); ulids.IsZero(userID) {
		return
	}

	// Get the topic ID from the URL and return a 400 response if
	// the topic ID is not a ULID.
	var topicID ulid.ULID
	if topicID, err = ulid.Parse(c.Param("topicID")); err != nil {
		sentry.Warn(c).Err(err).Str("topicID", c.Param("topicID")).Msg("could not parse topic id")
		c.JSON(http.StatusNotFound, api.ErrorResponse("topic not found"))
		return
	}

	// Verify topic exists in the organization.
	if err = db.VerifyOrg(c, orgID, topicID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("topic not found"))
			return
		}
		sentry.Warn(c).Err(err).Msg("could not check verification")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not verify organization"))
		return
	}

	// Bind the user request with JSON and return a 400 response if
	// binding is not successful.
	if err = c.BindJSON(&topic); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse topic update request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
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
			c.JSON(http.StatusNotFound, api.ErrorResponse("topic not found"))
			return
		}

		sentry.Error(c).Err(err).Msg("could not retrieve topic from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update topic"))
		return
	}

	// Ensure the new name is valid
	t.Name = topic.Name
	if err = t.Validate(); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse(err))
		return
	}

	// Check if we have to update the topic state
	// TODO: Do we need a dedicated endpoint for this?
	if topic.Status != t.Status() {
		// Topic state can only be set to readonly/archived
		if topic.Status != db.TopicStatusArchived {
			c.JSON(http.StatusBadRequest, api.ErrorResponse("topic state can only be set to Archived"))
			return
		}

		// Don't proceed if the topic is already being deleted
		if t.State == pb.TopicTombstone_DELETING {
			c.JSON(http.StatusBadRequest, api.ErrorResponse("topic is already being deleted"))
			return
		}

		// Get the access token for the Ensign request. This method handles logging and
		// error responses.
		var accessToken string
		if accessToken, err = s.EnsignProjectToken(ctx, userID, t.ProjectID); err != nil {
			sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
			api.ReplyQuarterdeckError(c, err)
			return
		}

		// Archive means a "soft" delete in Ensign (no data is destroyed)
		if err = s.ensign.InvokeOnce(accessToken).ArchiveTopic(ctx, t.ID.String()); err != nil {
			switch {
			case err.Error() == "not implemented yet":
				sentry.Warn(c).Err(err).Msg("this version of the Go SDK does not support topic archiving")
				c.JSON(http.StatusNotImplemented, api.ErrorResponse("archiving a topic is not supported"))
				return
			default:
				sentry.Debug(c).Err(err).Msg("tracing ensign error in tenant")
				c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not update topic"))
				return
			}
		}
		t.State = pb.TopicTombstone_READONLY
	}

	// Update topic in the database and return a 500 response if the topic
	// record cannot be updated.
	if err = db.UpdateTopic(ctx, t); err != nil {
		sentry.Error(c).Err(err).Msg("could not update topic in database after ensign update")
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
		sentry.Error(c).Err(err).Msg("could not get user claims from authenticated request")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse(api.ErrInvalidUserClaims))
		return
	}

	// orgID is required to verify that the user owns the topic
	var orgID ulid.ULID
	if orgID = orgIDFromContext(c); ulids.IsZero(orgID) {
		return
	}

	// Get the user ID from the context
	var userID ulid.ULID
	if userID = userIDFromContext(c); ulids.IsZero(userID) {
		return
	}

	// Get the topic ID from the URL and return a 400 response
	// if the ID is not parseable
	var topicID ulid.ULID
	if topicID, err = ulid.Parse(c.Param("topicID")); err != nil {
		sentry.Warn(c).Err(err).Str("topicID", c.Param("topicID")).Msg("could not parse topic id")
		c.JSON(http.StatusNotFound, api.ErrorResponse("topic not found"))
		return
	}

	// Verify topic exists in the organization.
	if err = db.VerifyOrg(c, orgID, topicID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("topic not found"))
			return
		}
		sentry.Warn(c).Err(err).Msg("could not check verification")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("could not verify organization"))
		return
	}

	// Parse the request body for the confirmation token
	confirm := &api.Confirmation{}
	if err = c.BindJSON(confirm); err != nil {
		sentry.Warn(c).Err(err).Msg("could not parse topic delete confirmation request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse(api.ErrUnparsable))
		return
	}

	// Sanity check that the ID in the request body matches the ID in the URL
	if confirm.ID != topicID.String() {
		c.Error(err)
		c.JSON(http.StatusBadRequest, api.ErrorResponse("id in request body does not match id in URL"))
		return
	}

	// Fetch the topic metadata from the database
	var topic *db.Topic
	if topic, err = db.RetrieveTopic(ctx, topicID); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse("topic not found"))
			return
		}

		sentry.Error(c).Err(err).Msg("could not retrieve topic from database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not delete topic"))
		return
	}

	// Send confirmation token if not provided
	if confirm.Token == "" {
		// Create a short-lived confirmation token in the database
		if topic.ConfirmDeleteToken, err = tokens.NewConfirmation(topic.ID); err != nil {
			sentry.Error(c).Err(err).Msg("could not generate confirmation token")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not generate confirmation token"))
			return
		}

		if err = db.UpdateTopic(ctx, topic); err != nil {
			sentry.Error(c).Err(err).Msg("could not update topic in database")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not generate confirmation token"))
			return
		}

		confirm.Name = topic.Name
		confirm.Token = topic.ConfirmDeleteToken
		c.JSON(http.StatusOK, confirm)
		return
	}

	// Check that the token is valid and has not expired
	token := &tokens.Confirmation{}
	if err = token.Decode(confirm.Token); err != nil {
		sentry.Warn(c).Err(err).Msg("could not decode topic delete confirmation token")
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

	// Get the access token for the Ensign request. This method handles logging and
	// error responses.
	var accessToken string
	if accessToken, err = s.EnsignProjectToken(ctx, userID, topic.ProjectID); err != nil {
		sentry.Debug(c).Err(err).Msg("tracing quarterdeck error in tenant")
		api.ReplyQuarterdeckError(c, err)
		return
	}

	// Request topic delete from Ensign, which will destroy the topic and all of its data
	if err = s.ensign.InvokeOnce(accessToken).DestroyTopic(ctx, topic.ID.String()); err != nil {
		// TODO: Update with the standard errors defined by the SDK
		switch {
		case err.Error() == "not implemented yet":
			sentry.Warn(c).Err(err).Msg("this version of the Go SDK does not support topic deletion")
			c.JSON(http.StatusNotImplemented, api.ErrorResponse("deleting a topic is not supported"))
			return
		default:
			sentry.Debug(c).Err(err).Msg("tracing ensign error in tenant")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not delete topic"))
			return
		}
	}

	// The delete request is asynchronous so just update the state in the database
	topic.State = pb.TopicTombstone_DELETING
	if err = db.UpdateTopic(ctx, topic); err != nil {
		sentry.Error(c).Err(err).Msg("could not update tombstone topic in database")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not delete topic"))
		return
	}

	// Set 202 for the response so the frontend knows the delete is in progress
	confirm.Name = topic.Name
	confirm.Status = topic.State.String()
	c.JSON(http.StatusAccepted, confirm)
}

// EnsignProjectToken is a helper method to request access to an Ensign project on
// behalf of a user. This type of access is different from API keys; it requires
// obtaining a short-lived access token from the Quarterdeck service by providing user
// credentials. It also only carries permissions for managing topics (e.g. no pub/sub)
// based on the permissions the user had when the token was issued. This method makes
// an external request to Quarterdeck but uses a cache to avoid repeated requests. This
// method only returns an error if the request to Quarterdeck fails.
func (s *Server) EnsignProjectToken(ctx context.Context, userID, projectID ulid.ULID) (_ string, err error) {
	// Get the access token from the cache if it exists
	if token, err := s.tokens.Get(userID, projectID); err == nil {
		return token, nil
	} else {
		log.Debug().Err(err).Msg("could not get access token from cache")
	}

	// Request a new access token from Quarterdeck
	req := &qd.Project{
		ProjectID: projectID,
	}
	var rep *qd.LoginReply
	if rep, err = s.quarterdeck.ProjectAccess(ctx, req); err != nil {
		return "", err
	}

	// Add the access token to the cache
	if err = s.tokens.Add(userID, projectID, rep.AccessToken); err != nil {
		log.Error().Err(err).Msg("could not add access token to cache")
	}

	return rep.AccessToken, nil
}

// Helper to fetch the userID from the gin context. This method also logs and returns
// any errors to allow endpoints to have consistent error handling.
func userIDFromContext(c *gin.Context) (userID ulid.ULID) {
	var (
		claims *tk.Claims
		err    error
	)
	if claims, err = middleware.GetClaims(c); err != nil {
		sentry.Error(c).Err(err).Msg("could not get user claims from context")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("user claims unavailable"))
		return ulid.ULID{}
	}

	if userID = claims.ParseOrgID(); ulids.IsZero(userID) {
		sentry.Error(c).Err(err).Msg("could not parse userID from claims")
		c.JSON(http.StatusUnauthorized, api.ErrorResponse("invalid user claims"))
		return ulid.ULID{}
	}

	return userID
}

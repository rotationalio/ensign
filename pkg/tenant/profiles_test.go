package tenant_test

import (
	"bytes"
	"context"
	"net/http"
	"time"

	qd "github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/middleware"
	"github.com/rotationalio/ensign/pkg/quarterdeck/mock"
	"github.com/rotationalio/ensign/pkg/quarterdeck/permissions"
	"github.com/rotationalio/ensign/pkg/quarterdeck/tokens"
	"github.com/rotationalio/ensign/pkg/tenant/api/v1"
	"github.com/rotationalio/ensign/pkg/tenant/db"
	"github.com/rotationalio/ensign/pkg/utils/responses"
	"github.com/rotationalio/ensign/pkg/utils/ulids"
	"github.com/trisacrypto/directory/pkg/trtl/pb/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *tenantTestSuite) TestProfileDetail() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Setup the trtl mock
	trtl := db.GetMock()
	defer trtl.Reset()

	// Create the member fixture
	member := &db.Member{
		OrgID:        ulids.New(),
		ID:           ulids.New(),
		Name:         "Cleon I",
		Email:        "cleon@empire.org",
		Organization: "Empire",
		Workspace:    "empire",
		Role:         permissions.RoleMember,
	}

	data, err := member.MarshalValue()
	require.NoError(err, "could not marshal the member fixture")

	key, err := member.Key()
	require.NoError(err, "could not create the member record key")

	// Trtl Get should return the member data
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		switch in.Namespace {
		case db.MembersNamespace:
			if !bytes.Equal(in.Key, key) {
				return nil, status.Errorf(codes.NotFound, "member not found")
			}
			return &pb.GetReply{Value: data}, nil
		default:
			return nil, status.Errorf(codes.Internal, "unknown namespace: %s", in.Namespace)
		}
	}

	validClaims := &tokens.Claims{
		Name:  member.Name,
		Email: member.Email,
		OrgID: member.OrgID.String(),
	}
	validClaims.Subject = member.ID.String()

	s.Run("Happy Path", func() {
		require.NoError(s.SetClientCredentials(validClaims))

		expected := &api.Member{
			ID:                member.ID.String(),
			Name:              member.Name,
			Email:             member.Email,
			Picture:           member.Picture(),
			Organization:      member.Organization,
			Workspace:         member.Workspace,
			ProfessionSegment: db.ProfessionSegmentUnspecified.String(),
			Role:              member.Role,
			OnboardingStatus:  db.MemberStatusOnboarding.String(),
		}

		// Make the request
		rep, err := s.client.ProfileDetail(ctx)
		require.NoError(err, "could not make the profile detail request")
		require.Equal(expected, rep, "response does not match expected")
	})

	s.Run("Missing user ID", func() {
		claims := &tokens.Claims{
			Name:  member.Name,
			Email: member.Email,
			OrgID: member.OrgID.String(),
		}
		require.NoError(s.SetClientCredentials(claims))

		// Should error if no user ID is present
		_, err := s.client.ProfileDetail(ctx)
		s.requireHTTPError(err, http.StatusInternalServerError)
	})

	s.Run("Missing org ID", func() {
		claims := &tokens.Claims{
			Name:  member.Name,
			Email: member.Email,
		}
		claims.Subject = member.ID.String()
		require.NoError(s.SetClientCredentials(claims))

		// Should error if no org ID is present
		_, err := s.client.ProfileDetail(ctx)
		s.requireHTTPError(err, http.StatusInternalServerError)
	})

	s.Run("Member not found", func() {
		claims := &tokens.Claims{
			Name:  member.Name,
			Email: member.Email,
			OrgID: member.OrgID.String(),
		}
		claims.Subject = ulids.New().String()

		// Should error if member is not found by ID
		require.NoError(s.SetClientCredentials(claims))
		_, err := s.client.ProfileDetail(ctx)
		s.requireHTTPError(err, http.StatusInternalServerError)
	})
}

func (s *tenantTestSuite) TestProfileUpdate() {
	require := s.Require()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Setup valid claims for tests
	orgID := ulids.New()
	memberID := ulids.New()
	validClaims := &tokens.Claims{
		Name:  "Hari Seldon",
		Email: "seldon@foundation",
		OrgID: orgID.String(),
	}
	validClaims.Subject = memberID.String()

	// Setup the trtl mock
	trtl := db.GetMock()
	defer trtl.Reset()

	// Trtl Get returns a byte encoded member fixture
	var data, key []byte
	trtl.OnGet = func(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
		switch in.Namespace {
		case db.MembersNamespace:
			if !bytes.Equal(in.Key, key) {
				return nil, status.Errorf(codes.NotFound, "member not found")
			}
			return &pb.GetReply{Value: data}, nil
		default:
			return nil, status.Errorf(codes.NotFound, "unknown namespace: %s", in.Namespace)
		}
	}

	// Trtl Put should return success for the correct namespace and key
	trtl.OnPut = func(ctx context.Context, in *pb.PutRequest) (*pb.PutReply, error) {
		switch in.Namespace {
		case db.MembersNamespace:
			if !bytes.Equal(in.Key, key) {
				return nil, status.Errorf(codes.NotFound, "member not found")
			}
			return &pb.PutReply{}, nil
		default:
			return nil, status.Errorf(codes.NotFound, "unknown namespace: %s", in.Namespace)
		}
	}

	s.Run("Invited User", func() {
		require.NoError(s.SetClientCSRFProtection())
		require.NoError(s.SetClientCredentials(validClaims))

		// Existing member fixture returned by the mock
		member := &db.Member{
			OrgID:        orgID,
			ID:           memberID,
			Name:         "Hari Seldon",
			Email:        "seldon@foundation",
			Organization: "Foundation",
			Workspace:    "foundation",
			Role:         permissions.RoleMember,
			Invited:      true,
			JoinedAt:     time.Now(),
		}

		var err error
		data, err = member.MarshalValue()
		require.NoError(err, "could not marshal the member fixture")

		key, err = member.Key()
		require.NoError(err, "could not create the member record key")

		// Create a request with updated name, org, and workspace
		req := &api.Member{
			ID:           member.ID.String(),
			Name:         "Raven Seldon",
			Email:        member.Email,
			Organization: "Second Foundation",
			Workspace:    "second-foundation",
			Role:         permissions.RoleMember,
		}

		// Invited user should not have their organization or workspace updated
		expected := &api.Member{
			ID:                member.ID.String(),
			Name:              req.Name,
			Email:             member.Email,
			Organization:      member.Organization,
			Workspace:         member.Workspace,
			ProfessionSegment: db.ProfessionSegmentUnspecified.String(),
			Role:              member.Role,
			Invited:           true,
			Picture:           member.Picture(),
			OnboardingStatus:  db.MemberStatusOnboarding.String(),
		}

		// Make the request
		rep, err := s.client.ProfileUpdate(ctx, req)
		require.NoError(err, "could not make the profile update request")
		rep.Created, rep.DateAdded = "", ""
		require.Equal(expected, rep, "response does not match expected")
	})

	s.Run("Invited User Done Onboarding", func() {
		require.NoError(s.SetClientCSRFProtection())
		require.NoError(s.SetClientCredentials(validClaims))

		// Existing member fixture returned by the mock
		member := &db.Member{
			OrgID:        orgID,
			ID:           memberID,
			Name:         "Hari Seldon",
			Email:        "seldon@foundation",
			Organization: "Foundation",
			Workspace:    "foundation",
			Role:         permissions.RoleMember,
			Invited:      true,
			JoinedAt:     time.Now(),
		}

		var err error
		data, err = member.MarshalValue()
		require.NoError(err, "could not marshal the member fixture")

		key, err = member.Key()
		require.NoError(err, "could not create the member record key")

		// Create a request with all fields completed
		req := &api.Member{
			ID:                member.ID.String(),
			Name:              "Raven Seldon",
			Email:             member.Email,
			ProfessionSegment: "Personal",
			DeveloperSegment:  []string{"Data science"},
			Role:              permissions.RoleMember,
		}

		// Invited user should be done onboarding
		expected := &api.Member{
			ID:                member.ID.String(),
			Name:              req.Name,
			Email:             member.Email,
			Organization:      member.Organization,
			Workspace:         member.Workspace,
			ProfessionSegment: req.ProfessionSegment,
			DeveloperSegment:  req.DeveloperSegment,
			Role:              member.Role,
			Invited:           true,
			Picture:           member.Picture(),
			OnboardingStatus:  db.MemberStatusActive.String(),
		}

		// Make the request
		rep, err := s.client.ProfileUpdate(ctx, req)
		require.NoError(err, "could not make the profile update request")
		rep.Created, rep.DateAdded = "", ""
		require.Equal(expected, rep, "response does not match expected")
	})

	s.Run("Organization Owner Clear Workspace", func() {
		require.NoError(s.SetClientCSRFProtection())
		require.NoError(s.SetClientCredentials(validClaims))

		// Existing member fixture returned by the mock
		member := &db.Member{
			OrgID:        orgID,
			ID:           memberID,
			Name:         "Hari Seldon",
			Email:        "seldon@foundation",
			Organization: "Foundation",
			Workspace:    "foundation",
			Role:         permissions.RoleOwner,
			JoinedAt:     time.Now(),
		}

		var err error
		data, err = member.MarshalValue()
		require.NoError(err, "could not marshal the member fixture")

		key, err = member.Key()
		require.NoError(err, "could not create the member record key")

		// Create a request with updated name and org, workspace is empty
		req := &api.Member{
			ID:           member.ID.String(),
			Name:         "Raven Seldon",
			Email:        member.Email,
			Organization: "Second Foundation",
			Role:         permissions.RoleMember,
		}

		// Organization owner should have name and organization updated, workspace is
		// cleared
		expected := &api.Member{
			ID:                member.ID.String(),
			Name:              req.Name,
			Email:             member.Email,
			Organization:      req.Organization,
			ProfessionSegment: db.ProfessionSegmentUnspecified.String(),
			Role:              permissions.RoleOwner,
			Picture:           member.Picture(),
			OnboardingStatus:  db.MemberStatusOnboarding.String(),
		}

		// Make the request
		rep, err := s.client.ProfileUpdate(ctx, req)
		require.NoError(err, "could not make the profile update request")
		rep.Created, rep.DateAdded = "", ""
		require.Equal(expected, rep, "response does not match expected")
	})

	s.Run("Organization Owner Partial Progress", func() {
		require.NoError(s.SetClientCSRFProtection())
		require.NoError(s.SetClientCredentials(validClaims))

		// Existing member fixture returned by the mock
		member := &db.Member{
			OrgID:        orgID,
			ID:           memberID,
			Name:         "Hari Seldon",
			Email:        "seldon@foundation",
			Organization: "Foundation",
			Workspace:    "foundation",
			Role:         permissions.RoleOwner,
			JoinedAt:     time.Now(),
		}

		var err error
		data, err = member.MarshalValue()
		require.NoError(err, "could not marshal the member fixture")

		key, err = member.Key()
		require.NoError(err, "could not create the member record key")

		// Create a request with updated name, org, and workspace
		req := &api.Member{
			ID:           member.ID.String(),
			Name:         "Raven Seldon",
			Email:        member.Email,
			Organization: "Second Foundation",
			Workspace:    "foundation",
			Role:         permissions.RoleMember,
		}

		// Quarterdeck returns that the workspace is still available
		workspace := &qd.Workspace{
			OrgID:  orgID,
			Name:   member.Organization,
			Domain: member.Workspace,
		}
		s.quarterdeck.OnWorkspace(mock.UseStatus(http.StatusOK), mock.UseJSONFixture(workspace), mock.RequireAuth())

		// Organization owner should have name and organization updated
		expected := &api.Member{
			ID:                member.ID.String(),
			Name:              req.Name,
			Email:             member.Email,
			Organization:      req.Organization,
			Workspace:         req.Workspace,
			ProfessionSegment: db.ProfessionSegmentUnspecified.String(),
			Role:              permissions.RoleOwner,
			Picture:           member.Picture(),
			OnboardingStatus:  db.MemberStatusOnboarding.String(),
		}

		// Make the request
		rep, err := s.client.ProfileUpdate(ctx, req)
		require.NoError(err, "could not make the profile update request")
		rep.Created, rep.DateAdded = "", ""
		require.Equal(expected, rep, "response does not match expected")
	})

	s.Run("Organization Owner Change Workspace", func() {
		require.NoError(s.SetClientCSRFProtection())
		require.NoError(s.SetClientCredentials(validClaims))

		// Existing member fixture returned by the mock
		member := &db.Member{
			OrgID:        orgID,
			ID:           memberID,
			Name:         "Hari Seldon",
			Email:        "seldon@foundation",
			Organization: "Foundation",
			Workspace:    "foundation",
			Role:         permissions.RoleOwner,
			JoinedAt:     time.Now(),
		}

		var err error
		data, err = member.MarshalValue()
		require.NoError(err, "could not marshal the member fixture")

		key, err = member.Key()
		require.NoError(err, "could not create the member record key")

		// Create a request with updated name, org, and workspace
		req := &api.Member{
			ID:           member.ID.String(),
			Name:         "Raven Seldon",
			Email:        member.Email,
			Organization: "Second Foundation",
			Workspace:    "second-foundation",
			Role:         permissions.RoleMember,
		}

		// Quarterdeck returns that the new workspace is available
		workspace := &qd.Workspace{
			IsAvailable: true,
		}
		s.quarterdeck.OnWorkspace(mock.UseStatus(http.StatusOK), mock.UseJSONFixture(workspace), mock.RequireAuth())

		// Organization owner should have name, organization, and workspace updated
		expected := &api.Member{
			ID:                member.ID.String(),
			Name:              req.Name,
			Email:             member.Email,
			Organization:      req.Organization,
			Workspace:         req.Workspace,
			ProfessionSegment: db.ProfessionSegmentUnspecified.String(),
			Role:              permissions.RoleOwner,
			Picture:           member.Picture(),
			OnboardingStatus:  db.MemberStatusOnboarding.String(),
		}

		// Make the request
		rep, err := s.client.ProfileUpdate(ctx, req)
		require.NoError(err, "could not make the profile update request")
		rep.Created, rep.DateAdded = "", ""
		require.Equal(expected, rep, "response does not match expected")
	})

	s.Run("Organization Owner Done Onboarding", func() {
		require.NoError(s.SetClientCSRFProtection())
		require.NoError(s.SetClientCredentials(validClaims))

		// Existing member fixture returned by the mock
		member := &db.Member{
			OrgID:        orgID,
			ID:           memberID,
			Name:         "Hari Seldon",
			Email:        "seldon@foundation",
			Organization: "Foundation",
			Workspace:    "foundation",
			Role:         permissions.RoleOwner,
			JoinedAt:     time.Now(),
		}

		var err error
		data, err = member.MarshalValue()
		require.NoError(err, "could not marshal the member fixture")

		key, err = member.Key()
		require.NoError(err, "could not create the member record key")

		// Create a request with updated name, org, and workspace
		req := &api.Member{
			ID:                member.ID.String(),
			Name:              "Raven Seldon",
			Email:             member.Email,
			Organization:      "Second Foundation",
			Workspace:         "second-foundation",
			ProfessionSegment: "Personal",
			DeveloperSegment:  []string{"Data science"},
			Role:              permissions.RoleMember,
		}

		// Quarterdeck returns that the new workspace is available
		workspace := &qd.Workspace{
			IsAvailable: true,
		}
		s.quarterdeck.OnWorkspace(mock.UseStatus(http.StatusOK), mock.UseJSONFixture(workspace), mock.RequireAuth())

		// Organization owner should be done onboarding
		expected := &api.Member{
			ID:                member.ID.String(),
			Name:              req.Name,
			Email:             member.Email,
			Organization:      req.Organization,
			Workspace:         req.Workspace,
			ProfessionSegment: req.ProfessionSegment,
			DeveloperSegment:  req.DeveloperSegment,
			Role:              permissions.RoleOwner,
			Picture:           member.Picture(),
			OnboardingStatus:  db.MemberStatusActive.String(),
		}

		// Quarterdeck mock should return success for the organization update
		qdReply := &qd.Organization{
			ID:     orgID,
			Name:   req.Organization,
			Domain: req.Workspace,
		}
		s.quarterdeck.OnOrganizationsUpdate(orgID.String(), mock.UseStatus(http.StatusOK), mock.UseJSONFixture(qdReply), mock.RequireAuth())

		// Make the request
		rep, err := s.client.ProfileUpdate(ctx, req)
		require.NoError(err, "could not make the profile update request")
		rep.Created, rep.DateAdded = "", ""
		require.Equal(expected, rep, "response does not match expected")
	})

	s.Run("Organization Owner Workspace Taken", func() {
		require.NoError(s.SetClientCSRFProtection())
		require.NoError(s.SetClientCredentials(validClaims))

		// Existing member fixture returned by the mock
		member := &db.Member{
			OrgID:        orgID,
			ID:           memberID,
			Name:         "Hari Seldon",
			Email:        "seldon@foundation",
			Organization: "Foundation",
			Workspace:    "foundation",
			Role:         permissions.RoleOwner,
			JoinedAt:     time.Now(),
		}

		var err error
		data, err = member.MarshalValue()
		require.NoError(err, "could not marshal the member fixture")

		key, err = member.Key()
		require.NoError(err, "could not create the member record key")

		// Organization owner is still onboarding
		req := &api.Member{
			ID:           member.ID.String(),
			Name:         "Raven Seldon",
			Email:        member.Email,
			Organization: "Second Foundation",
			Workspace:    "second-foundation",
			Role:         permissions.RoleMember,
		}

		// Quarterdeck returns that the new workspace is available
		workspace := &qd.Workspace{}
		s.quarterdeck.OnWorkspace(mock.UseStatus(http.StatusOK), mock.UseJSONFixture(workspace), mock.RequireAuth())

		// Should return a 409 if the workspace is taken
		_, err = s.client.ProfileUpdate(ctx, req)
		s.requireHTTPError(err, http.StatusConflict)
	})

	s.Run("No CSRF Token", func() {
		require.NoError(s.ClearClientCSRFProtection())
		require.NoError(s.SetClientCredentials(validClaims))

		// Should error if no CSRF token is present
		_, err := s.client.ProfileUpdate(ctx, &api.Member{})
		s.requireError(err, http.StatusForbidden, middleware.ErrCSRFVerification.Error(), "expected CSRF error")
	})

	s.Run("Missing user ID", func() {
		require.NoError(s.SetClientCSRFProtection())
		claims := &tokens.Claims{
			Name:  "Hari Seldon",
			Email: "seldon@foundation",
			OrgID: orgID.String(),
		}
		require.NoError(s.SetClientCredentials(claims))

		// Should error if no user ID is present
		_, err := s.client.ProfileUpdate(ctx, &api.Member{})
		s.requireHTTPError(err, http.StatusInternalServerError)
	})

	s.Run("Missing org ID", func() {
		require.NoError(s.SetClientCSRFProtection())
		claims := &tokens.Claims{
			Name:  "Hari Seldon",
			Email: "seldon@foundation",
		}
		claims.Subject = memberID.String()
		require.NoError(s.SetClientCredentials(claims))

		// Should error if no org ID is present
		_, err := s.client.ProfileUpdate(ctx, &api.Member{})
		s.requireHTTPError(err, http.StatusInternalServerError)
	})

	s.Run("Member not found", func() {
		claims := &tokens.Claims{
			Name:  "Hari Seldon",
			Email: "seldon@foundation",
			OrgID: orgID.String(),
		}
		claims.Subject = ulids.New().String()

		// Should error if member is not found by ID
		require.NoError(s.SetClientCSRFProtection())
		require.NoError(s.SetClientCredentials(claims))
		_, err := s.client.ProfileUpdate(ctx, &api.Member{ID: claims.Subject})
		s.requireHTTPError(err, http.StatusInternalServerError)
	})

	s.Run("Invalid Profession", func() {
		require.NoError(s.SetClientCSRFProtection())
		require.NoError(s.SetClientCredentials(validClaims))

		// Existing member fixture returned by the mock
		member := &db.Member{
			OrgID:        orgID,
			ID:           memberID,
			Name:         "Hari Seldon",
			Email:        "seldon@foundation",
			Organization: "Foundation",
			Workspace:    "foundation",
			Role:         permissions.RoleOwner,
			JoinedAt:     time.Now(),
		}

		var err error
		data, err = member.MarshalValue()
		require.NoError(err, "could not marshal the member fixture")

		key, err = member.Key()
		require.NoError(err, "could not create the member record key")

		// Should error if the profession in the request is invalid
		req := &api.Member{
			ID:                member.ID.String(),
			Workspace:         "not a valid workspace",
			ProfessionSegment: "not a valid profession",
			Role:              permissions.RoleOwner,
		}
		_, err = s.client.ProfileUpdate(ctx, req)
		s.requireError(err, http.StatusBadRequest, db.ErrProfessionUnknown.Error(), "expected unknown profession error")
	})

	s.Run("Invalid Developer", func() {
		require.NoError(s.SetClientCSRFProtection())
		require.NoError(s.SetClientCredentials(validClaims))

		// Existing member fixture returned by the mock
		member := &db.Member{
			OrgID:        orgID,
			ID:           memberID,
			Name:         "Hari Seldon",
			Email:        "seldon@foundation",
			Organization: "Foundation",
			Workspace:    "foundation",
			Role:         permissions.RoleOwner,
			JoinedAt:     time.Now(),
		}

		var err error
		data, err = member.MarshalValue()
		require.NoError(err, "could not marshal the member fixture")

		key, err = member.Key()
		require.NoError(err, "could not create the member record key")

		// Should error if the profession in the request is invalid
		req := &api.Member{
			ID:               member.ID.String(),
			Workspace:        "not a valid workspace",
			DeveloperSegment: []string{"not a valid developer"},
			Role:             permissions.RoleOwner,
		}
		_, err = s.client.ProfileUpdate(ctx, req)
		s.requireError(err, http.StatusBadRequest, db.ErrDeveloperUnknown.Error(), "expected unknown developer error")
	})

	s.Run("Invalid fields", func() {
		require.NoError(s.SetClientCSRFProtection())
		require.NoError(s.SetClientCredentials(validClaims))

		// Existing member fixture returned by the mock
		member := &db.Member{
			OrgID:        orgID,
			ID:           memberID,
			Name:         "Hari Seldon",
			Email:        "seldon@foundation",
			Organization: "Foundation",
			Workspace:    "foundation",
			Role:         permissions.RoleOwner,
			JoinedAt:     time.Now(),
		}

		var err error
		data, err = member.MarshalValue()
		require.NoError(err, "could not marshal the member fixture")

		key, err = member.Key()
		require.NoError(err, "could not create the member record key")

		// Should error if there are validation errors in the request
		req := &api.Member{
			ID:               member.ID.String(),
			Workspace:        "not a valid workspace",
			DeveloperSegment: []string{"Data science", ""},
			Role:             permissions.RoleOwner,
		}
		expected := &api.FieldValidationErrors{
			{
				Field: "workspace",
				Err:   db.ErrInvalidWorkspace.Error(),
				Index: -1,
			},
			{
				Field: "developer_segment",
				Err:   db.ErrDeveloperUnspecified.Error(),
				Index: 1,
			},
		}
		_, err = s.client.ProfileUpdate(ctx, req)
		s.requireError(err, http.StatusBadRequest, expected.Error(), "wrong validation errors returned")
	})

	s.Run("Quarterdeck Workspace Lookup Error", func() {
		require.NoError(s.SetClientCSRFProtection())
		require.NoError(s.SetClientCredentials(validClaims))

		// Existing member fixture returned by the mock
		member := &db.Member{
			OrgID:    orgID,
			ID:       memberID,
			Name:     "Hari Seldon",
			Email:    "seldon@foundation",
			Role:     permissions.RoleOwner,
			JoinedAt: time.Now(),
		}

		var err error
		data, err = member.MarshalValue()
		require.NoError(err, "could not marshal the member fixture")

		key, err = member.Key()
		require.NoError(err, "could not create the member record key")

		// Should error if workspace lookup returns an error
		s.quarterdeck.OnWorkspace(mock.UseError(http.StatusInternalServerError, responses.ErrSomethingWentWrong), mock.RequireAuth())

		req := &api.Member{
			ID:                member.ID.String(),
			Name:              "Raven Seldon",
			Email:             member.Email,
			Organization:      "Second Foundation",
			Workspace:         "second-foundation",
			ProfessionSegment: "Personal",
			DeveloperSegment:  []string{"Data science"},
			Role:              permissions.RoleMember,
		}
		_, err = s.client.ProfileUpdate(ctx, req)
		s.requireHTTPError(err, http.StatusInternalServerError)
	})

	s.Run("Quarterderk Organization Update Error", func() {
		require.NoError(s.SetClientCSRFProtection())
		require.NoError(s.SetClientCredentials(validClaims))

		// Existing member fixture returned by the mock
		member := &db.Member{
			OrgID:    orgID,
			ID:       memberID,
			Name:     "Hari Seldon",
			Email:    "seldon@foundation",
			Role:     permissions.RoleOwner,
			JoinedAt: time.Now(),
		}

		var err error
		data, err = member.MarshalValue()
		require.NoError(err, "could not marshal the member fixture")

		key, err = member.Key()
		require.NoError(err, "could not create the member record key")

		// Should error if organization update returns an error
		s.quarterdeck.OnWorkspace(mock.UseStatus(http.StatusOK), mock.UseJSONFixture(&qd.Workspace{IsAvailable: true}), mock.RequireAuth())
		s.quarterdeck.OnOrganizationsUpdate(orgID.String(), mock.UseError(http.StatusNotFound, responses.ErrOrganizationNotFound), mock.RequireAuth())

		req := &api.Member{
			ID:                member.ID.String(),
			Name:              "Raven Seldon",
			Email:             member.Email,
			Organization:      "Second Foundation",
			Workspace:         "second-foundation",
			ProfessionSegment: "Personal",
			DeveloperSegment:  []string{"Data science"},
			Role:              permissions.RoleMember,
		}
		_, err = s.client.ProfileUpdate(ctx, req)
		s.requireHTTPError(err, http.StatusNotFound)
	})
}

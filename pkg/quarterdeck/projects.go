package quarterdeck

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/quarterdeck/api/v1"
	"github.com/rotationalio/ensign/pkg/quarterdeck/db"
	"github.com/rs/zerolog/log"
)

// TODO: document
// TODO: actually implement this resource endpoint
// TODO: implement pagination
// HACK: this is just a quick hack to get us going, it should filter projects based on
// the authenticated user and organization instead of just returning everyting.
func (s *Server) ProjectList(c *gin.Context) {
	var (
		err  error
		rows *sql.Rows
		out  *api.ProjectList
	)

	// Fetch the projects from the database
	var tx *sql.Tx
	if tx, err = db.BeginTx(c.Request.Context(), &sql.TxOptions{ReadOnly: true}); err != nil {
		log.Error().Err(err).Msg("could not start database transaction")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not fetch projects"))
		return
	}
	defer tx.Rollback()

	out = &api.ProjectList{Projects: make([]*api.Project, 0)}
	if rows, err = tx.Query(`SELECT p.id, p.slug, p.name, p.description, o.domain, u.email, p.created, p.modified FROM projects p JOIN organizations o ON o.id=p.organization_id JOIN users u ON u.id=p.created_by;`); err != nil {
		log.Error().Err(err).Msg("could not list projects")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not fetch projects"))
		return
	}
	defer rows.Close()

	for rows.Next() {
		p := &api.Project{}
		if err = rows.Scan(&p.ID, &p.Slug, &p.Name, &p.Description, &p.Organization, &p.Owner, &p.Created, &p.Modified); err != nil {
			log.Error().Err(err).Msg("could not scan project")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not fetch projects"))
			return
		}
		out.Projects = append(out.Projects, p)
	}

	tx.Commit()
	c.JSON(http.StatusOK, out)
}

// TODO: document
// TODO: actually implement this resource endpoint
// HACK: this is just a quick hack to get us going: it creates a project and organization
func (s *Server) ProjectCreate(c *gin.Context) {
	var (
		err     error
		project *api.Project
	)

	if err = c.BindJSON(&project); err != nil {
		log.Warn().Err(err).Msg("could not parse create project request")
		c.JSON(http.StatusBadRequest, api.ErrorResponse("could not parse request"))
		return
	}

	// TODO: add project validation

	// Insert project into the database
	var tx *sql.Tx
	if tx, err = db.BeginTx(c.Request.Context(), &sql.TxOptions{ReadOnly: false}); err != nil {
		log.Error().Err(err).Msg("could not start database transaction")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create project"))
		return
	}
	defer tx.Rollback()

	// Get or create the organization ID for the project
	var orgID int64
	if err = tx.QueryRow(`SELECT id FROM organizations WHERE domain=$1`, project.Organization).Scan(&orgID); err != nil {
		if err == sql.ErrNoRows {
			var result sql.Result
			if result, err = tx.Exec(`INSERT INTO organizations (name, domain, created, modified) VALUES ($1, $2, datetime('now'), datetime('now'));`, "Unknown", project.Organization); err != nil {
				log.Error().Err(err).Msg("could not create organization")
				c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create project"))
				return
			}

			// HACK: should not ignore the error here, just trying to move quickly.
			orgID, _ = result.LastInsertId()
		} else {
			log.Error().Err(err).Msg("could not fetch organization ID")
			c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create project"))
			return
		}
	}

	var (
		projectID int64
		result    sql.Result
	)
	if result, err = tx.Exec(`INSERT INTO projects (slug, name, description, organization_id, created_by, created, modified) VALUES ($1, $2, $3, $4, (SELECT id FROM users WHERE email=$5), datetime('now'), datetime('now'))`, project.Slug, project.Name, project.Description, orgID, project.Owner); err != nil {
		log.Error().Err(err).Msg("could not create project")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create project"))
		return
	}

	// HACK: should not ignore error here
	projectID, _ = result.LastInsertId()

	// Populate return response
	if err = tx.QueryRow(`SELECT p.id, p.slug, p.name, p.description, o.domain, u.email, p.created, p.modified FROM projects p JOIN organizations o ON o.id=p.organization_id JOIN users u ON u.id=p.created_by;`, projectID).Scan(&project.ID, &project.Slug, &project.Name, &project.Description, &project.Organization, &project.Owner, &project.Created, &project.Modified); err != nil {
		log.Error().Err(err).Msg("could not create project")
		c.JSON(http.StatusInternalServerError, api.ErrorResponse("could not create project"))
		return
	}

	tx.Commit()
	c.JSON(http.StatusCreated, project)
}

func (s *Server) ProjectDetail(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, api.ErrorResponse("not yet implemented"))
}

func (s *Server) ProjectUpdate(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, api.ErrorResponse("not yet implemented"))
}

func (s *Server) ProjectDelete(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, api.ErrorResponse("not yet implemented"))
}

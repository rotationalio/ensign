package uptime

import (
	"embed"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rotationalio/ensign/pkg/uptime/db"
	"github.com/rotationalio/ensign/pkg/uptime/health"
	"github.com/rotationalio/ensign/pkg/uptime/incident"
	"github.com/rotationalio/ensign/pkg/uptime/services"
	"github.com/rs/zerolog/log"
)

// content holds our static web server content.
//
//go:embed all:templates
//go:embed all:static
var content embed.FS

type CSSColor string
type CSSIcon string

const (
	CSSPrimary   CSSColor = "primary"
	CSSSecondary CSSColor = "secondary"
	CSSSuccess   CSSColor = "success"
	CSSDanger    CSSColor = "danger"
	CSSWarning   CSSColor = "warning"
	CSSInfo      CSSColor = "info"
	CSSLight     CSSColor = "light"
	CSSDark      CSSColor = "dark"
)

const (
	CSSUnknown     CSSIcon = "circle-question"
	CSSOnline      CSSIcon = "circle-check"
	CSSDegraded    CSSIcon = "battery-quarter"
	CSSPartial     CSSIcon = "triangle-exclamation"
	CSSOutage      CSSIcon = "skull-crossbones"
	CSSMaintenance CSSIcon = "wrench"
)

func IconFromStatus(s health.Status) CSSIcon {
	switch s {
	case health.Online:
		return CSSOnline
	case health.Maintenance:
		return CSSMaintenance
	case health.Stopping, health.Degraded:
		return CSSDegraded
	case health.Unhealthy:
		return CSSPartial
	case health.Offline, health.Outage:
		return CSSOutage
	default:
		return CSSUnknown
	}
}

func ColorFromStatus(s health.Status) CSSColor {
	switch s {
	case health.Online:
		return CSSSuccess
	case health.Maintenance:
		return CSSInfo
	case health.Stopping, health.Degraded:
		return CSSSecondary
	case health.Unhealthy:
		return CSSWarning
	case health.Offline, health.Outage:
		return CSSDanger
	default:
		return CSSSecondary
	}
}

// StatusPageContext is used to render the content into the index page for the uptime
// HTML status page (info page). Note that the template uses Bootstrap CSS classes for
// status colors and other properties, which is enforced using the CSSColor type.
type StatusPageContext struct {
	StatusMessage   string
	StatusColor     CSSColor
	ServiceGroups   []ServiceGroupContext
	IncidentHistory []IncidentDayContext
}

type ServiceGroupContext struct {
	Title         string
	ServiceStates []ServiceStateContext
}

type ServiceStateContext struct {
	Title       string
	StatusColor CSSColor
	StatusIcon  CSSIcon
}

type IncidentDayContext struct {
	Date      time.Time
	Incidents []IncidentContext
}

type IncidentContext struct {
	Description string
	StartTime   time.Time
	EndTime     time.Time
	StatusColor CSSColor
	StatusIcon  CSSIcon
}

func (i IncidentContext) TimeFormat() string {
	if i.EndTime.IsZero() || i.EndTime.Equal(i.StartTime) {
		return i.StartTime.Format("Jan 02, 2006 at 15:04 MST")
	}

	// Make sure the timezones are the same for start and end time
	endTime := i.EndTime.In(i.StartTime.Location())
	if DateEqual(i.StartTime, i.EndTime) {
		return fmt.Sprintf("%s - %s %s", i.StartTime.Format("Jan 02, 2006 from 15:04"), endTime.Format("15:04"), i.StartTime.Format("MST"))
	}

	if i.StartTime.Year() == i.EndTime.Year() {
		return fmt.Sprintf("from %s - %s %s", i.StartTime.Format("Jan 02, 15:04"), endTime.Format("Jan 02, 15:04"), i.StartTime.Format("2006 MST"))
	}

	return fmt.Sprintf("%s - %s", i.StartTime.Format("Jan 02 2006, 15:04 MST"), i.EndTime.Format("Jan 02 2006, 15:04 MST"))
}

func DateEqual(date1, date2 time.Time) bool {
	date2 = date2.In(date1.Location())
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func (s *Server) Index(c *gin.Context) {
	// Create a context to render the web page with
	status := &StatusPageContext{
		StatusMessage: "Unknown Rotational Systems Status",
		StatusColor:   CSSSecondary,
	}

	// Load the service states from the db
	serviceInfo := &services.Info{}
	if err := db.Get(db.KeyCurrentStatus, serviceInfo); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			log.Warn().Err(err).Msg("no monitors have been checked yet")
		} else {
			log.Error().Err(err).Msg("could not retrieve service info from database")
		}
	}

	// Create the services contexts
	var worstStatus health.Status
	status.ServiceGroups = make([]ServiceGroupContext, 0, len(serviceInfo.Groups))
	for _, group := range serviceInfo.Groups {
		sgroup := ServiceGroupContext{
			Title:         group.Title,
			ServiceStates: make([]ServiceStateContext, 0, len(group.Services)),
		}

		for _, service := range group.Services {
			// Global Description of Rotational Status
			if service.Status > worstStatus {
				worstStatus = service.Status
				status.StatusColor = ColorFromStatus(service.Status)
				switch service.Status {
				case health.Online:
					status.StatusMessage = "All Rotational Systems Operational"
				case health.Maintenance:
					status.StatusMessage = "Ongoing Maintenance: Some Services may be Temporarily Unavailable"
				case health.Stopping, health.Degraded:
					status.StatusMessage = "Some Rotational Systems are Experiencing Degraded Performance"
				case health.Unhealthy:
					status.StatusMessage = "Partial Outages Detected: Rotational Systems are Unhealthy"
				case health.Offline, health.Outage:
					status.StatusMessage = "Major Outages Detected: Rotational Systems are Unavailable"
				default:
					status.StatusMessage = "Unknown Rotational Systems Status"
				}
			}

			// Create the Service Context
			sstate := ServiceStateContext{
				Title:       service.Title,
				StatusColor: ColorFromStatus(service.Status),
				StatusIcon:  IconFromStatus(service.Status),
			}
			sgroup.ServiceStates = append(sgroup.ServiceStates, sstate)
		}
		status.ServiceGroups = append(status.ServiceGroups, sgroup)
	}

	// Fetch Incidents from the database
	days, err := incident.LastWeek()
	if err != nil {
		log.Error().Err(err).Msg("could not fetch incidents from db")
	}

	status.IncidentHistory = make([]IncidentDayContext, 0, len(days))
	for _, day := range days {
		idc := IncidentDayContext{
			Date:      day.Date,
			Incidents: make([]IncidentContext, 0, len(day.Incidents)),
		}

		for _, incident := range day.Incidents {
			idc.Incidents = append(idc.Incidents, IncidentContext{
				Description: incident.Description,
				StartTime:   incident.StartTime,
				EndTime:     incident.EndTime,
				StatusColor: ColorFromStatus(incident.Status),
				StatusIcon:  IconFromStatus(incident.Status),
			})
		}

		status.IncidentHistory = append(status.IncidentHistory, idc)
	}

	c.HTML(http.StatusOK, "index.html", status)
}

func (s *Server) Services(c *gin.Context) {
	// Create a context to render the web page with
	status := &StatusPageContext{
		StatusMessage: "Unknown Rotational Systems Status",
		StatusColor:   CSSSecondary,
	}

	// Load the service states from the db
	serviceInfo := &services.Info{}
	if err := db.Get(db.KeyCurrentStatus, serviceInfo); err != nil {
		if errors.Is(err, db.ErrNotFound) {
			log.Warn().Err(err).Msg("no monitors have been checked yet")
		} else {
			log.Error().Err(err).Msg("could not retrieve service info from database")
		}
	}

	// Create the services contexts
	var worstStatus health.Status
	status.ServiceGroups = make([]ServiceGroupContext, 0, len(serviceInfo.Groups))
	for _, group := range serviceInfo.Groups {
		sgroup := ServiceGroupContext{
			Title:         group.Title,
			ServiceStates: make([]ServiceStateContext, 0, len(group.Services)),
		}

		for _, service := range group.Services {
			// Global Description of Rotational Status
			if service.Status > worstStatus {
				worstStatus = service.Status
				status.StatusColor = ColorFromStatus(service.Status)
				switch service.Status {
				case health.Online:
					status.StatusMessage = "All Rotational Systems Operational"
				case health.Maintenance:
					status.StatusMessage = "Ongoing Maintenance: Some Services may be Temporarily Unavailable"
				case health.Stopping, health.Degraded:
					status.StatusMessage = "Some Rotational Systems are Experiencing Degraded Performance"
				case health.Unhealthy:
					status.StatusMessage = "Partial Outages Detected: Rotational Systems are Unhealthy"
				case health.Offline, health.Outage:
					status.StatusMessage = "Major Outages Detected: Rotational Systems are Unavailable"
				default:
					status.StatusMessage = "Unknown Rotational Systems Status"
				}
			}

			// Create the Service Context
			sstate := ServiceStateContext{
				Title:       service.Title,
				StatusColor: ColorFromStatus(service.Status),
				StatusIcon:  IconFromStatus(service.Status),
			}
			sgroup.ServiceStates = append(sgroup.ServiceStates, sstate)
		}
		status.ServiceGroups = append(status.ServiceGroups, sgroup)
	}

	c.HTML(http.StatusOK, "services.html", status)
}

func (s *Server) Incidents(c *gin.Context) {
	// Create a context to render the web page with
	status := &StatusPageContext{
		StatusMessage: "Unknown Rotational Systems Status",
		StatusColor:   CSSSecondary,
	}

	// Fetch Incidents from the database
	days, err := incident.LastWeek()
	if err != nil {
		log.Error().Err(err).Msg("could not fetch incidents from db")
	}

	status.IncidentHistory = make([]IncidentDayContext, 0, len(days))
	for _, day := range days {
		idc := IncidentDayContext{
			Date:      day.Date,
			Incidents: make([]IncidentContext, 0, len(day.Incidents)),
		}

		for _, incident := range day.Incidents {
			idc.Incidents = append(idc.Incidents, IncidentContext{
				Description: incident.Description,
				StartTime:   incident.StartTime,
				EndTime:     incident.EndTime,
				StatusColor: ColorFromStatus(incident.Status),
				StatusIcon:  IconFromStatus(incident.Status),
			})
		}

		status.IncidentHistory = append(status.IncidentHistory, idc)
	}

	c.HTML(http.StatusOK, "incidents.html", status)
}

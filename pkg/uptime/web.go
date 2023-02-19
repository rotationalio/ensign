package uptime

import (
	"embed"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
	CSSOnline      CSSIcon = "circle-check"
	CSSDegraded    CSSIcon = "battery-quarter"
	CSSPartial     CSSIcon = "triangle-exclamation"
	CSSOutage      CSSIcon = "skull-crossbones"
	CSSMaintenance CSSIcon = "wrench"
)

// StatusPageContext is used to render the content into the index page for the uptime
// HTML status page (info page). Note that the template uses Bootstrap CSS classes for
// status colors and other properties, which is enforced using the CSSColor type.
type StatusPageContext struct {
	StatusMessage   string
	StatusColor     CSSColor
	ServiceStates   []ServiceStateContext
	IncidentHistory []IncidentDayContext
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
	// TODO: load the context from the database
	status := &StatusPageContext{
		StatusMessage: "All Rotational Systems Operational",
		StatusColor:   CSSSuccess,
		ServiceStates: []ServiceStateContext{
			{
				"Ensign Eventing API",
				CSSSuccess, CSSOnline,
			},
			{
				"Ensign Placement API",
				CSSSecondary, CSSDegraded,
			},
			{
				"Tenant (Beacon BFF) API",
				CSSDanger, CSSOutage,
			},
			{
				"Quarterdeck Authentication API",
				CSSWarning, CSSPartial,
			},
			{
				"Rotational Frontend API",
				CSSSuccess, CSSOnline,
			},
			{
				"Trtl Replicated Database",
				CSSInfo, CSSMaintenance,
			},
		},
		IncidentHistory: []IncidentDayContext{
			{
				Date:      time.Date(2023, 2, 17, 0, 0, 0, 0, time.UTC),
				Incidents: []IncidentContext{},
			},
			{
				Date: time.Date(2023, 2, 16, 0, 0, 0, 0, time.UTC),
				Incidents: []IncidentContext{
					{
						Description: "Detected version change from v1.4.3 to v1.5.0",
						StartTime:   time.Date(2023, 2, 16, 16, 41, 31, 313432, time.UTC),
						StatusColor: CSSInfo,
						StatusIcon:  CSSMaintenance,
					},
					{
						Description: "Major service outage",
						StartTime:   time.Date(2022, 12, 23, 8, 32, 02, 495123, time.UTC),
						EndTime:     time.Date(2023, 2, 16, 5, 21, 49, 0, time.Local),
						StatusColor: CSSDanger,
						StatusIcon:  CSSOutage,
					},
				},
			},
			{
				Date:      time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC),
				Incidents: []IncidentContext{},
			},
			{
				Date: time.Date(2023, 2, 14, 0, 0, 0, 0, time.UTC),
				Incidents: []IncidentContext{
					{
						Description: "Quarterdeck has slowed down",
						StartTime:   time.Date(2023, 2, 14, 9, 32, 02, 495123, time.Local),
						EndTime:     time.Date(2023, 2, 14, 18, 21, 49, 0, time.Local),
						StatusColor: CSSSecondary,
						StatusIcon:  CSSDegraded,
					},
					{
						Description: "Partial Ensign service outage",
						StartTime:   time.Date(2023, 2, 13, 8, 32, 02, 495123, time.UTC),
						EndTime:     time.Date(2023, 2, 14, 5, 21, 49, 0, time.Local),
						StatusColor: CSSWarning,
						StatusIcon:  CSSPartial,
					},
				},
			},
			{
				Date:      time.Date(2023, 2, 13, 0, 0, 0, 0, time.UTC),
				Incidents: []IncidentContext{},
			},
			{
				Date:      time.Date(2023, 2, 12, 0, 0, 0, 0, time.UTC),
				Incidents: []IncidentContext{},
			},
			{
				Date:      time.Date(2023, 2, 11, 0, 0, 0, 0, time.UTC),
				Incidents: []IncidentContext{},
			},
		},
	}

	c.HTML(http.StatusOK, "index.html", status)
}

func (s *Server) NotFound(c *gin.Context) {
	c.String(http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

func (s *Server) NotAllowed(c *gin.Context) {
	c.String(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
}

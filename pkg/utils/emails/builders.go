package emails

import (
	"bytes"
	"embed"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/rotationalio/ensign/pkg/utils/sendgrid"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Email templates must be provided in this directory and are loaded at compile time
// using go:embed.
const templatesDir = "templates"

var (
	//go:embed templates/*.html templates/*.txt
	files     embed.FS
	templates map[string]*template.Template
)

// Load templates when the package is imported
func init() {
	var (
		err           error
		templateFiles []fs.DirEntry
	)

	templates = make(map[string]*template.Template)
	if templateFiles, err = fs.ReadDir(files, templatesDir); err != nil {
		panic(err)
	}

	for _, file := range templateFiles {
		if file.IsDir() {
			continue
		}

		// Each template will be accessible by its base name in the global map
		templates[file.Name()] = template.Must(template.ParseFS(files, filepath.Join(templatesDir, file.Name())))
	}
}

//===========================================================================
// Template Contexts
//===========================================================================

const (
	UnknownDate = "unknown date"
	DateFormat  = "Monday, January 2, 2006"
)

// EmailData includes data fields that are common to all the email builders such as the
// subject and sender/recipient information.
type EmailData struct {
	Subject   string           `json:"-"`
	Sender    sendgrid.Contact `json:"-"`
	Recipient sendgrid.Contact `json:"-"`
}

// Validate that all required data is present to assemble a sendable email.
func (e EmailData) Validate() error {
	switch {
	case e.Subject == "":
		return ErrMissingSubject
	case e.Sender.Email == "":
		return ErrMissingSender
	case e.Recipient.Email == "":
		return ErrMissingRecipient
	}
	return nil
}

// Build creates a new email from pre-rendered templates.
func (e EmailData) Build(text, html string) (msg *mail.SGMailV3, err error) {
	if err = e.Validate(); err != nil {
		return nil, err
	}

	return mail.NewSingleEmail(
		e.Sender.NewEmail(),
		e.Subject,
		e.Recipient.NewEmail(),
		text,
		html,
	), nil
}

// WelcomeData is used to complete the welcome email template
type WelcomeData struct {
	EmailData
	FirstName    string
	LastName     string
	Email        string
	Organization string
	Domain       string
}

// VerifyEmailData is used to complete the verify email template
type VerifyEmailData struct {
	EmailData
	FullName  string
	VerifyURL string
}

// InviteData is used to complete the invite email template
type InviteData struct {
	EmailData
	Email       string
	InviterName string
	OrgName     string
	Role        string
	InviteURL   string
}

// DailyUsersData is used to complete the daily users email template
type DailyUsersData struct {
	EmailData
	Date                time.Time         `json:"date"`
	InactiveDate        time.Time         `json:"inactive_date"`
	Domain              string            `json:"domain"`
	EnsignDashboardLink string            `json:"dashboard_url"`
	NewUsers            int               `json:"new_users"`
	DailyUsers          int               `json:"daily_users"`
	ActiveUsers         int               `json:"active_users"`
	InactiveUsers       int               `json:"inactive_users"`
	APIKeys             int               `json:"api_keys"`
	ActiveKeys          int               `json:"active_keys"`
	InactiveKeys        int               `json:"inactive_keys"`
	RevokedKeys         int               `json:"revoked_keys"`
	Organizations       int               `json:"organizations"`
	NewOrganizations    int               `json:"new_organizations"`
	Projects            int               `json:"projects"`
	NewProjects         int               `json:"new_projects"`
	NewAccounts         []*NewAccountData `json:"-"`
}

// NewAccountData describes user accounts that were created in the last 24 hours. An
// account is an instance of a user assigned to an organization. This includes the cases
// where a user registers and creates an organization, where a user is invited to an
// existing organization and creates an account, and where an existing user is invited
// to a second organization. The organization data is from the perspective of the entire
// organization not just the users' apikeys, projects, invitations, etc.
type NewAccountData struct {
	Name          string    `json:"name"`           // name of the user
	Email         string    `json:"email"`          // email address of the user
	EmailVerified bool      `json:"email_verified"` // if the user has verified their email address
	Role          string    `json:"role"`           // role of the user in the organization
	LastLogin     time.Time `json:"last_login"`     // timestamp the user logged in
	Created       time.Time `json:"created"`        // timestamp the user was added to the org
	Organization  string    `json:"organization"`   // name of the organization (workspace)
	Domain        string    `json:"domain"`         // domain of the organization
	Projects      int       `json:"projects"`       // number of projects in the organization
	APIKeys       int       `json:"apikeys"`        // number of api keys in the organization
	Users         int       `json:"users"`          // number of users in the organization
	Invitations   int       `json:"invitations"`    // number of user invitations in the organization
}

func (d DailyUsersData) TabTable() string {
	var builder strings.Builder
	w := tabwriter.NewWriter(&builder, 2, 3, 2, ' ', 0)
	fmt.Fprintf(w, "New Users:\t%d\tDaily Users:\t%d\n", d.NewUsers, d.DailyUsers)
	fmt.Fprintf(w, "Active Users:\t%d\tInactive Users:\t%d\n", d.ActiveUsers, d.InactiveUsers)
	fmt.Fprintf(w, "API Keys:\t%d\tRevoked API Keys:\t%d\n", d.APIKeys, d.RevokedKeys)
	fmt.Fprintf(w, "Active API Keys:\t%d\tInactive API Keys:\t%d\n", d.ActiveKeys, d.InactiveKeys)
	fmt.Fprintf(w, "New Organizations:\t%d\tOrganizations:\t%d\n", d.NewOrganizations, d.Organizations)
	fmt.Fprintf(w, "New Projects:\t%d\tProjects:\t%d\n", d.NewProjects, d.Projects)
	w.Flush()
	return builder.String()
}

func (d DailyUsersData) FormattedDate() string {
	return d.Date.Format(DateFormat)
}

func (d DailyUsersData) FormattedInactiveDate() string {
	return d.InactiveDate.Format(DateFormat)
}

func (d DailyUsersData) NewAccountsCSV() (_ []byte, err error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	// Write the header of the CSV file
	header := []string{
		"name", "email", "email_verified", "role", "last_login", "created",
		"organization", "domain", "projects", "apikeys", "users", "invitations",
	}
	if err = w.Write(header); err != nil {
		return nil, err
	}

	for _, account := range d.NewAccounts {
		row := []string{
			account.Name, account.Email, strconv.FormatBool(account.EmailVerified), account.Role,
			account.LastLogin.Format(time.RFC3339), account.Created.Format(time.RFC3339),
			account.Organization, account.Domain, strconv.Itoa(account.Projects),
			strconv.Itoa(account.APIKeys), strconv.Itoa(account.Users),
			strconv.Itoa(account.Invitations),
		}
		if err = w.Write(row); err != nil {
			return nil, err
		}
	}

	w.Flush()
	return buf.Bytes(), nil
}

//===========================================================================
// Email Builders
//===========================================================================

// WelcomeEmail creates a welcome email for a new user
func WelcomeEmail(data WelcomeData) (message *mail.SGMailV3, err error) {
	var text, html string
	if text, html, err = Render("welcome", data); err != nil {
		return nil, err
	}
	data.Subject = WelcomeRE
	return data.Build(text, html)
}

// VerifyEmail creates an email to verify a user's email address
func VerifyEmail(data VerifyEmailData) (message *mail.SGMailV3, err error) {
	var text, html string
	if text, html, err = Render("verify_email", data); err != nil {
		return nil, err
	}
	data.Subject = VerifyEmailRE
	return data.Build(text, html)
}

// InviteEmail creates an email to invite a user to join an organization
func InviteEmail(data InviteData) (message *mail.SGMailV3, err error) {
	var text, html string
	if text, html, err = Render("invite", data); err != nil {
		return nil, err
	}
	data.Subject = fmt.Sprintf(InviteRE, data.InviterName)
	return data.Build(text, html)
}

// DailyUsersEmail creates an email to send to admins that reports the PLG status
func DailyUsersEmail(data DailyUsersData) (message *mail.SGMailV3, err error) {
	var text, html string
	if text, html, err = Render("daily_users", data); err != nil {
		return nil, err
	}
	data.Subject = fmt.Sprintf(DailyUsersRE, data.Domain, data.Date.Format("January 2, 2006"))
	return data.Build(text, html)
}

//===========================================================================
// Template Builders
//===========================================================================

// Render returns the text and html executed templates for the specified name and data.
// Ensure that the extension is not supplied to the render method.
func Render(name string, data interface{}) (text, html string, err error) {
	if text, err = render(name+".txt", data); err != nil {
		return "", "", err
	}

	if html, err = render(name+".html", data); err != nil {
		return "", "", err
	}

	return text, html, nil
}

func render(name string, data interface{}) (_ string, err error) {
	var (
		ok bool
		t  *template.Template
	)

	if t, ok = templates[name]; !ok {
		return "", fmt.Errorf("could not find %q in templates", name)
	}

	buf := &strings.Builder{}
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// AttachData onto an email as a file with the specified mimetype
func AttachData(message *mail.SGMailV3, data []byte, filename, mimetype string) error {
	// Encode the data to attach to the email
	encoded := base64.StdEncoding.EncodeToString(data)

	// Create the attachment
	attach := mail.NewAttachment()
	attach.SetContent(encoded)
	attach.SetType(mimetype)
	attach.SetFilename(filename)
	attach.SetDisposition("attachment")
	message.AddAttachment(attach)
	return nil
}

// LoadAttachment onto email from a file on disk.
func LoadAttachment(message *mail.SGMailV3, attachmentPath string) (err error) {
	// Read and encode the attachment data
	var data []byte
	if data, err = os.ReadFile(attachmentPath); err != nil {
		return err
	}

	// Detect the mimetype either from the extension or using detect content type
	var mimetype string
	switch filepath.Ext(attachmentPath) {
	case ".zip":
		mimetype = "application/zip"
	case ".json":
		mimetype = "application/json"
	case ".csv":
		mimetype = "text/csv"
	case ".pdf":
		mimetype = "application/pdf"
	case ".tgz", ".gz":
		mimetype = "application/gzip"
	default:
		mimetype = http.DetectContentType(data)
	}

	// Create the attachment
	return AttachData(message, data, filepath.Base(attachmentPath), mimetype)
}

// AttachJSON by marshaling the specified data into human-readable data and encode and
// attach it to the email as a file.
func AttachJSON(message *mail.SGMailV3, data []byte, filename string) (err error) {
	return AttachData(message, data, filename, "application/json")
}

// AttachCSV by encoding the csv data and attaching it to the email as a file.
func AttachCSV(message *mail.SGMailV3, data []byte, filename string) (err error) {
	return AttachData(message, data, filename, "text/csv")
}

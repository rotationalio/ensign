package emails

import (
	"embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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
	Subject   string
	Sender    sendgrid.Contact
	Recipient sendgrid.Contact
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

// LoadAttachment onto email from a file on disk.
func LoadAttachment(message *mail.SGMailV3, attachmentPath string) (err error) {
	// Read and encode the attachment data
	var data []byte
	if data, err = os.ReadFile(attachmentPath); err != nil {
		return err
	}
	encoded := base64.StdEncoding.EncodeToString(data)

	// Create the attachment
	// TODO: detect mimetype rather than assuming zip
	attach := mail.NewAttachment()
	attach.SetContent(encoded)
	attach.SetType("application/zip")
	attach.SetFilename(filepath.Base(attachmentPath))
	attach.SetDisposition("attachment")
	message.AddAttachment(attach)
	return nil
}

// AttachJSON by marshaling the specified data into human-readable data and encode and
// attach it to the email as a file.
func AttachJSON(message *mail.SGMailV3, data []byte, filename string) (err error) {
	// Encode the data to attach to the email
	encoded := base64.StdEncoding.EncodeToString(data)

	// Create the attachment
	attach := mail.NewAttachment()
	attach.SetContent(encoded)
	attach.SetType("application/json")
	attach.SetFilename(filename)
	attach.SetDisposition("attachment")
	message.AddAttachment(attach)
	return nil
}

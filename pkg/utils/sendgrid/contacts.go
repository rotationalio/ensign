package sendgrid

const (
	Host     = "https://api.sendgrid.com"
	Contacts = "/v3/marketing/contacts"
)

type AddContact struct {
	ListIDs  []string   `json:"list_ids"`
	Contacts []*Contact `json:"contacts"`
}

type Contact struct {
	FirstName    string        `json:"first_name"`
	LastName     string        `json:"last_name"`
	Email        string        `json:"email"`
	Country      string        `json:"country"`
	CustomFields *CustomFields `json:"custom_fields"`
}

// TODO: make custom fields request to get field IDs rather than hardcoding.
type CustomFields struct {
	Title                string `json:"e1_T"`
	Organization         string `json:"e2_T"`
	CloudServiceProvider string `json:"e3_T"`
}

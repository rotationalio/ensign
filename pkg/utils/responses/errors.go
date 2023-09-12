package responses

// String errors intended to pass back from the server up to
// a **human** Beacon user so these need to make sense to our customers!
var (
	ErrTryLoginAgain             = response("Unable to login with those details - please try again!")
	ErrTryRegisterAgain          = response("Unable to register with those details - please try again!")
	ErrTryOrganizationAgain      = response("Unable to create or access that organization - please try again!")
	ErrTryProfileAgain           = response("Unable to create or access user profile - please try again!")
	ErrTryProjectAgain           = response("Unable to create or access that project - please try again!")
	ErrFixProjectDetails         = response("Unable to create a project with those details - please correct them and try again!")
	ErrNeedProjectPermission     = response("Unable to access project - please request permission from your team owner.")
	ErrMemberNotFound            = response("Team member with the specified ID was not found.")
	ErrMissingOrganizationName   = response("Organization name is required.")
	ErrMissingOrganizationDomain = response("Organization domain is required.")
	ErrDomainAlreadyExists       = response("An organization with that workspace already exists.")
	ErrOrganizationNotFound      = response("Organization with the specified ID was not found.")
	ErrTenantNotFound            = response("Tenant with the specified ID was not found.")
	ErrProjectNotFound           = response("Project with the specified ID was not found.")
	ErrTopicNotFound             = response("Topic with the specified ID was not found.")
	ErrLogBackIn                 = response("Logged out of your account - please log back in!")
	ErrVerifyEmail               = response("Please verify your email address and try again!")
	ErrVerificationFailed        = response("Email verification failed. Please contact support@rotational.io for assistance.")
	ErrRequestNewInvite          = response("Invalid invitation link - please request a new one!")
	ErrSomethingWentWrong        = response("Oops - something went wrong!")

	AllResponses = map[string]struct{}{}
)

// response creates a standard error message to ensure uniqueness and testability for
// external paackages
func response(msg string) string {
	if _, ok := AllResponses[msg]; ok {
		panic("duplicate error response defined: " + msg)
	}
	AllResponses[msg] = struct{}{}
	return msg
}

package responses

// String errors intended to pass back from the server up to
// a **human** Beacon user so these need to make sense to our customers!
var (
	ErrTryLoginAgain         = "Unable to login with those details - please try again!"
	ErrTryProjectAgain       = "Unable to create or access that project - please try again!"
	ErrFixProjectDetails     = "Unable to create a project with those details - please correct them and try again!"
	ErrNeedProjectPermission = "Unable to access project - please request permission from your team owner."
	ErrMemberNotFound        = "Team member with the specified ID was not found."
	ErrTenantNotFound        = "Tenant with the specified ID was not found."
	ErrProjectNotFound       = "Project with the specified ID was not found."
	ErrTopicNotFound         = "Topic with the specified ID was not found."
	ErrLogBackIn             = "Logged out of your account - please log back in!"
	ErrVerifyEmail           = "Please verify your email address and try again!"
	ErrVerificationFailed    = "Email verification failed. Please contact support@rotational.io for assistance."
	ErrRequestNewInvite      = "Invalid invitation link - please request a new one!"
	ErrSomethingWentWrong    = "Oops - something went wrong!"
)

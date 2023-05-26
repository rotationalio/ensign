package quarterdeck_test

import (
	"time"

	"github.com/rotationalio/ensign/pkg/utils/emails"
	"github.com/rotationalio/ensign/pkg/utils/emails/mock"
)

func (s *quarterdeckTestSuite) TestSendDailyUsers() {
	require := s.Require()
	defer mock.Reset()

	data := &emails.DailyUsersData{
		Date:                time.Date(2023, 4, 7, 0, 0, 0, 0, time.UTC),
		InactiveDate:        time.Date(2023, 3, 8, 0, 0, 0, 0, time.UTC),
		Domain:              "ensign.local",
		EnsignDashboardLink: "http://grafana.ensign.local/dashboards/ensign",
		NewUsers:            2,
		DailyUsers:          8,
		ActiveUsers:         102,
		InactiveUsers:       3,
		APIKeys:             58,
		ActiveKeys:          52,
		InactiveKeys:        6,
		RevokedKeys:         12,
		Organizations:       87,
		NewOrganizations:    1,
		Projects:            87,
		NewProjects:         1,
		NewAccounts: []*emails.NewAccountData{
			{
				Name:          "Wiley E. Coyote",
				Email:         "wiley@acme.co",
				EmailVerified: true,
				Role:          "owner",
				LastLogin:     "2023-04-06T19:21:39Z",
				Created:       "2023-04-06T12:02:52Z",
				Organization:  "Acme, Inc.",
				Domain:        "acme.co",
				Projects:      3,
				APIKeys:       7,
				Invitations:   3,
				Users:         2,
			},
			{
				Name:          "Julie Smith Lee",
				Email:         "jlee@foundations.io",
				EmailVerified: true,
				Role:          "member",
				LastLogin:     "2023-04-06T08:22:27Z",
				Created:       "2023-04-06T08:21:01Z",
				Organization:  "Foundations",
				Domain:        "foundations.io",
				Projects:      1,
				APIKeys:       1,
				Invitations:   8,
				Users:         1,
			},
		},
	}

	err := s.srv.SendDailyUsers(data)
	require.NoError(err, "could not send daily users report")

	// Check that there are two attachments
	messages := []*mock.EmailMeta{
		{
			To:          s.conf.SendGrid.AdminEmail,
			From:        s.conf.SendGrid.FromEmail,
			Subject:     "Daily PLG Report for : April 7, 2023",
			Attachments: 2,
		},
	}
	mock.CheckEmails(s.T(), messages)
}

func (s *quarterdeckTestSuite) TestSendDailyUsersNoNewAccounts() {
	require := s.Require()
	defer mock.Reset()

	data := &emails.DailyUsersData{
		Date:                time.Date(2023, 4, 7, 0, 0, 0, 0, time.UTC),
		InactiveDate:        time.Date(2023, 3, 8, 0, 0, 0, 0, time.UTC),
		Domain:              "ensign.local",
		EnsignDashboardLink: "http://grafana.ensign.local/dashboards/ensign",
		NewUsers:            0,
		DailyUsers:          8,
		ActiveUsers:         102,
		InactiveUsers:       3,
		APIKeys:             58,
		ActiveKeys:          52,
		InactiveKeys:        6,
		RevokedKeys:         12,
		Organizations:       87,
		NewOrganizations:    0,
		Projects:            87,
		NewProjects:         1,
	}

	err := s.srv.SendDailyUsers(data)
	require.NoError(err, "could not send daily users report")

	// Check that there is only 1 attachment
	messages := []*mock.EmailMeta{
		{
			To:          s.conf.SendGrid.AdminEmail,
			From:        s.conf.SendGrid.FromEmail,
			Subject:     "Daily PLG Report for : April 7, 2023",
			Attachments: 1,
		},
	}
	mock.CheckEmails(s.T(), messages)
}

Feature: Reset Password

I want to reset my password

Scenario: Resetting a password

Given I am on the reset password page
When I click the submit button without entering a new password
Then I should see error messages
When I enter a new password
And I do not enter the password confirmation
Then I should see an error message
When I enter a confirmation password that does not match the password
Then I should see an error message that the passwords do not match
When I enter a confirmation password that matches the password
And I click the submit button
Then I should be directed to the login page
And I should see a message that my password has been reset
When I log in with my new password
Then I should be directed to the dashboard

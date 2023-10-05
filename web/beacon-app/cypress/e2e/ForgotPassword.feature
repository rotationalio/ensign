Feature: Forgot Password

I want to submit a request to reset my password

Scenario: Completing the forgot password form

Given I am on the login page
When I click the forgot password link
Then I should be directed to the forgot password page
When I click the submit button without entering an email address
Then I should see a message informing me that an email address is required
When I enter an invalid email address
Then I should see a message informing me that the email address is invalid
When I enter a valid email address
And I click the submit button
Then I should be directed to the reset password verification page
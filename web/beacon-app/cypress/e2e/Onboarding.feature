Feature: New User Onboarding

I want to complete onboarding after I register a new account.

Scenario: New user onboarding

Given I'm on the login page
When I log into Beacon
Then I should be directed to the onboarding form
And I should see the onboarding sidebar
And I should see the first step of the onboarding form
And I should not see the Back button
When I click the next button without entering an organization name
Then I should see that the organization name is required
And I should not be able to continue to the second step
When I enter an organization name and click next
Then I should be directed to the second step of the onboarding form
And I should see the Back button
Then I should see a default workspace URL value
When I click the Back button
Then I should be directed to the first step of the onboarding form
And I should see the organization name that I entered
When I click next to return to the second step of the onboarding form
When I delete the default workspace URL
And I click next without entering a workspace URL
Then I should see that the workspace URL is required
And I should not be able to continue to the third step
When I enter a workspace URL and click next
Then I should be directed to the third step of the onboarding form
When I click the Back button on the third step of the onboarding form
Then I should be directed to the second step of the onboarding form
And I should see the workspace URL I entered
When I click to return to the third step of the onboarding form
And I click next without entering a name
Then I should see that the name is required
When I enter a name into the name input field and click next
Then I should be directed to the fourth step of the onboarding form
When I click the Back button on the fourth step of the onboarding form
Then I should be directed to the third step of the onboarding form
And I should see the name I entered
When I click to return to the fourth step of the onboarding form
Then I should see the profession segment options
And I should see the developer segment options
When I click next before selection a profession option or developer option
Then I should see that a profession segment option is required
And I should see that at least one developer segment option is required
When I select a profession option and not a developer option
And I click the next button to continue
Then I should see that at least one developer option is required
When I select 3 developer options
Then I should see that I cannot select any more developer options
When I click next to submit the onboarding form
Then I should be directed to the dashboard
And I should see the onboarding sidebar has been replaced with the regular sidebar

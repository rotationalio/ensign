Feature: Invited User Onboarding

I want to complete onboarding after receiving an invitation to join an organization.

Scenario: Onboarding for an invited user

Given I'm on the login page
When I log into Beacon
Then I should be directed to the onboarding form
And I should see the onboarding sidebar
And I should see my email address
When I click log out in the topbar
Then I should be directed to the login page
When I log in a second time
Then I should be directed back to the onboarding form
And I should see step 3 of the onboarding form
When I click next without entering a name
Then I should see that the name is required
When I enter a name into the name input field and click next
Then I should be directed to the fourth step of the onboarding form
When I click the Back button on the fourth step of the onboarding form
Then I should be directed back to the third step of the onboarding form
And I should see the name I entered
When I click to return to the fourth step of the onboarding form
Then I should see the professional segment options
And I should see the developer segment options
When I click next before selecting a professional option or developer option
Then I should see that a professional segment option is required
And I should see that at least one developer segment option is required
When I select a professional option and not a developer option
And I click the next button to continue
Then I should see that at least one developer option is required
When I select a first developer option
And I click a second developer option
And I click a third developer option
Then I should see that I cannot select any more developer options
When I click next to submit the onboarding form
Then I should be directed to the dashboard
And I should see the onboarding sidebar has been replaced with the regular sidebar
And I should see the name of the organization I joined
When I click the log out button
Then I should be directed to the login page
When I log into Beacon again
Then I should be directed to the dashboard and not see the onboarding workflow
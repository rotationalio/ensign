Feature: Team page

I want to navigate to the Team page

Scenario: Navigating the Team page

Given I login to Beacon
When I'm logged in
And I click team
Then I navigate to the team page
Then I should see the team table

And I should see the add team member button 
When I click the add team member button
Then I should see the invite new team member modal
When I click invite without entering an email address
Then I should see an error message
When I enter a valid email address
And I click the invite button
Then I should see a new team member appear in the teams table with a pending status

When I click the close button
Then I should not see the invite new team member modal

When I click the actions icon
Then I should see change role
When I click change role
Then I should see the change role modal
And the team member and current role should be populated
When I click save without changing the member's role I should see an error message
Then when I select a new role and click save I the member's role should update
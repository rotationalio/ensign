Feature: Team page

I want to navigate to the Team page

Scenario: Navigating the Team page

Given I'm logged into Beacon
When I navigate to the team page
Then I should see the team table
And I should see the team data
Then I should see the add team member button 
When I click the add team member button
Then I should see the invite new team member modal
When I click the close button
Then I should not see the invite new team member modal
When I click invite without entering an email address
Then I should see an error message
When I enter an invalid email address
Then I should see an error message
When I enter a valid email address
And I click the invite button
Then I should see a new team member appear in the teams table with a pending status
When I click the actions icon
Then I should see the actions menu
When I click change role
Then I should see the change role modal
When I click the close button
Then I should not see the change role modal
Then I should the team member and current role populated in the form when I want to change a role
When I click save without changing the member's role 
Then I should see an error message
When I select a new role and click save 
Then the member's role should update in the member table
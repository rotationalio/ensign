Feature: Team page

I want to navigate to the Team page

Scenario: Navigating the Team page

Given I'm logged into Beacon
And I click team
Then I navigate to the team page
Then I should see the team table
And I should see the team data
When I click the actions icon
And I click change role
Then I should see the change role modal
When I click the close button
Then I should not see the change role modal
Then I should see the team member and current role populated in the form when I want to change a role
When I click save without changing the member's role 
Then I should see a status error message
When I select a new role and click save 
Then the member's role should update in the member table
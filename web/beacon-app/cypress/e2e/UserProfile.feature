Feature: User profile page

I want to navigate to the user profile page

Scenario: Navigating the user profile page

Given I login to Beacon
When I'm logged in
And I click Profile
Then I navigate to the profile page
Then I should see the user profile
Then I should see the organizations table

When I click the cancel account button
Then I should see cancel account modal

When I click the close button
Then I should not see the cancel account modal


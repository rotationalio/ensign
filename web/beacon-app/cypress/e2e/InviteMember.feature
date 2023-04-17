Feature: Invite member

  I want to be able to invite a new member

  Scenario: Login to Beacon App
   
    Given I'm logged in
    When I navigate to the team page
    And I click on invite team member button
    And Add the member email address
    Then I should see the invited user in my team list



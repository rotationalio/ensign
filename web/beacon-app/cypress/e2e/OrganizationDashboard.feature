Feature: Display correct information on organization dashboard

  Scenario: Data are accurate
   
    Given I'm logged in to Beacon UI
    When I go to the organization dashboard
    Then I should see correct data

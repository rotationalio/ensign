Feature: Display correct information on organization dashboard

  Scenario: Data are accurate
   
    Given I'm logged in to Beacon UI
    When I go to the dashboard
    Then I should see correct data
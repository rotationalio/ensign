Feature: Project Detail

I want to navigate the project detail page

Scenario: Navigating the Beacon project detail page
    Given I'm  navigating to beacon app
    When I'm logged in
    And I Click on manage project button
    Then I'm redirected to project detail page
    Then I should see project detail 
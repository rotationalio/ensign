Feature: Beacon Main Page

I want to navigate the Beacon main page

Scenario: Navigating the Beacon main page
    Given I am on the Beacon homepage
    When I'm logged in
    Then I should see an avatar
    And I should see an org name
    And I should see Projects in the sidebar
    And I should see Team in the sidebar
    And I should see a link to Docs in the sidebar
    And I should see a link to Support in the sidebar
    And I should see a link to Profile in the sidebar
    And I should see a link to the About page in the sidebar footer
    Then I should be able to visit the About page if I click the link
    And I should see a link to the Contact Us page in the sidebar footer
    Then I should be able to visit the Contact Us page if I click the link
    And I should see a link to the Server Status page in the sidebar footer
    Then I should be able to visit the Server Status page if I click the link
    
    And I should see the Welcome component
    When I click the Start button
    Then I should see the create project modal
    When I click the close button
    Then I should not see the create project modal
    When I click on the Start button again
    And I fill in the project name
    And I fill in the project description
    And I click the Create Project button
    Then I should be redirected to the projects page
    And I should see the new project in the project list table
    Then I should see the number of projects increase to 1 on the projects page
    When I go back to the main page
    And I should see the number of projects increase to 1 on the main page

    And I should see the Set Up A New Project component
    And I should see the Access Resources component
    And I should see the Access button
    Then I should be able to visit the external documentation site
    
    When I click on the Create button
    Then I should see the create project modal
    When I fill in the project name
    And I fill in the project description
    And I click the Create Project button
    Then I should be redirected to the projects page
    And I should see the new project in the project list table
    And I should see the number of projects increase to 2 on the projects page
    When I go back to the main page
    Then I should see the number of projects increase increase to 2 on the main page
    And I should not see the Welcome component

    When I see the settings button
    And I click the settings button
    Then I should visit the settings page
    When I go back to the main page
    And I click the logout button
    Then I should log out of the Beacon home page

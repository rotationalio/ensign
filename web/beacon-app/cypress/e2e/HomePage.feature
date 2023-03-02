Feature: Beacon Main Page

I want to navigate the Beacon main page

Scenario: Navigating the Beacon main page
    Given I am on the Beacon homepage
    When I'm logged in
    Then I should see an avatar
    And I should see an org name
    And I should see a link to Docs in the sidebar
    And I should see a link to Profile in the sidebar
    And I should see a link to the About page in the sidebar footer
    Then I should be able to visit the About page if I click the link
    And I should see a link to the Contact Us page in the sidebar footer
    Then I should be able to visit the Contact Us page if I click the link
    And I should see a link to the Server Status page in the sidebar footer
    And I should see quickview data
    When I see the Manage project button
    Then I should not be able to click it
    And I should see the Create API Key button
    And I should see the View Docs button
    And I should be able to visit the external documentation site
    And I should see the settings button
    When I click the settings button
    Then I should visit the settings page
    And I should go back to the main page
    Then I should see the log out button
    When I click the logout button
    Then I should log out of the Beacon home page

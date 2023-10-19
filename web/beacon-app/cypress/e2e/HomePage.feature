Feature: Beacon Main Page

I want to navigate the Beacon main page

Scenario: Navigating the Beacon main page
    Given I am on the Beacon homepage
    When I'm logged in
    Then I should see an avatar
    And I should see an org name
    And I should see Projects in the sidebar
    And I should see Team in the sidebar
    And I should see a link to Profile in the sidebar

    And I should see a link to Ensign U in the sidebar
    And I should see a link to Use Cases in the sidebar
    And I should see a link to Docs in the sidebar
    And I should see a link to Data Playground in the sidebar
    And I should see a link to SDKs in the sidebar
    And I should see a link to Support in the sidebar

    And I should see a link to the About page in the sidebar footer
    Then I should be able to visit the About page if I click the link
    And I should see a link to the Contact Us page in the sidebar footer
    Then I should be able to visit the Contact Us page if I click the link
    And I should see a link to the Server Status page in the sidebar footer
    Then I should be able to visit the Server Status page if I click the link
    
    And I should see the Welcome component
    And I should see the welcome to Ensign video
    When I click on the welcome video
    Then I should see a modal open with a playable version of the video
    When I click the close button the modal
    Then I should not see the modal with the video

    And I should see the Set Up A New Project component
    And I should see the Create Project button
    When I click the Create Project button
    Then I should see the Create Project modal
    When I click the close button in the Create Project modal
    Then I should not see the Create Project modal

    And I should see the Starter Videos component
    And I should see thumbnails for the starter Videos
    When I click on a thumbnail
    Then I should see a modal open with a playable version of the video
    When I click the close button the modal
    Then I should not see the modal with the video

    And I should see the Schedule Office Hours icon
    Then I should see that I will be able to visit the Schedule Office Hours page if I click the icon
    When I see the settings button
    And I click the settings button
    Then I should see the settings page
    When I return to the main page
    And I click the logout button
    Then I should be logged out of the Beacon home page

Feature: Beacon Main Page

I want to navigate the Beacon main page

Scenario: Navigating the Beacon main page
    Given I am on the Beacon homepage
    When I'm logged in
    Then I should see the org name
    And I should see an avatar in the sidebar
    When I click on the avatar
    Then I should see a list of orgs I belong to
    
    When I click Projects in the sidebar
    Then I should be taken to the Projects page
    When I click Team in the sidebar
    Then I should be taken to the Team page
    When I click Profile in the sidebar
    Then I should be taken to the Profile page

    When I return to the home page
    Then I should see external links in the sidebar
    
    When I see the Welcome component
    Then I should see the Welcome to Ensign video
    When I click on the welcome video
    Then I should see a modal open with a playable version of the video
    When I click the close button to close the modal
    Then I should not see the modal with the video

    When I see the Set Up A New Project component
    And I click the Create Project button
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

    And I should see the Schedule Office Hours icon in the top bar
    Then I should see that I will be able to visit the Schedule Office Hours page if I click the icon
    And I should see the menu icon in the top bar
    When I click the memu icon
    Then I should see settings in the menu
    When I click settings
    Then I should be taken to the settings page
    When I return to the main page
    When I click the logout button in the menu
    Then I should be logged out of the Beacon home page

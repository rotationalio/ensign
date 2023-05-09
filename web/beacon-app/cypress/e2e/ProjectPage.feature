Feature: Project Page

I want to navigate the Project page

Scenario: Navigating the Project Page

Given I'm logged into Beacon
And I click Projects
Then I navigate to the project page
Then I should see the quick view data
And I should see the project list table
When I click on the Create Project button
Then I should see the create project modal
When I click the close button
Then I should not see the create project modal
When I click on the Create Project button again
And I fill in the project name
And I fill in the project description
And I click the Create Project button
Then I should see the new project in the project list table
And I should see the number of projects in the quick view data increase by 1
And when I go back to the home page
Then I should see the number of projects in the quick view data increase by 1

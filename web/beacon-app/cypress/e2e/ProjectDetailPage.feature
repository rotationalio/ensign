Feature: Project Detail Page

I want to navigate a project detail page

Scenario: Navigating a project detail page

Given I'm logged into Beacon
And I click Projects
Then I navigate to the project page
And I click on a project in the project table list
Then I should see the project detail page for the project
And I should see the project name at the top of the page
#And I should see the project's details when I hover over the hint icon next to the project name

When I see the project setup component
Then I should see that a project has been created

When I click on the cogwheel
Then I should see the project detail actions menu
When I click Delete Project
Then I should see the delete project modal
And I should not see the delete project modal when I click the close button

When I click Edit Project
Then I should see the Edit Project modal
And I should not see the Edit Project modal when I click the close button
When I re-open the Edit Project modal
Then I should see the current project's name
When I enter a new project name
And I change the project's description
And I click save
Then I should see the new project name
And I should see the updated project description

When I click Change Owner
Then I should see the Change Owner modal
And I should not see the Change Owner modal when I click the close button
When I change the project's owner
And I click the Save button
Then I should see the new project owner

When I see the API Keys component
#Then I should see more details about API keys when I hover over the hint icon
And I should see the API key list table
When I click the + New Key button
Then I should see the Generate API Key modal
When I click the close button I should not see the Generate API Key modal
When I re-open the Generate API Key modal
Then I should see an error if I try to create an API key without entering a Key Name
When I enter a Key Name
And I click the Generate API Key button
Then I should see the Your API Key modal
And I should be able to copy the Client ID
And I should be able to copy the Client Secret
And I should be able to download the API key details
When I confirm that have read the info on the Your API Key modal
Then I should see the new API key in the API key list table
And I should see that an API key has been created

When I see the Topics component
#Then I should see more details about topics when I hover over the hint icon
And I should see the topic list table
When I click the + New Topic button
Then I should see the New Topic modal
When I click the close button I should not see the New Topic modal
When I re-open the New Topic modal
Then I should see an error if I try to create a topic without entering a Topic Name
And I should see an error if I type an invalid Topic Name
When I enter a valid Topic Name
And I click the Create Topic button
Then I should see the new topic in the topic list table
And I should not see the project setup component

When I go back to the projects page
Then I should see the number of Topics increase by 1
And I should see the number of API keys increase by 1

When I go back to the main page
Then I should see the number of Topics increase by 1
And I should see the number of API keys increase by 1

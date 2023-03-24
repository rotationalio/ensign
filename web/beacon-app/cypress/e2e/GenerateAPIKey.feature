Feature: Generate API Key

I want to generate an API Key

Scenario: Generating an API Key workflow

Given I log in to Beacon
When I'm logged in
Then I should see the Create API Key button
When I click the Create API Key button
Then I should see the Generate API Key modal
And I should complete the Generate API Key modal
When I click the Generate API Key button
Then I should see the API Key confirmation modal
And I should be able to copy and download the client id and client secret
When I confirm that I have saved the client id and client secret
Then I should see a green check mark under the Create API Key button
And I should not be able to create another API key from the home page
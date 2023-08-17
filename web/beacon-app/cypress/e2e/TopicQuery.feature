Feature: Topic Query

I want to query EnSQL

Scenario: Submitting a query to EnSQL

Given I'm logged into Beacon
When I see the topic query section of the Topic Detail page
Then I should see the EnSQL link in the topic query instructions
And I should see the SDKs link in the topic query instructions

When I see the topic query input field
Then I should see the default topic query
And I should see the query button
And I should see the clear button
And I should see no query results
And I should see no viewing event results
And I should see an empty meta data table
And I should see NA listed for the mime type and event type
When I see the topic query results view
Then I should see the no query result message
And I should see disabled pagination buttons

When I click the query button
Then I should see 10 query results out of 11 total results
And I should view event 1 of 10
And I should see the mime type for event 1
And I should see the event type for event 1
And I should see the topic query result for event 1
And I should see that the previous button is disabled
And I should see that the next button is enabled

When I click the next button
Then I should view event 2 of 10
And I should see the mime type for event 2
And I should see the event type for event 2
And I should see that the previous button is enabled

When I click to view a result that could not be parsed
Then I should see the could not parse message next to the mime type
And I should see base 64 encoded data in the results view

When I click to the last result
Then I should see that the next button is disabled
And I should see that the previous button is still enabled
When I click the previous button
Then I should view event 9 of 10

When I click the clear button
Then I should see the default result view with the no query result message
And I should not see a value in the input field
And I should see that the pagination buttons are disabled
When I click query without typing a query
Then I should see the validation error message
When I type a query into the input field
Then I should not see the validation error message


Feature: Topic Detail Page

I want to navigate the Topic Detail Page

Scenario: Navigating the Topic Detail Page

When I am on the Topic Detail Page
Then I should see the topic name in the header component
When I hover over the topic detail tooltip
Then I should see the topic's details
When I move the cursor from the topic detail tooltip
Then I should not see the topic's details

And I should see the cogwheel icon in the header component
When I click the cogwheel icon
Then I should see a menu with menu items for Archive Topic, Delete Topic, and Clone Topic
When I click Archive Topic
Then I should see the Archive Topic modal
When I click x to close the Archive Topic modal
Then I should not see the Archive Topic Modal
When I click Delete Topic
Then I should see the Delete Topic modal
When I click x to close the Delete Topic modal
Then I should not see the Delete Topic Modal
When I click Clone Topic
Then I should see the Clone Topic modal
When I click x to close the Clone Topic modal
Then I should not see the Clone Topic Modal

And I should see 4 cards with metrics for the topic

And I should see the Topic Query compoent 
Then I should see the Topic Query carat toggle is open by default and pointed down
And I should see the Topic Query text instructions
When I click on the Topic Query title the carat toggle should be closed and pointed up
And I should still see the Topic Query text instructions
And I should not be able to see the Topic Query content
When I click on the Topic Query title again, the content should be visible

And I should see the Advanced Topic Policy Management compoent 
Then I should see the Advanced Topic Policy Management carat toggle is open by default and pointed down
And I should see the Advanced Topic Policy Management content
When I click on the Advanced Topic Policy Management title the carat toggle should be closed and pointed up
And I should not see the content
When I click on the Advanced Topic Policy Management title again, the content should be visible
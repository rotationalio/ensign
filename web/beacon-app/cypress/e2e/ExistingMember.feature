Feature: Accept member invitation

I want to accept member invitation

Scenario: Existing user when he has already an account

Given I've already an account
Then I should display login page



Scenario: Existing user when he hasn't an account

Given I've not an account
Then I should display registration page
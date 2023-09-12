Feature: Register as new tenant

  I want to be able to register as new tenant to Beacon App

  Scenario: Register to Beacon App

    Given I open the registration page
    When I click the Create Free Account button
    Then I should see the form error messages
    When I complete the registration form
    And I submit the registration form
    Then I should see the verify account page
    And I should see my email address in the verification email message

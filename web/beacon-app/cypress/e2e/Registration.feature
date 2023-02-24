Feature: Register as new tenant

  I want to be able to register as new tenant to Beacon App

  Scenario: Register to Beacon App

    Given I open the registration page
    When I fill correct informations
    And I submit the registration form

    Then I'm redirected on verify account page

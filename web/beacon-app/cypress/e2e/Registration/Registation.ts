import { And, Given, Then, When } from 'cypress-cucumber-preprocessor/steps';

Given('I open the registration page', () => {
  cy.intercept('http://localhost:8080/v1/register', {});
  cy.visit('/register');
});

When('I fill correct informations', () => {
  cy.findByRole('textbox', { name: /name \(required\)/i }).type('Holly Golythl');
  cy.findByRole('textbox', { name: /email address \(required\)/i }).type('holly.golth+1@ly.com');
  cy.findByTestId('password').type('holly.golth.ly1A');
  cy.findByTestId('pwcheck').type('holly.golth.ly1A');
  cy.findByTestId('organization').type('Team Diamonds');
  cy.findByTestId('domain').type('breakfast.tiffany.io');
  cy.get('[type="checkbox"]').check({ force: true });
});

And('I submit the registration form', () => {
  cy.get('form').submit();
});

Then("I'm redirected on verify account page", () => {
  cy.location('pathname').should('eq', '/verify-account');
});

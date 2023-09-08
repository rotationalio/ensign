import { And, Given, Then, When } from 'cypress-cucumber-preprocessor/steps';

beforeEach(function () {
  cy.fixture('user').then((user) => {
      this.user = user;
  });
})


Given('I open the registration page', () => {
  cy.visit('/register');
});

When('I click the Create Free Account button', () => {
  cy.get('[data-cy="submit-bttn"]').click();
});

Then('I should see the form error messages', () => {
  cy.get('[data-cy="email"]').siblings('div').should('have.text', 'The email address is required.');
  cy.get('[data-cy="password"]').siblings('div').should('have.text', 'The password is required.');
  cy.get('[data-cy="pwcheck"]').siblings('div').should('have.text', 'Please re-enter your password to confirm.');
});

When('I complete the reigstration form', function () {
  cy.get('[data-cy="email"]').type(this.user.email);
  cy.get('[data-cy="password"]').type(this.user.password);
  cy.get('[data-cy="pwcheck"]').type(this.user.password);
});

And('I submit the registration form', () => {
  cy.get('[data-cy="submit-bttn"]').click();
});

Then("I should see the verify account page", () => {
  cy.location('pathname').should('eq', '/verify-account');
});

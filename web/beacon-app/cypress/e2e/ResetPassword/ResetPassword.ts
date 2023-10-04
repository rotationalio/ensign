import { And, Given, Then, When } from 'cypress-cucumber-preprocessor/steps';

beforeEach(function () {
    cy.fixture('user').then((user) => {
        this.user = user;
    });
});

Given('I am on the reset password page', function () {
    cy.visit(`/${this.user.reset_password_url}`);
});

When('I click the submit button without entering a new password', () => {
    cy.get('[data-cy="reset-password-submit-bttn"]').should('exist').click();
});

Then('I should see error messages', () => {
    cy.get('[data-cy="password-error-msg"]')
    .should('exist')
    .and('have.text', 'Password is required.');

    cy.get('[data-cy="pwcheck-error-msg"]')
    .should('exist')
    .and('have.text', 'Please re-enter your password to confirm.');
});

When('I enter a new password', function () {
    cy.get('[data-cy="password"]')
    .should('exist')
    .type(this.user.password);
});

And('I do not enter the password confirmation', () => {
    cy.get('[data-cy="reset-password-submit-bttn"]').should('exist').click();
});

Then('I should see an error message', () => {
    cy.get('[data-cy="pwcheck-error-msg"]')
    .should('exist')
    .and('have.text', 'Please re-enter your password to confirm.');
});

When('I enter a confirmation password that does not match the password', function () {
    cy.get('[data-cy="pwcheck"]')
    .should('exist')
    .type(this.user.invalid_pwcheck);
});

Then('I should see an error message that the passwords do not match', () => {
    cy.get('[data-cy="pwcheck-error-msg"]')
    .should('exist')
    .and('have.text', 'The passwords must match.');
});

When('I enter a confirmation password that matches the password', function () {
    cy.get('[data-cy="pwcheck"]')
    .should('exist')
    .clear()
    .type(this.user.pwcheck);
});

And('I click the submit button', () => {
    cy.get('[data-cy="reset-password-submit-bttn"]').should('exist').click();
});

Then('I should be directed to the login page', () => {
    cy.location('pathname').should('eq', '/');
});

And('I should see a message that my password has been reset', () => {
    cy.findByRole('status')
      .should('have.text', 'Your password has been reset successfully. Please log in with your new password.')
});

When('I log in with my new password', function () {
    cy.loginWith({ email: this.user.email, password: this.user.password });
});

Then('I should be directed to the dashboard', () => {
    cy.location('pathname').should('eq', '/app');
});
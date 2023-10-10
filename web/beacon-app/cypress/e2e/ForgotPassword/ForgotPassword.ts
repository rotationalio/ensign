import { Given, Then, When, And } from "cypress-cucumber-preprocessor/steps";

beforeEach(function () {
    cy.fixture('user').then((user) => {
        this.user = user;
    });
});

Given('I am on the login page', () => {
    cy.visit('/');
});

When('I click the forgot password link', () => {
    cy.get('[data-cy="forgot-password-link"]').should('exist').click();
});

Then('I should be directed to the forgot password page', () => {
    cy.location('pathname').should('eq', '/forgot-password');
});

When('I click the submit button without entering an email address', () => {
    cy.get('[data-cy="forgot-password-submit-bttn"]').should('exist').click();
});

Then('I should see a message informing me that an email address is required', () => {
    cy.get('[data-cy="forgot-password-email-error"]')
    .should('exist')
    .and('have.text', 'Email is required.');
});

When('I enter an invalid email address', function () {
    cy.get('[data-cy="forgot-password-email-input"]')
    .should('exist')
    .type(this.user.invalid_email)
    .click();
});

Then('I should see a message informing me that the email address is invalid', () => {
    cy.get('[data-cy="forgot-password-email-error"]')
    .should('exist')
    .and('have.text', 'Please enter a valid email address.');
});

When('I enter a valid email address', function () {
    cy.get('[data-cy="forgot-password-email-input"]')
    .should('exist')
    .clear()
    .type(this.user.email);
});

And('I click the submit button', () => {
    cy.get('[data-cy="forgot-password-submit-bttn"]').should('exist').click();
});

Then('I should be directed to the reset password verification page', () => {
    cy.location('pathname').should('eq', '/reset-verification');
});
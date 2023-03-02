/// <reference types="cypress" />
/// <reference types="@testing-library/cypress" />
// ***********************************************
// This example commands.ts shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************
//
//
// -- This is a parent command --
// Cypress.Commands.add('login', (email, password) => { ... })
//
//
// -- This is a child command --
// Cypress.Commands.add('drag', { prevSubject: 'element'}, (subject, options) => { ... })
//
//
// -- This is a dual command --
// Cypress.Commands.add('dismiss', { prevSubject: 'optional'}, (subject, options) => { ... })
//
//
// -- This will overwrite an existing command --
// Cypress.Commands.overwrite('visit', (originalFn, url, options) => { ... })
//

import '@testing-library/cypress/add-commands';

Cypress.Commands.add('loginWith', ({ email, password }) =>
  cy
    .visit('/')
    .get('input[name="email"]')
    .type(email)
    .get('input[name="password"]')
    .type(password)
    .get('[data-testid="login-button"]')
    .click()
);

declare global {
  namespace Cypress {
    interface Chainable {
      loginWith({ email, password }: any): Chainable<JQuery<HTMLElement>>;
    }
  }
}

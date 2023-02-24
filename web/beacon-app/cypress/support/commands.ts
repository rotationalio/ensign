/// <reference types="cypress" />
/// <reference types="@testing-library/cypress" />
/// <reference types="cypress-localstorage-commands" />
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
// declare global {
//   namespace Cypress {
//     interface Chainable {
//       loginWith(data: any): Chainable<void>;
//       // drag(subject: string, options?: Partial<TypeOptions>): Chainable<Element>
//       // dismiss(subject: string, options?: Partial<TypeOptions>): Chainable<Element>
//       // visit(originalFn: CommandOriginalFn, url: string, options: Partial<VisitOptions>): Chainable<Element>
//     }
//   }
// }

import '@testing-library/cypress/add-commands';

Cypress.Commands.add('login', () => {
  cy.visit('/');
  cy.get('form').within(() => {
    cy.fixture('user.json').then((user) => {
      cy.get('input[name="email"]').type(user.email);
      cy.get('input[name="password"]').type(user.password);
      cy.root().submit();
    });
  });

  cy.location('pathname').should('eq', '/app');
});

Cypress.Commands.add('getBySel', (selector, ...args) => {
  return cy.get(`[data-cy=${selector}]`, ...args);
});

declare global {
  namespace Cypress {
    interface Chainable {
      login(): Chainable<void>;
      register(data: { email: string; password: string; username: string }): Chainable<void>;
      getBySel(dataTestAttribute: string, args?: any): Chainable<JQuery<HTMLElement>>;

      //   drag(subject: string, options?: Partial<TypeOptions>): Chainable<Element>
      //   dismiss(subject: string, options?: Partial<TypeOptions>): Chainable<Element>
      //   visit(originalFn: CommandOriginalFn, url: string, options: Partial<VisitOptions>): Chainable<Element>
    }
  }
}

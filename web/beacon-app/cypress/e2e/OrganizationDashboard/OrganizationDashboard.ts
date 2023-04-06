/* eslint-disable testing-library/await-async-query */
/* eslint-disable testing-library/prefer-screen-queries */
/* eslint-disable testing-library/no-debugging-utils */
import { Given, Then, When } from 'cypress-cucumber-preprocessor/steps';

Given("I'm logged in to Beacon UI", () => {
  cy.login();
});

When('I go to the organization dashboard', () => {
  cy.visit('/app/organization');
});

Then('I should see correct data', () => {
  cy.get('[data-cy="card_list_item_table"]').within(($subject) => {
    cy.wrap($subject).find('tr').should('have.length', 4);

    cy.findByTestId('cardlistitem-0').should('exist');
    cy.findByTestId('cardlistitem-1').should('exist');
    cy.findByTestId('cardlistitem-2').should('exist');
    cy.findByTestId('cardlistitem-3').should('exist');
  });

  cy.getBySel('tenants_table').within(($subject) => {
    cy.wrap($subject)
      .find('tbody')
      .within(($tr) => {
        cy.wrap($tr).find('tr').should('have.length', 1);

        cy.wrap($tr).find('td').should('have.length', 3).eq(0).should('exist');
        cy.wrap($tr).find('td').should('have.length', 3).eq(1).should('exist');
        cy.wrap($tr).find('td').should('have.length', 3).eq(2).should('exist');
      });
  });
});

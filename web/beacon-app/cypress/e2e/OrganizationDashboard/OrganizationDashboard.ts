/* eslint-disable testing-library/await-async-query */
/* eslint-disable testing-library/prefer-screen-queries */
/* eslint-disable testing-library/no-debugging-utils */
import { Given, Then, When } from 'cypress-cucumber-preprocessor/steps';

Given("I'm logged in to Beacon UI", () => {
  cy.login();
});

When('I go to the organization dashboard', () => {
  cy.intercept('GET', 'http://localhost:8080/v1/tenant', (req) => {
    // Retourner les données depuis vos fixtures
    req.reply({
      fixture: 'tenants.json',
    });
  }).as('tenants');

  cy.intercept('GET', 'http://localhost:8080/v1/organization/01GVZPCHMB28KM7BPPWH5R48HW', (req) => {
    // Retourner les données depuis vos fixtures
    req.reply({
      fixture: 'organization.json',
    });
  }).as('organization');

  cy.intercept(
    'GET',
    'http://localhost:8080/v1/tenant/01GVZPCHRJ3J6WRWPY1YGQ1J13/projects',
    (req) => {
      // Retourner les données depuis vos fixtures
      req.reply({
        fixture: 'projects.json',
      });
    }
  ).as('projects');

  cy.visit('/app/organization');
});

Then('I should see correct data', () => {
  cy.get('[data-cy="card_list_item_table"]').within(($subject) => {
    cy.wrap($subject).find('tr').should('have.length', 4);

    cy.fixture('organization.json').then((organization) => {
      cy.findByTestId('cardlistitem-0').should('have.text', organization.name);
      cy.findByTestId('cardlistitem-1').should('have.text', organization.id);
      cy.findByTestId('cardlistitem-2').should('have.text', organization.owner);
      cy.findByTestId('cardlistitem-3').should('have.text', '20/03/2023 14:22');
    });
  });

  cy.getBySel('tenants_table').within(($subject) => {
    cy.wrap($subject)
      .find('tbody')
      .within(($tr) => {
        cy.wrap($tr).find('tr').should('have.length', 1);

        cy.wrap($tr)
          .find('td')
          .should('have.length', 3)
          .eq(0)
          .should('have.text', 'alford-and-klein-co');
        cy.wrap($tr).find('td').should('have.length', 3).eq(1).should('have.text', 'development');
        cy.wrap($tr)
          .find('td')
          .should('have.length', 3)
          .eq(2)
          .should('have.text', '20/03/2023 14:22');
      });
  });
});

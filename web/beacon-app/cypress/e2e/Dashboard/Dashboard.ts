/* eslint-disable testing-library/await-async-query */
/* eslint-disable testing-library/prefer-screen-queries */
/* eslint-disable testing-library/no-debugging-utils */
import { Given, Then, When } from 'cypress-cucumber-preprocessor/steps';

Given("I'm logged in to Beacon UI", () => {
  cy.login();
});

When('I go to the dashboard', () => {
  cy.intercept('GET', 'http://localhost:8080/v1/tenant', (req) => {
    req.reply({
      fixture: 'tenants.json',
    });
  }).as('tenants');
  cy.intercept('GET', 'http://localhost:8080/v1/organization/01GVZPCHMB28KM7BPPWH5R48HW', (req) => {
    req.reply({
      fixture: 'organization.json',
    });
  }).as('organization');
  cy.intercept(
    'GET',
    'http://localhost:8080/v1/tenant/01GVZPCHRJ3J6WRWPY1YGQ1J13/projects',
    (req) => {
      req.reply({
        fixture: 'projects.json',
      });
    }
  ).as('projects');

  cy.intercept('GET', 'localhost:8080/v1/tenant/01GVZPCHRJ3J6WRWPY1YGQ1J13/stats', (req) => {
    req.reply({
      fixture: 'stats.json',
    });
  }).as('stats');
});

Then('I should see correct data', () => {
  cy.get('[data-cy="quick_view"]').within(($subject) => {
    cy.fixture('stats.json').then((stats) => {
      cy.wrap($subject)
        .find('div')
        .should('have.length', 4)
        .each(($item, $index) => {
          cy.wrap($item)
            .get(`[data-cy=${stats[$index].name}_value]`)
            .should('have.text', stats[$index].value);
        });
    });
  });

  cy.get('[data-testid="manage"]').click();
  cy.location('pathname').should('eq', '/app/projects/01GVZPCGX8F4V1DP4E8EMPAVSM');
});

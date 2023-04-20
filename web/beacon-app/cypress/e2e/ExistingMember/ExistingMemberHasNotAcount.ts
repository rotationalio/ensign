import { Given, Then } from 'cypress-cucumber-preprocessor/steps';

Given("I've not an account", () => {
  cy.visit('/invite?token=L-2BLXocLL-wpe4yOcd_otf6d-vHt0Zs8wPGgqFQCJU');
});

Then('I should display registration page', () => {
  // eslint-disable-next-line testing-library/await-async-query, testing-library/prefer-screen-queries
  cy.findByText("You've Been Invited!").should('exist');
  // cy.getBySel('cy.get('[data-testid="email"]')')
  cy.get('[data-testid="inviter_name"]').should('exist');
  cy.get('[data-testid="org_name"]').should('exist');
  cy.get('[data-testid="role"]').should('exist');
  cy.get('input').should('have.length', 5);
});

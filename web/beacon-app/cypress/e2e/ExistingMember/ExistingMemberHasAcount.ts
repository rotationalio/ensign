import { Given, Then } from 'cypress-cucumber-preprocessor/steps';

Given("I've already an account", () => {
  cy.visit('invite?token=nBW0OQr8yQxRWxo4aVGPCTIBsEjgaXIH5dlnJ2IWzkU');
});

Then('I should display login page', () => {
  // eslint-disable-next-line testing-library/await-async-query, testing-library/prefer-screen-queries
  cy.findByText("You've Been Invited!").should('exist');
  // cy.getBySel('cy.get('[data-testid="email"]')')
  cy.get('[data-testid="inviter_name"]').should('exist');
  cy.get('[data-testid="org_name"]').should('exist');
  cy.get('[data-testid="role"]').should('exist');
});

import { And, Given, Then, When } from 'cypress-cucumber-preprocessor/steps';

Given('I am on the Beacon app', () => {
  cy.visit('/');
});

When("I'm logged in", () => {
  cy.loginWith({ email: 'test@ensign.com', password: 'Abc123Def$56' });
  cy.url().should('include', 'app');
  cy.getCookies().should('exist');
  cy.getCookie('bc_atk').should('exist');
});

And('I Click on manage project button', () => {
  cy.get('[data-testid="manage"]').click();
});

Then("I'm redirected to project detail page", () => {
  // get project id from data-testid and check if the url is correct with the project id

  cy.location('pathname').should('include', 'app/projects');
});

Then('I see the project detail page', () => {
  cy.get('[data-testid="project-detail"]').should('exist');
});

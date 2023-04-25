import { And, Given, Then, When } from 'cypress-cucumber-preprocessor/steps';

Given("I'm logged in", () => {
  cy.login();
});

When('I navigate to the team page', () => {
  cy.visit('/app/team');
});

And('I click on invite team member button', () => {
  cy.getBySel('add-team-member').click();
});

And('Add the member email address', () => {
  cy.getBySel('member-email-address').type('fake@email.com');
  cy.get('[data-cy="member-email-address"]').type('fake@email.com');
  cy.get('form').submit();
});

Then('I should see the invited user in my team list', () => {
  // eslint-disable-next-line testing-library/prefer-screen-queries, testing-library/await-async-query
  cy.findByText('fake@email.com').should('exist');
});

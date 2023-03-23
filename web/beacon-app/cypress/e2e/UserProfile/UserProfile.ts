import { Given, Then, When, And } from "cypress-cucumber-preprocessor/steps";

Given('I login to Beacon', () => {
    cy.loginWith( { email: 'test@testing.com', password:'Abc123Def$56'})
  });

  When('I\'m logged in', () => {
    cy.url().should('include', 'app')
      cy.getCookies().should('exist')
      cy.getCookie('bc_atk').should('exist')
  });

  And('I click Profile', () => {
    cy.contains('a', 'Profile').click()
  })

  Then('I navigate to the profile page', () => {
    cy.location('pathname').should('include', 'app/profile')
  })

  Then('I should see the user profile', () => {
    cy.get('[data-testid="cardlistitem"]').should('exist')
  })

  Then('I should see the organizations table', () => {
    cy.get('[data-testid="orgTable"]').should('exist')
  })

  When('I click the cancel account button', () => {
    cy.get('[data-testid="blueBars"]').click()
    cy.get('[data-testid="cancelButton"]').click()
  })

  Then('I should see cancel account modal', () => {
    cy.get('[data-testid="cancelAcctModal"]').should('exist')
  })

  When('I click the close button', () => {
    cy.get('[data-testid="closeButton"]').click()
  })

  Then('I should not see the cancel account modal', () => {
    cy.get('[data-testid="cancelAcctModal"]').should('not.exist')
  })


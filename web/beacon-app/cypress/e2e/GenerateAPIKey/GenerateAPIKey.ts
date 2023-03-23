import { Given, Then, When, And } from "cypress-cucumber-preprocessor/steps";

Given('I log in to Beacon', () => {
    cy.loginWith( { email: 'test@test.com', password:'Abc123Def$56'})
  });

  When('I\'m logged in', () => {
    cy.url().should('include', 'app')
      cy.getCookies().should('exist')
      cy.getCookie('bc_atk').should('exist')
  });

  Then('I should see the Create API Key button', () => {
    cy.get('[data-testid="apikey"]').should('exist')
    cy.get('[data-testid="key"]').should('be.enabled')
  })

  When('I click the Create API Key button', () => {
    cy.get('[data-testid="key"]').click()
  })

  Then('I should see the Generate API Key modal', () => {
    cy.get('[data-testid="keyModal"]').should('exist')
  })

  And('I should complete the Generate API Key modal', () => {
    cy.get('[data-testid="keyName"]').type('Test')
    cy.get('[data-testid="generateKey"]').click()
  })

  When('I click the Generate API Key button', () => {
    cy.get('[data-testid="generateKey"]').click()
  })

  Then('I should see the API Key confirmation modal', () => {
    cy.get('[data-testid="keyCreated"]').should('exist')
  })

  And('I should be able to copy and download the client id and client secret', () => {
    cy.get('[data-testid="clientId"]').should('exist')
    cy.get('[data-testid="copyID"]').should('exist')
    cy.get('[data-testid="clientSecret"]').should('exist')
    cy.get('[data-testid="copySecret"]').should('exist')
    cy.get('[data-testid="download"]').should('exist')
  })

  When('I confirm that I have saved the client id and client secret', () => {
    cy.get('[data-testid="closeKey"]').click()
  })

  Then('I should see a green check mark under the Create API Key button', () => {
    cy.get('[data-testid="checkmark"]').should('exist')
  })

  And('I should not be able to create another API key from the home page', () => {
    cy.get('[data-testid="key"]').should('not.be.enabled')
  })
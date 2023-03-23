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

  And('I should see the user data', () => {
    cy.get('tr>td').eq(0).should('have.text', 'User ID')
    cy.get('[data-testid="cardlistitem-0"]').should('exist')
    cy.get('tr>td').eq(2).should('have.text', 'Name')
    cy.get('[data-testid="cardlistitem-1"]').should('have.text', 'Danielle Maxwell')
    cy.get('tr>td').eq(4).should('have.text', 'Role')
    cy.get('[data-testid="cardlistitem-2"]').should('have.text', 'Owner')
    cy.get('tr>td').eq(6).should('have.text', 'Date Created')
    cy.get('[data-testid="cardlistitem-3"]').should('have.text', '03/20/2023, 05:52 PM')
  })

  Then('I should see the organizations table', () => {
    cy.get('[data-testid="orgTable"]').should('exist')
  })

  And('I should see the organization data', () => {
    cy.get('tr>th>div').eq(0).should('have.text', "Organization ID")
    cy.get('tr>td').eq(8).should('exist')
    cy.get('tr>th>div').eq(1).should('have.text', "Organization Name")
    cy.get('tr>td').eq(9).should('have.text', 'X Team')
    cy.get('tr>th>div').eq(2).should('have.text', "Organization Owner")
    cy.get('tr>td').eq(10).should('have.text', 'Danielle Maxwell')
    cy.get('tr>th>div').eq(3).should('have.text', "Projects")
    cy.get('tr>td').eq(11).should('have.text', '1')
    cy.get('tr>th>div').eq(4).should('have.text', "Date Created")
    cy.get('tr>td').eq(12).should('have.text', '03/20/2023, 05:52 PM')
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


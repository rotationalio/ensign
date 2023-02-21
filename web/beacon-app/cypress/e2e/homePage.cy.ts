describe('Beacon dashboard', () => {
  beforeEach(() => {
    cy.visit('/');
    cy.findByTestId('email').type('test@ensign.com')
    cy.findByTestId('password').type('Abc123Def$56')
    cy.findByTestId('login').click()
  });
  
  it('displays an avatar in sidebar', () => {
    cy.findByTestId('gravatar').should('exist')
    cy.findByTestId('gravatar').invoke('attr', 'src').should('contain', 'https://www.gravatar.com/')
  })

  // TODO: Add test to ensure org name matches org name included in registration
 it('displays organization name in sidebar', () => {
    cy.findByTestId('orgName').should('exist')
  })
  
  it('displays quickview data', () => {
    
    const data = [
      'Active Projects',
      'Topics',
      'API Keys',
      'Data Storage'
    ]
  
    cy.get('h5').each(($e, index) => {
      cy.wrap($e).should('have.text', data[index])
      cy.wrap($e).siblings().should('not.be.empty')
    })
  })

  it('displays not allowed cursor when user hovers over manage project button', () => {
    cy.findByTestId('manage').should('have.class', 'cursor-not-allowed')
  })

  it('displays error message when project data is not available', () => {
    cy.findByText(/No data available, please try again later or contact support./i).should('have.text', 'No data available, please try again later or contact support.')
  })

  it('displays view documentation button with a link to external documentation site', () => {
    cy.contains('a', 'View Docs').should('have.text', 'View Docs').and('have.attr', 'href', 'https://ensign.rotational.dev/getting-started/')
  })

  it('navigates user to settings page', () => {
    cy.findByTestId('menu').click()
    cy.findByText(/settings/i).click()
    cy.location('pathname').should('eq', '/app/organization')
  })

  it('navigates user to log in page after log out', () => {
    cy.findByTestId('menu').click()
    cy.findByText(/logout/i).click()
    cy.location('pathname').should('eq', '/')
  })

/*    it('includes link in nav bar that navigates users to external documentation site', () => {
    cy.contains('a', 'Docs').should('have.text', 'Docs ').and('have.attr', 'href', 'https://ensign.rotational.dev/getting-started/')
  })

  it('includes link in nav bar that navigates users to profile', () => {
    cy.contains('a', 'Profile').should('have.text', 'Profile ').and('have.attr', 'href', '/app/profile')
  }) */

 /*   it('includes link in nav bar that navigates user to Rotational about page', () => {
    cy.contains('a', 'About').should('have.text', 'About').and('have.attr', 'href', 'https://rotational.io/about')
  })

  it('includes link in nav bar that navigates user to Rotational contact us page', () => {
  cy.contains('a', 'Contact Us').should('have.text', 'Contact Us').and('have.attr', 'href', 'https://rotational.io/contact')
  }) */

  // TODO: Test support link
  
  // TODO: Test server status link
})
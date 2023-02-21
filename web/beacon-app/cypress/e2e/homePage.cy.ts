describe('Beacon dashboard', () => {
  beforeEach(() => {
    cy.visit('/');
    cy.findByTestId('email').type('danielle@rotational.io')
    cy.findByTestId('password').type('Abc123Def$56')
    cy.findByTestId('login').click()
  });
  
  it('displays an avatar in sidebar', () => {
    cy.findByRole('img').should('exist')
  })

 /* it('displays user name in sidebar', () => {
    cy
    .findByTestId('')
    .should('exist')
    .and('')
  }) */

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

   /*  it('includes link in nav bar that navigates users to external documentation site', () => {
    cy.findByTestId('').should('exist').click()
    cy.visit('https://ensign.rotational.dev/getting-started/')
  }) */

  // TODO: Test support link

  // TODO: Test profile link

 /*  it('includes link in nav bar that navigates user to Rotational about page', () => {
    cy.findByText(/About/i).should('exist').click()
    cy.visit('https://rotational.io/about')
  })

  it('includes link in nav bar that navigates user to Rotational contact us page', () => {
    cy.findByText(/Contact/i).should('exist').click()
    cy.visit('https://rotational.io/contact')
  }) */
  
  // TODO: Test server status link

})
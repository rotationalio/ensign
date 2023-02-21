describe('Beacon dashboard', () => {
  beforeEach(() => {
    cy.visit('/app');
    cy.findByTestId('email').type('danielle@rotational.io')
    cy.findByTestId('password').type('Abc123Def$56')
    cy.findByTestId('login').click()
  });
  
  it('displays an avatar in sidebar', () => {
    cy.findByRole('img').should('exist')
  })

 /*  it('displays user name in sidebar', () => {
    cy
    .findByTestId('')
    .should('exist')
    .and('')
  }) */

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

  it('displays not allowed cursor when user hovers over manage project button', () => {
    cy.findByTestId('manage').should('have.class', 'cursor-not-allowed')
  })

  it('displays error message when project data is not available', () => {
    cy.findByText(/No data available, please try again later or contact support./i).should('have.text', 'No data available, please try again later or contact support.')
  })

  it('displays view documentation button navigates users to external documentation site when view docs is clicked', () => {
    cy.findByText(/View Docs/i).should('exist').click()
    cy.visit('https://ensign.rotational.dev/getting-started/')
  })

  it('navigates user to settings page', () => {
    cy.findByTestId('menu').click()
    cy.findByText(/settings/i).click()
    cy.visit('/app/organizations')
    // TODO: Test return to dashboard from settings page
  })

  it('logs out user and navigates to log in page and returns to main dashboard page', () => {
    cy.findByTestId('menu').click()
    cy.findByText(/logout/i).click()
    cy.visit('/app')
  })

  // TODO: Test stats data in Quick view

  // TODO: Test server status link

  // TODO: Test Create API Key

})
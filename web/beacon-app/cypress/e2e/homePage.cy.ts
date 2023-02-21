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

  it('displays not allowed cursor when user hovers over manage project button', () => {
    cy.findByTestId('manage').should('have.class', 'cursor-not-allowed')
  })

  it('displays error message when project data is not available', () => {
    cy.findByText(/No data available, please try again later or contact support./i).should('have.text', 'No data available, please try again later or contact support.')
  })

  it('displays view documentation button', () => {
    cy.findByText(/View Docs/i).should('exist')
  })

  it('navigates users to external documentation site when view docs is clicked', () => {
    cy.findByText(/View Docs/i).click()
    cy.visit('https://ensign.rotational.dev/getting-started/')
  })

/*   it('logs out user and navigates to log in page', () => {
    cy
    .findByTestId('menu')
    .click()
    .findByText(/logout/i)
    .click()
    .visit('/')
  }) */

  // TODO: Test server status link

})
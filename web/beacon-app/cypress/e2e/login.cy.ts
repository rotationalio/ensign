describe('Beacon login page', () => {
  beforeEach(() => {
    cy.visit('/');
  });

  it('displays email and password input fields and log in button', () => {
    cy.findByTestId('email').should('exist')
    cy.findByTestId('password').should('exist')
    cy.findByTestId('login').should('have.text', 'Log in').and('not.be.disabled')
  })

  it('navigates users without an account to getting started page', () => {
    cy.findByTestId('register').click()
    cy.location('pathname').should('eq', '/register')
  })

  it('returns an error if log in details are not provided', () => {
    cy.findByTestId('login').click()
    cy.findByText(/Email is required/i).should('have.text', "Email is required")
    cy.findByText(/Password is required/i).should('have.text', "Password is required")
  })

  it('returns an error if email address is not valid', () => {
    cy.findByTestId('email').type('test')
    cy.findByTestId('password').click()
    cy.findByText(/Email is invalid/i).should('have.text', "Email is invalid")
  })

  it('returns an error if email address is valid but password is not provided', () => {
    cy.findByTestId('email').type('test@example.com')
    cy.findByTestId('login').click()
    cy.findByText(/Password is required/i).should('have.text', "Password is required")
  })

  it('directs users to dashboard after successful login', () => {
    cy.findByTestId('email').type('danielle@rotational.io')
    cy.findByTestId('password').type('Abc123Def$56')
    cy.findByTestId('login').click()
    cy.location('pathname').should('eq', '/app')
  })
})
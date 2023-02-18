import cypressConfig from "../../cypress.config"


describe('template spec', () => {
  it('opens base url', () => {
    cy.visit('/app')
  })
})
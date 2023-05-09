import { And, Given, Then, When } from 'cypress-cucumber-preprocessor/steps';

Given("I'm logged into Beacon", () => {
    cy.loginWith( { email: 'test2@test.com', password:'Abc123Def$56'})
});

And("I click Projects", () => {
    cy.contains('div', 'Projects').should('exist').click()
});

Then("I navigate to the project page", () => {
    cy.location('pathname').should('include', 'app/projects')
});

Then("I should see the quick view data", () => {
    cy.contains('h5', 'Projects').should('exist')
    cy.contains('h5', 'Topics').should('exist')
    cy.contains('h5', 'Keys').should('exist')
    cy.contains('h5', 'Storage').should('exist')
});

And("I should see the project list table", () => {
    cy.get('[data-testid="projectsTable"]').should("exist")
    cy.get('th>div').eq(0).should('have.text', 'Project Name')
    cy.get('th>div').eq(1).should('have.text', 'Description')
    cy.get('th>div').eq(2).should('have.text', 'Status')
    cy.get('th>div').eq(3).should('have.text', 'Date Created')
});

When("I click on the Create Project button", () => {
    cy.get('[data-testid="create-project-btn"]').should("exist").click()
});

Then("I should see the create project modal", () => {
    cy.get('[data-testid="newProjectModal"]').should("exist")
});

/* When("I click the close button", () => {
    cy.get('div>button>svg').eq(1).click()
});

Then("I should not see the create project modal", () => {
    cy.get('[data-testid="newProjectModa"]').should("not.exist")
});
 
When("I click on the Create Project button again", () => {
    cy.get('[data-testid="create-project-btn"]').should("exist").click()
}); */

And("I fill in the project name", () => {
    cy.get('[data-cy="project-name"]').should('exist').type('My first project')
});

And("I fill in the project description", () => {
    cy.get('[data-cy="project-description"]').should('exist').type('A new project to test Ensign!')
});

And("I click the Create Project button", () => {
    cy.get('[data-cy="NewProjectButton"]').should("exist").click()
});

Then("I should see the new project in the project list table", () => {
    cy.get('tbody>tr>td').eq(4).should('have.text', 'My first project')
    cy.get('tbody>tr>td').eq(5).should('have.text', 'A new project to test Ensign!')
    cy.get('tbody>tr>td').eq(7).should('have.text', '05/09/2023')
});

And("I should see the number of projects in the quick view data increase by 1", () => {
    cy.contains('p', '3').should('exist')
});

When("I go back to the home page", () => {
    cy.go('back')
    cy.location('pathname').should('include', 'app')
});

Then("I should see the number of projects in the quick view data increase by 1", () => {
    cy.contains('p', '3').should('exist')
});
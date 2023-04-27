import { Given, And, Then, When } from "cypress-cucumber-preprocessor/steps";

Given("I'm logged into Beacon", () => {
    cy.loginWith( { email: 'test1@test.com', password:'Abc123Def$56'})
});

And("I click team", () => {
    cy.contains('div', 'Team').should('exist').click()
});

Then('I navigate to the team page', () => {
    cy.location('pathname').should('include', 'app/team')
  })

Then("I should see the team table", () => {
    cy.get('[data-testid="teamTable"]').should("exist");
});

And("I should see the team data", () => {
    cy.get('th>div').eq(0).should('have.text', 'Name')
    cy.get('th>div').eq(1).should('have.text', 'Email Address')
    cy.get('th>div').eq(2).should('have.text', 'Role')
    cy.get('th>div').eq(3).should('have.text', 'Status')
    cy.get('th>div').eq(4).should('have.text', 'Last Activity')
    cy.get('th>div').eq(5).should('have.text', 'Joined Date')
    cy.get('th>div').eq(6).should('have.text', 'Actions')
});

When("I click the actions icon", () => {
    cy.get('td>div>button').eq(0).click();
});

And("I click change role", () => {
    cy.get('div>ul>li').eq(4).click();
});

Then("I should see the change role modal", () => {
    cy.get('[data-testid="keyCreated"]').should("exist");
});

When("I click the close button", () => {
    cy.get('[data-testid="closeButton"]').click()
});

Then("I should not see the change role modal", () => {
    cy.get('[data-testid="keyCreated"]').should("not.exist")
});

Then("I should see the team member and current role populated in the form when I want to change a role", () => {
    cy.get('td>div>button').eq(0).click();
    cy.get('div>ul>li').eq(4).click();
    cy.get('[data-testid="teamMemberName"]').should('be.not.empty')
    cy.get('[data-testid="teamMemberRole"]').should('be.not.empty')
});

When("I click save without changing the member's role", () => {
    cy.get('[data-cy="saveNewRole"]').click();
});

Then("I should see a status error message", () => {
    cy.findByRole('status').should('have.text', 'organization must have at least one owner')
});

When("I select a new role and click save", () => {
    cy.get('td>div>button').eq(1).click();
    cy.get('div>ul>li').eq(4).click();
    cy.get('#role').click()
    cy.get('#react-select-7-option-1').click();
    cy.get('[data-cy="saveNewRole"]').click();
})

Then("the member's role should update in the member table", () => {
    cy.get('tr>td').eq(9).should('have.text', 'Member')
})

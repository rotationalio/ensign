import { Given, And, Then, When } from "cypress-cucumber-preprocessor/steps";

Given("I'm logged into Beacon", () => {
    cy.loginWith( { email: 'test1@test.com', password:'Abc123Def$56'})
});

When("I navigate to the team page", () => {
    cy.visit("/app/team");
});

Then("I should see the team table", () => {
    cy.get('[data-testid="teamTable"]').should("exist");
});

And("I should see the team data", () => {
    cy.get('tr>td').eq(0).should('have.text', 'Name')
    cy.get('[data-testid="cardlistitem-0"]').should('exist')
    cy.get('tr>td').eq(2).should('have.text', 'Email Address')
    cy.get('[data-testid="cardlistitem-1"]').should('exist')
    cy.get('tr>td').eq(4).should('have.text', 'Role')
    cy.get('[data-testid="cardlistitem-2"]').should('exist')
    cy.get('tr>td').eq(6).should('have.text', 'Status')
    cy.get('[data-testid="cardlistitem-3"]').should('exist')
    cy.get('tr>td').eq(8).should('have.text', 'Last Activity')
    cy.get('[data-testid="cardlistitem-4"]').should('exist')
    cy.get('tr>td').eq(10).should('have.text', 'Joined Date')
    cy.get('[data-testid="cardlistitem-5"]').should('exist')
    cy.get('tr>td').eq(12).should('have.text', 'Actions')
    cy.get('[data-testid="cardlistitem-6"]').should('exist')
});

Then("I should see the add team member button", () => {
    cy.get('[data-cy="add-team-member"]').should("exist");
});

When("I click the add team member button", () => {
    cy.get('[data-cy="add-team-member"]').click();
});

Then("I should see the invite new team member modal", () => {
    cy.get('[data-testid="memberCreationModal"]').should("exist");
});

When("I click the close button", () => {
    cy.get('[data-testid="closeButton"]').click()
});

Then("I should not see the invite new team member modal", () => {
    cy.get('[data-testid="memberCreationModal"]').should("not.exist");
});

When("I click invite without entering an email address", () => {
    cy.get('[data-cy="inviteMemberButton"]').click();
});

Then("I should see an error message", () => {
    cy.get('input>div').eq(0).should('have.text', 'Email is required')
});

When("I enter an invvalid email address", () => {
    cy.get('[data-cy="inviteMemberButton"]').type("testing123@test");
});

Then("I should see an error message", () => {
    cy.get('input>div').eq(0).should('have.text', 'Email is invalid.')

});

When("I enter an valid email address", () => {
    cy.get('[data-cy="inviteMemberButton"]').type("testing123@test.com");
});

And("I click invite the invite button", () => {
    cy.get('[data-cy="inviteMemberButton"]').click();
});

Then("I should see a new team member appear in the teams table with a pending status", () => {
    cy.get('tr>td').eq(4).should('have.text', 'Status')
    cy.get('[data-testid="cardlistitem-2"]').should('have.text', 'Pending')
});

When("I click the actions icon", () => {
    cy.get('[data-testid="cardlistitem-5"]').click();
});

Then("I should see the actions menu", () => {
    cy.findByRole('menuItem').should('exist')
});

When("I click change role", () => {
    cy.get('ul>li').eq(0).click();
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

Then("I should the team member and current role populated in the form when I want to change a role", () => {
    cy.get('[data-testid="teamMemberName"]').should('be.not.empty')
    cy.get('[data-testid="teamMemberRole"]').should('be.not.empty')
});

When("I click save without changing the member's role", () => {
    cy.get('[data-cy="saveNewRole"]').click();
});

Then("I should see an error message", () => {
    cy.findByRole('status').should('have.text', 'team member already has the requested role')
});

When("I select a new role and click save", () => {
    cy.get('select').select('Member').should('have.value', 'Member');
    cy.get('[data-cy="saveNewRole"]').click();
})

Then("the member's role should update in the member table", () => {
    cy.get('[data-testid="cardlistitem-2"]').should('have.text', 'Member')
})

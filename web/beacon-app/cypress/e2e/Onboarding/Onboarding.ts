import { And, Given, Then, When } from 'cypress-cucumber-preprocessor/steps';

beforeEach(function () {
    cy.fixture('user').then((user) => {
        this.user = user;
    });
})

Given('I\'m on the login page', () => {
    cy.visit('/');
});

When('I log into Beacon', function () {
    cy.loginWith({ email: this.user.email, password: this.user.password });
});

Then('I should be directed to the onboarding form', () => {
    cy.location('pathname').should('eq', '/app/onboarding');
});

And('I should see the onboarding sidebar', () => {
    cy.get('[data-cy="onboarding-sidebar"]').should('exist');
});

And('I should see my email address', function () {
    cy.get('[data-cy="user-email"]')
      .should('exist')
      .and('have.text', this.user.email);
});

When('I click log out in the topbar', () => {
    cy.get('[data-cy="log-out-bttn"]').click();
});

Then('I should be directed to the login page', () => {
    cy.location('pathname').should('eq', '/');
});

When('I log in a second time', function () {
    cy.loginWith({ email: this.user.email, password: this.user.password });
});

Then('I should be directed back to the onboarding form', () => {
    cy.location('pathname').should('eq', '/app/onboarding');
});

And('I should see the first step of the onboarding form', () => {
    cy.get('[data-cy="step-counter"]')
      .should('exist')
      .and('include.text', 'Step 1 of 4');
});

And('I should not see the Back button', () => {
    cy.get('[data-cy="back-bttn"]').should('not.exist');
});

And('I should see a default team name', () => {
    cy.get('[data-cy="team-name"]')
      .should('exist')
      .and('not.have.value', '');
});

When('I remove the default team name', () => {
    cy.get('[data-cy="team-name"]').clear();
});

And('I click the next button without entering a team name', () => {
    cy.get('[data-cy="next-bttn"]').click();
});

Then('I should see that the team name is required', () => {
    cy.get('[data-cy="team-name-error"]')
      .should('exist')
      .and('have.text', 'Team or organization name is required.');
});

And('I should not be able to continue to the second step', () => {
    cy.get('[data-cy="step-counter"]')
      .should('exist')
      .and('not.include.text', 'Step 2 of 4');
});

When('I enter a team name and click next', function () {
    cy.get('[data-cy="team-name"]').click().type(this.user.onboarding.team_name);
    cy.get('[data-cy="next-bttn"]').click();
});

Then('I should be directed to the second step of the onboarding form', () => {
    cy.get('[data-cy="step-counter"]')
      .should('exist')
      .and('include.text', 'Step 2 of 4');
});

And('I should see the Back button', () => {
    cy.get('[data-cy="back-bttn"]').should('exist');
});

Then('I should see a default workspace URL value', () => {
    cy.get('[data-cy="step-counter"]')
      .should('exist')
      .and('include.text', 'Step 2 of 4');

    cy.get('[data-cy="workspace"]')
       .should('exist')
       .and('not.have.value', '');
});

When('I click the Back button', () => {
    cy.get('[data-cy="back-bttn"]').click();
});

Then('I should be directed to the first step of the onboarding form', () => {
    cy.get('[data-cy="step-counter"]')
      .should('exist')
      .and('include.text', 'Step 1 of 4');
});

And('I should see the team name that I entered', function () {
    cy.get('[data-cy="team-name"]')
      .should('exist')
      .and('have.value', this.user.onboarding.team_name);
});

When('I click next to return to the second step of the onboarding form', () => {
    cy.get('[data-cy="next-bttn"]').click();
});

And('I delete the default workspace URL', () => {
    cy.get('[data-cy="workspace"]').clear();
});

And('I click next without entering a workspace URL', () => {
    cy.get('[data-cy="next-bttn"]').click();
});

Then('I should see that the workspace URL is required', () => {
    cy.get('[data-cy="workspace-error"]')
      .should('exist')
      .and('have.text', 'Workspace name is required.');
});

When('I enter an invalid workspace URL and click next', function () {
    cy.get('[data-cy="workspace"]')
      .clear()
      .type(this.user.onboarding.workspace.invalid_one);

    cy.get('[data-cy="next-bttn"]').click();
});

Then('I should see a validation error message', function () {
    cy.get('[data-cy="workspace-error"]')
      .should('have.text', this.user.onboarding.workspace.invalid_error);
});

When('I enter another invalid workspace URL and click next', function () {
    cy.get('[data-cy="workspace"]')
      .clear()
      .type(this.user.onboarding.workspace.invalid_two);

      cy.get('[data-cy="next-bttn"]').click();
});

Then('I should see another validation error message', function () {
    cy.get('[data-cy="workspace-error"]')
      .should('have.text', this.user.onboarding.workspace.invalid_error);
});

When('I enter an existing workspace URL and click next', function () {
    cy.get('[data-cy="workspace"]')
      .clear()
      .type(this.user.onboarding.workspace.exist);

    cy.get('[data-cy="next-bttn"]').click();
});

Then('I should get an error message that the workspace URL is already taken', function () {
    cy.findByRole('status')
      .should('exist')
      .and('have.text', this.user.onboarding.workspace.exist_error);
});

And('I should not be able to continue to the third step', () => {
    cy.get('[data-cy="next-bttn"]').click();
    cy.get('[data-cy="step-counter"]')
      .should('exist')
      .and('not.include.text', 'Step 3 of 4');
});

When('I enter a valid workspace URL', function () {
    cy.get('[data-cy="workspace"]')
      .clear()
      .type(this.user.onboarding.workspace.valid);
});

And('I click next to continue to the third step', () => {
    cy.get('[data-cy="next-bttn"]').click();
});

Then('I should be directed to the third step of the onboarding form', () => {
    cy.get('[data-cy="step-counter"]')
      .should('exist')
      .and('include.text', 'Step 3 of 4');
});

When('I click the Back button on the third step of the onboarding form', () => {
    cy.get('[data-cy="back-bttn"]').click();
});

Then('I should be directed to the second step of the onboarding form', () => {
    cy.get('[data-cy="step-counter"]')
      .should('exist')
      .and('include.text', 'Step 2 of 4');
});

And('I should see the workspace URL I entered', function () {
    cy.get('[data-cy="workspace"]')
      .should('exist')
      .and('have.value', this.user.onboarding.workspace.valid);
});

When('I click to return to the third step of the onboarding form', () => {
    cy.get('[data-cy="next-bttn"]').click();
});

And('I click next without entering a name', () => {
    cy.get('[data-cy="step-counter"]')
        .should('exist')
        .and('include.text', 'Step 3 of 4');
    
    cy.get('[data-cy="next-bttn"]').click();
});

Then('I should see that the name is required', () => {
    cy.get('[data-cy="user-name-error"]')
      .should('exist')
      .and('have.text', 'Name is required.');
});

When('I enter a name into the name input field and click next', function () {
    cy.get('[data-cy="user-name"]').type(this.user.onboarding.user_name);
    cy.get('[data-cy="next-bttn"]').click();
});

Then('I should be directed to the fourth step of the onboarding form', () => {
    cy.get('[data-cy="step-counter"]')
      .should('exist')
      .and('include.text', 'Step 4 of 4');
});

When('I click the Back button on the fourth step of the onboarding form', () => {
    cy.get('[data-cy="back-bttn"]').click();
});

Then('I should be directed back to the third step of the onboarding form', () => {
    cy.get('[data-cy="step-counter"]')
      .should('exist')
      .and('include.text', 'Step 3 of 4');
});

And('I should see the name I entered', function () {
    cy.get('[data-cy="user-name"]')
      .should('exist')
      .and('have.value', this.user.onboarding.user_name);
});

When('I click to return to the fourth step of the onboarding form', () => {
    cy.get('[data-cy="next-bttn"]').click();
});

Then('I should see the professional segment options', () => {
    cy.get('[data-cy="step-counter"]')
      .should('exist')
      .and('include.text', 'Step 4 of 4');

    cy.get('[data-cy="profession-segment"]').should('exist');
});

And('I should see the developer segment options', () => {
    cy.get('[data-cy="developer-segment"]').should('exist');
});

When('I click next before selecting a professional option or developer option', () => {
    cy.get('[data-cy="next-bttn"]').click();
});

Then('I should see that a professional segment option is required', () => {
    cy.get('[data-cy="profession-segment-error"]')
      .should('exist')
      .and('have.text', 'Please select one option.');
});

And('I should see that at least one developer segment option is required', () => {
    cy.get('[data-cy="developer-segment-error"]')
      .should('exist')
      .and('have.text', 'Please select at least one option.');
});

When('I select a professional option and not a developer option', () => {
    cy.get('[data-cy="profession-work"]').click({force: true});
});

And('I click the next button to continue', () => {
    cy.get('[data-cy="next-bttn"]').click();
});

Then('I should see that at least one developer option is required', () => {
    cy.get('[data-cy="developer-segment-error"]')
      .should('exist')
      .and('have.text', 'Please select at least one option.');
});

When('I select a first developer option', function () {
    cy.get('[id="developer_segment').click({multiple: true});

    cy.findByText(this.user.onboarding.dev_segment.option_one)
      .should('exist')
      .click();
    
});

And('I click a second developer option', function () {
    cy.get('[id="developer_segment').click({multiple: true});

    cy.findByText(this.user.onboarding.dev_segment.option_two)
      .should('exist')
      .click();
});

And('I click a third developer option', function () {
    cy.get('[id="developer_segment').click({multiple: true});

    cy.findByText(this.user.onboarding.dev_segment.option_three)
      .should('exist')
      .click();
});

Then('I should see that I cannot select any more developer options', function () {
    cy.get('[id="developer_segment').click({multiple: true});

    cy.findByText(this.user.onboarding.dev_segment.option_four)
    .should('exist')
    .and('have.attr', 'aria-disabled', 'true');
});

When('I click next to submit the onboarding form', () => {
    cy.get('[data-cy="developer-segment').click();
    cy.wait(1000);
    cy.get('[data-cy="next-bttn"]').click();
});

Then('I should be directed to the dashboard', () => {
    cy.location('pathname').should('eq', '/app');
});

And('I should see the onboarding sidebar has been replaced with the regular sidebar', () => {
    cy.get('[data-cy="onboarding-sidebar"]').should('not.exist');
    cy.get('[data-cy="sidebar"]').should('exist');
});

When('I click the log out button', () => {
    cy.get('[data-cy="menu"]').click({force: true})
    cy.findByText('Logout').click()
  });

Then('I should be directed to the login page', () => {
    cy.location('pathname').should('eq', '/');
});

When('I log into Beacon again', function () {
    cy.loginWith({ email: this.user.email, password: this.user.password });
});

Then('I should be directed to the dashboard and not see the onboarding workflow', () => {
    cy.location('pathname').should('eq', '/app');
});
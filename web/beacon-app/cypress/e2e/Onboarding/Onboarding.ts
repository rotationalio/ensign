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

And('I should see the first step of the onboarding form', () => {
    cy.get('[data-cy="step-counter"]')
      .should('exist')
      .and('include.text', 'Step 1 of 4');
});

And('I should not see the Back button', () => {
    cy.get('[data-cy="back-bttn"]').should('not.exist');
});

When('I remove the default team name', () => {
    cy.get('[data-cy="organization-name"]').clear();
});

And('I click the next button without entering a team name', () => {
    cy.get('[data-cy="next-bttn"]').click();
});

Then('I should see that the team name is required', () => {
    cy.get('[data-cy="organization-name-error"]')
      .should('exist')
      .and('have.text', 'Team or organization name is required.');
});

And('I should not be able to continue to the second step', () => {
    cy.get('[data-cy="step-counter"]')
      .should('exist')
      .and('not.include.text', 'Step 2 of 4');
});

When('I enter a team name and click next', function () {
    cy.get('[data-cy="organization-name"]').click().type(this.user.onboarding.org_name);
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

    cy.get('[data-cy="workspace-url"]')
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

And('I should see the organization name that I entered', function () {
    cy.get('[data-cy="organization-name"]')
      .should('exist')
      .and('have.value', this.user.onboarding.org_name);
});

When('I click next to return to the second step of the onboarding form', () => {
    cy.get('[data-cy="next-bttn"]').click();
});

When('I delete the default workspace URL', () => {
    cy.get('[data-cy="workspace-url"]').clear();
});

And('I click next without entering a workspace URL', () => {
    cy.get('[data-cy="next-bttn"]').click();
});

Then('I should see that the workspace URL is required', () => {
    cy.get('[data-cy="workspace-url-error"]')
      .should('exist')
      .and('have.text', 'Workspace name is required.');
});

And('I should not be able to continue to the third step', () => {
    cy.get('[data-cy="next-bttn"]').click();
    cy.get('[data-cy="step-counter"]')
      .should('exist')
      .and('not.include.text', 'Step 3 of 4');
});

When('I enter a workspace URL', function () {
    cy.get('[data-cy="workspace-url"]').type(this.user.onboarding.workspace_url);
});

And('I click next to continue to the third step', () => {
    cy.wait(5000);
    cy.get('[data-cy="next-bttn"]').click();
});

Then('I should be directed to the third step of the onboarding form', () => {
    cy.get('[data-cy="step-counter"]')
      .should('exist')
      .and('include.text', 'Step 3 of 4');
});

When(' I click the Back button on the third step of the onboarding form', () => {
    cy.get('[data-cy="back-bttn"]').click();
});

Then('I should be directed to the second step of the onboarding form', () => {
    cy.get('[data-cy="step-counter"]')
      .should('exist')
      .and('include.text', 'Step 2 of 4');
});

And('I should see the workspace URL I entered', function () {
    cy.get('[data-cy="workspace-url"]')
      .should('exist')
      .and('have.value', this.user.onboarding.workspace_url);
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

Then('I should be directed to the third step of the onboarding form', () => {
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

When('I select a first developer option', () => {
    cy.get('[id="developer_segment').should('exist').click({multiple: true});
    cy.get('[id="react-select-7-listbox').should('exist')
    cy.get('[id="react-select-7-option-0').should('exist').click();
});

Then('I select a second developer option', () => {
    cy.get('[id="developer_segment').should('exist').click({multiple: true});
    cy.get('[id="react-select-7-option-1').should('exist').click();
});

And('I select a third developer option', () => {
    cy.get('[id="developer_segment').should('exist').click({multiple: true});
    cy.get('[id="react-select-5-listbox').should('exist')
    cy.get('[id="react-select-5-option-2').should('exist').click();
});

Then('I should see that I cannot select any more developer options', () => {
    cy.get('[id="developer_segment').should('exist').click({multiple: true});
    cy.get('[id="react-select-5-option-3').should('be.disabled');
});

When('I click next to submit the onboarding form', () => {
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
import { And, Given, Then, When } from 'cypress-cucumber-preprocessor/steps';

beforeEach(function () {
    cy.fixture('user').then((user) => {
        this.user = user;
    });
});

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

And('I should see the name of the team I have been invited to join', function () {
    cy.get('[data-cy="sidebar-team-name"]')
      .should('exist')
      .and('contain.text', this.user.onboarding.team_name);
});

And('I should see my email address', function () {
    cy.get('[data-cy="user-email"]')
      .should('exist')
      .and('have.text', this.user.email);
});

And('I should see the option to log out', () => {
    cy.get('[data-cy="log-out-bttn"]').should('exist');
});

And('I should see step 3 of the onboarding form', () => {
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

And('I should not be able to edit the workspace URL', function () {
    cy.get('[data-cy="workspace"]')
      .should('exist')
      .and('be.disabled');
});

/* When('I click the Back button on the second step of the onboarding form', () => {
    cy.get('[data-cy="back-bttn"]').click();
});

Then('I should be directed to the first step of the onboarding form', () => {
    cy.get('[data-cy="step-counter"]')
      .should('exist')
      .and('include.text', 'Step 1 of 4');
});

And('I should not be able to edit the team name', () => {
    cy.get('[data-cy="team-name"]')
      .should('exist')
      .and('be.disabled');
}); */

When('I return to the third step of the onboarding form', () => {
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

And('I click a second developer option', () => {
    cy.get('[id="developer_segment').click({multiple: true});

    cy.findByText('Data engineering').should('exist').click();
});

And('I click a third developer option', function () {
    cy.get('[id="developer_segment').click({multiple: true});

    cy.findByText(this.user.onboarding.dev_segment.option_two)
      .should('exist')
      .click();
});

Then('I should see that I cannot select any more developer options', function () {
    cy.get('[id="developer_segment').click({multiple: true});
    cy.findByText(this.user.onboarding.dev_segment.option_three)
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

And('I should see the name of the organization I joined', function() {
    cy.get('[data-cy="org-name"]')
      .should('exist')
      .and('have.text', this.user.onboarding.team_name);
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
import { Given, Then, When, And } from 'cypress-cucumber-preprocessor/steps';

beforeEach(function () {
  cy.fixture('user').then((user) => {
      this.user = user;
  });
})

Given('I am on the Beacon homepage', function () {
  cy.loginWith( { email: this.user.email, password: this.user.password })
});

When('I\'m logged in', () => {
  cy.url().should('include', 'app')
});

Then('I should see the org name', function () {
  cy.get('[data-cy="org-name"]').should('exist').and('have.text', this.user.org_name)
});

And('I should see an avatar in the sidebar', () => {
  cy.get('[data-cy="avatar"]').should('be.visible')
});

When('I click on the avatar', () => {
  cy.get('[data-cy="avatar"]').click()
});

Then('I should see a list of orgs I belong to', () => {
  cy.get('[data-cy="org-menu"]').should('be.visible')
});


When('I click Projects in the sidebar', () => {
  cy.contains('span', 'Projects').should('be.visible').click()
});

Then('I should be taken to the Projects page', () => {
  cy.location('pathname').should('eq', '/app/projects')
});

When('I click Team in the sidebar', () => {
  cy.go('back')
  cy.contains('span', 'Team').should('be.visible').click()
});

Then('I should be taken to the Team page', () => {
  cy.location('pathname').should('eq', '/app/team')
});

When('I click Profile in the sidebar', () => {
  cy.go('back')
  cy.contains('span', 'Profile').should('be.visible').click()
});

Then('I should be taken to the Profile page', () => {
  cy.location('pathname').should('eq', '/app/profile')
});

When('I return to the home page', () => {
  cy.go('back')
});

Then('I should see external links in the sidebar', () => {
  cy.contains('a', 'Ensign U')
  .should('have.text', 'Ensign U ')
  .and('have.attr', 'href')
  .and('eq', 'https://rotational.io/blog/')

  cy.contains('a', 'Use Cases')
  .should('have.text', 'Use Cases ')
  .and('have.attr', 'href')
  .and('eq', 'https://ensign.rotational.dev/eventing/use_cases/')

  cy.contains('a', 'Docs')
  .should('have.text', 'Docs ')
  .and('have.attr', 'href')
  .and('eq', 'https://ensign.rotational.dev/getting-started/')

  cy.contains('a', 'Data Playground')
  .should('have.text', 'Data Playground ')
  .and('have.attr', 'href')
  .and('eq', 'https://rotational.io/data-playground/')

  cy.contains('a', 'SDKs')
  .should('have.text', 'SDKs ')
  .and('have.attr', 'href')
  .and('eq', 'https://ensign.rotational.dev/sdk/')

  cy.contains('a', 'Support')
  .should('have.text', 'Support')
  .and('have.attr', 'href')
  .and('eq', 'mailto:support@rotational.io')

  cy.get('ul li:first').should('have.text', 'About ')
  cy.get('ul>li>a').eq(0).should('have.attr', 'href').and('eq', 'https://rotational.io/about')

  cy.get('ul>li').eq(1).should('have.text', 'Contact Us ')
  cy.get('ul>li>a').eq(1).should('have.attr', 'href').and('eq', 'https://rotational.io/contact')

  cy.get('ul>li').eq(2).should('have.text', 'Server Status ')
  cy.get('ul>li>a').eq(2).should('have.attr', 'href').and('eq', 'https://status.rotational.dev/')
});

When('I see the Welcome component', () => {
  cy.get('[data-cy="ensign-welcome"]').should('be.visible')
});

Then('I should see the Welcome to Ensign video', () => {
  cy.get('[data-cy="welcome-video"]').should('be.visible')
});

When('I click on the welcome video', () => {
  cy.get('[data-cy="welcome-video-bttn"]').should('be.visible').click()
});

Then('I should see a modal open with a playable version of the video', () => {
  cy.get('.modal-video').should('be.visible')
});

When('I click the close button to close the modal', () => {
  cy.get('.modal-video-close-btn').should('be.visible').click()
});

Then('I should not see the modal with the video', () => {
  cy.get('.modal-video').should('not.be.visible')
});

When('I  see the Set Up A New Project component', () => {
  cy.get('[data-cy="setup-new-project"]').should('be.visible')
});

And('I click the Create Project button', () => {
  cy.get('[data-cy="create-project-bttn"]').should('be.visible').click()
});

Then('I should see the create project modal', () => {
    cy.get('[data-testid="newProjectModal"]').should('be.visible')
});

When('I see the settings button', () => {
  cy.get('[data-testid="menu"]').click()
  cy.findByText('Settings').should('be.visible')
});

And('I click the settings button', () => {
  cy.findByText('Settings').click()
});

Then('I should visit the settings page', () => {
  cy.location('pathname').should('eq', '/app/profile')
});

When('I return to the main page', () => {
  cy.go('back')
});

And('I click the logout button', () => {
  cy.get('[data-testid="menu"]').click({force: true})
  cy.findByText('Logout').click()
});

Then('I should be logged out of the Beacon home page', () => {
  cy.location('pathname').should('eq', '/')
});
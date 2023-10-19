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

And('I should see an avatar', () => {
  cy.get('[data-cy="avatar"]').should('exist')
});

When('I click on the avatar', () => {
  cy.get('[data-cy="avatar"]').click()
});

Then('I should see a list of orgs I belong to', function () {
  cy.findByRole('menuitem').eq(0).should('have.text', this.user.org_two)
  cy.findByRole('menuitem').eq(1).should('have.text', this.user.org_three)
});

And('I should see Projects in the sidebar', () => {
  cy.contains('div', 'Projects').should('be.visible')
});

And('I should see Team in the sidebar', () => {
  cy.contains('div', 'Team').should('be.visible')
});

And('I should see a link to Profile in the sidebar', () => {
  cy.contains('a', 'Profile').should('have.text', 'Profile ')
});

And('I should see Ensign U in the sidebar', () => {
  cy.contains('a', 'Ensign U').should('have.text', 'Ensign U ')
});

And('I should see a link to Docs in the sidebar', () => {
  cy.contains('a', 'Docs').should('have.text', 'Docs ')
});

And('I should see a link to Support in the sidebar', () => {
  cy.contains('a', 'Support').should('have.text', 'Support')
});

And('I should see a link to the About page in the sidebar footer', () => {
    cy.get('ul li:first').should('have.text', 'About ')
  });

Then('I should be able to visit the About page if I click the link', () => {
    cy.get('ul>li>a').eq(0).should('have.attr', 'href').and('eq', 'https://rotational.io/about')
  });

And('I should see a link to the Contact Us page in the sidebar footer', () => {
    cy.get('ul>li').eq(1).should('have.text', 'Contact Us ')
  });

Then('I should be able to visit the Contact Us page if I click the link', () => {
    cy.get('ul>li>a').eq(1).should('have.attr', 'href').and('eq', 'https://rotational.io/contact')
  });

And('I should see a link to the Server Status page in the sidebar footer', () => {
    cy.get('ul>li').eq(2).should('have.text', 'Server Status ')
  });

  Then('I should be able to visit the Server Status page if I click the link', () => {
    cy.get('ul>li>a').eq(2).should('have.attr', 'href').and('eq', 'https://status.rotational.dev/')
  });

And('I should see the Welcome component', () => {
  cy.get('[data-cy="projWelcome"]').should('be.visible')
})

And('I should see the Set Up A New Project component', () => {
  cy.get('[data-testid="cardlistitem"]').eq(0).should('be.visible')
});

  When('I click on the Create Project button', () => {
    cy.get('[data-testid="set-new-project"]').should('be.visible').click()
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
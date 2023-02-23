import { Given, Then, When, And } from 'cypress-cucumber-preprocessor/steps';

Given('I am on the Beacon homepage', () => {
  cy.loginWith( { email: 'test@ensign.com', password:'Abc123Def$56'})
});

When('I\'m logged in', () => {
  cy.url().should('include', 'app')
    cy.getCookies().should('exist')
    cy.getCookie('bc_atk').should('exist')
});

Then('I should see an avatar', () => {
  cy.get('[data-testid="avatar"]').should('exist')
  cy.get('[data-testid="avatar"]').invoke('attr', 'src').should('contain', 'https://www.gravatar.com/')
});

// TODO: Add test to ensure org name matches org name included in registration
And('I should see an org name', () => {
  cy.get('[data-testid="orgName"]').should('exist')
});

And('I should see a link to Docs in the sidebar', () => {
    cy.contains('a', 'Docs').should('have.text', 'Docs ')
  });

And('I should see a link to Profile in the sidebar', () => {
    cy.contains('a', 'Profile').should('have.text', 'Profile ')
  });

  And('I should see a link to the About page in the sidebar footer', () => {
    cy.get('ul li:first').should('have.text', 'About')
  });

  Then('I should be able to visit the About page if I click the link', () => {
    cy.get('ul li:first').should('have.html', '<a href="https://rotational.io/about">About</a>')
  });

  And('I should see a link to the Contact Us page in the sidebar footer', () => {
    cy.get('ul>li').eq(1).should('have.text', 'Contact Us')
  });

  Then('I should be able to visit the Contact Us page if I click the link', () => {
    cy.get('ul>li').eq(1).should('have.html', '<a href="https://rotational.io/contact">Contact Us</a>')
  });

  // TODO: Test server link when it is added
  And('I should see a link to the Server Status page in the sidebar footer', () => {
    cy.get('ul>li').eq(2).should('have.text', 'Server Status')
  });

And('I should see quickview data', () => {
  const data = [
    'Active Projects',
    'Topics',
    'API Keys',
    'Data Storage'
  ]
    cy.get('h5').each(($e, index) => {
    cy.wrap($e).should('have.text', data[index])
    cy.wrap($e).siblings().should('not.be.empty')
  })
});

When('I see the Manage project button', () => {
  cy.get('[data-testid="manage"]').should('have.text', "Manage Project")
});

Then('I should not be able to click it', () => {
  cy.get('[data-testid="manage"]').should('have.css', 'cursor', 'not-allowed')
});

And('I should see the Create API Key button', () => {
  cy.get('[data-testid="key"]').should('have.text', 'Create API Key')
  });

  // TODO: Add API Key tests

And('I should see the View Docs button', () => {
  cy.get('[data-testid="viewDocs"]').should('have.text', 'View Docs')
  });

Then('I should be able to visit the external documentation site', () => {
    cy.get('[data-testid="viewDocsLink"]').should('have.attr', 'href').and('eq', 'https://ensign.rotational.dev/getting-started/')
  });

And('I should see the settings button', () => {
  cy.get('[data-testid="menu"]').click()
  cy.get('[data-testid="settings"]').should('exist')
  });

When('I click the settings button', () => {
  cy.get('[data-testid="settings"]').click()
});

Then('I should visit the settings page', () => {
  cy.location('pathname').should('eq', '/app/organization')
});

And('I should go back to the main page', () => {
  cy.go('back')
});

Then('I should see the log out button', () => {
  cy.get('[data-testid="menu"]').click()
  cy.get('[data-testid="logoutButton"]').should('exist')
  });

When('I click the logout button', () => {
  cy.get('[data-testid="logoutButton"]').click()
});

Then('I should log out of the Beacon home page', () => {
  cy.location('pathname').should('eq', '/')
});
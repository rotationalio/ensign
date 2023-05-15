import { Given, Then, When, And } from 'cypress-cucumber-preprocessor/steps';

Given('I am on the Beacon homepage', () => {
  cy.loginWith( { email: 'test7@test.com', password:'Abc123Def$56'})
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

And('I should see an org name', () => {
  cy.get('[data-testid="orgName"]').should('exist').should('have.text', 'GoldTeam')
});

And('I should see Projects in the sidebar', () => {
  cy.contains('div', 'Projects').should('be.visible')
})

And('I should see Team in the sidebar', () => {
  cy.contains('div', 'Team').should('be.visible')
})

And('I should see a link to Docs in the sidebar', () => {
    cy.contains('a', 'Docs').should('have.text', 'Docs ')
  });

  And('I should see a link to Support in the sidebar', () => {
    cy.contains('a', 'Support').should('have.text', 'Support')
  });


And('I should see a link to Profile in the sidebar', () => {
    cy.contains('a', 'Profile').should('have.text', 'Profile ')
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

And('I should see quickview data', () => {
  const data = [
    'Projects',
    'Topics',
    'API Keys',
    'Data Storage'
  ]
    cy.get('h5').each(($e, index) => {
    cy.wrap($e).should('have.text', data[index])
    cy.wrap($e).siblings().should('not.be.empty')
  })
});

And('I should see the Welcome component', () => {
  cy.get('[data-cy="projWelcome"]').should('be.visible')
})

When('I click the Start button', () => {
  cy.get('[data-cy="startSetupBttn"]').should('be.visible').click()
})

Then('I should see the create project modal', () => {
  cy.get('[data-testid="newProjectModal"]').should('be.visible')
});

When('I click the close button', () => {
cy.get('[data-testid="newProjectModal"]').within(() => {
  cy.get('button>svg').click()
})
});

Then('I should not see the create project modal', () => {
  cy.get('[data-testid="newProjectModal"]').should('not.exist')
});

When('I click on the Start button again', () => {
  cy.get('[data-cy="startSetupBttn"]').should('be.visible').click()
});

And('I fill in the project name', () => {
  cy.get('[data-cy="project-name"]').should('exist').type('My first project')
});

And('I fill in the project description', () => {
  cy.get('[data-cy="project-description"]').should('exist').type('A new project to test Ensign!')
});

And('I click the Create Project button', () => {
  cy.get('[data-cy="NewProjectButton"]').should('exist').click()
});

Then('I should be redirected to the projects page', () => {
cy.location('pathname').should('include', 'app/projects')
});

And('I should see the new project in the project list table', () => {
cy.get('[data-cy="projectTable"]').should('exist').within(() => {
  cy.get('tbody>tr>td').eq(0).should('have.text', 'My first project')
      cy.get('tbody>tr>td').eq(1).should('have.text', 'A new project to test Ensign!')
      cy.get('tbody>tr>td').eq(2).should('have.text', 'Incomplete')
      cy.get('tbody>tr>td').eq(3).should('have.text', '05/12/2023')
  })
});

Then('I should see the number of projects increase to 1 on the projects page', () => {
cy.get('h5').eq(0).should('have.text', 'Projects').siblings('p').should('have.text', 1)
});

When('I go back to the main page', () => {
cy.go('back')
});

And('I should see the number of projects increase to 1 on the main page', () => {
cy.get('h5').eq(0).should('have.text', 'Projects').siblings('p').should('have.text', 1)
});

And('I should see the Set Up A New Project component', () => {
  cy.get('[data-testid="cardlistitem"]').eq(0).should('be.visible')
});

And('I should see the Access Resources component', () => {
  cy.get('[data-testid="cardlistitem"]').eq(1).should('be.visible')
});

And('I should see the Access button', () => {
  cy.get('[data-testid="viewDocs"]').should('have.text', 'Access')
});

Then('I should be able to visit the external documentation site', () => {
    cy.get('[data-testid="viewDocsLink"]').should('have.attr', 'href').and('eq', 'https://ensign.rotational.dev/getting-started/')
});

  When('I click on the Create button', () => {
    cy.get('[data-testid="set-new-project"]').should('be.visible').click()
});

Then('I should see the create project modal', () => {
    cy.get('[data-testid="newProjectModal"]').should('be.visible')
});

When('I fill in the project name', () => {
    cy.get('[data-cy="project-name"]').should('exist').type('My second project')
});

And('I fill in the project description', () => {
    cy.get('[data-cy="project-description"]').should('exist').type('One more project!')
});

And('I click the Create Project button', () => {
    cy.get('[data-cy="NewProjectButton"]').should('exist').click()
});

Then('I should be redirected to the projects page', () => {
  cy.location('pathname').should('include', 'app/projects')
});

And('I should see the new project in the project list table', () => {
  cy.get('[data-cy="projectTable"]').should('be.visible').within(() => {
    cy.get('tbody>tr>td').eq(5).should('have.text', 'My second project')
        cy.get('tbody>tr>td').eq(6).should('have.text', 'One more project!')
        cy.get('tbody>tr>td').eq(7).should('have.text', 'Incomplete')
        cy.get('tbody>tr>td').eq(8).should('have.text', '05/12/2023')   
  })
});

And('I should see the number of projects increase to 2 on the projects page', () => {
  cy.get('h5').eq(0).should('have.text', 'Projects').siblings('p').should('have.text', 2)
});

When('I go back to the main page', () => {
  cy.go('back')
});

Then('I should see the number of projects increase increase to 2 on the main page', () => {
  cy.get('h5').eq(0).should('have.text', 'Projects').siblings('p').should('have.text', 2)
});

And('I should not see the Welcome component', () => {
  cy.get('[data-cy="projWelcome"]').should('not.exist')
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

And('I go back to the main page', () => {
  cy.go('back')
});

And('I click the logout button', () => {
  cy.get('[data-testid="menu"]').click({force: true})
  cy.findByText('Logout').click()
});

Then('I should log out of the Beacon home page', () => {
  cy.location('pathname').should('eq', '/')
});
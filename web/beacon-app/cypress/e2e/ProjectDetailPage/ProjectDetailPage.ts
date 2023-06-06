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

And("I click on a project in the project table list", () => {
    cy.get('tbody>tr').eq(7).click()
});

Then("I should see the project detail page for the project", () => {
    cy.location('pathname').should('include', 'app/projects/01H06EBW4ZQSKTB5PV79CWBQJX')
});

And("I should see the project name at the top of the page", () => {
    cy.get('[data-cy="project-name"]').should('have.text', 'Project Test')
});

/* And("I should see the project's details when I hover over the hint icon next to the project name", () => {
    cy.get('[data-cy="detailHint"]').children('svg').realHover()
    cy.get('[data-cy="prjDetail"]').should('exist')
}); */

When("I see the project setup component", () => {
    cy.get('[data-testid="project-setup"]').should('be.visible')
});

Then("I should see that a project has been created", () => {
    cy.get('[data-testid="project-created"]').should('be.visible')
});

When("I click on the cogwheel", () => {
    cy.get('[data-cy="detailActions"]').click()
});

Then("I should see the project detail actions menu", () => {
    cy.contains('Delete Project')
    cy.contains('Edit Project')
    cy.contains('Change Owner')
});

When("I click Delete Project", () => {
    cy.get('[data-testid="cancelButton"]').click()
});

Then("I should see the delete project modal", () => {
    cy.get('[data-testid="delete-prj-modal"]')
});

And("I should not see the delete project modal when I click the close button", () => {
    cy.get('[data-testid="delete-prj-modal"]').within(() => {
        cy.get('button>svg').click()
    })
    cy.get('[data-testid="delete-prj-modal"]').should('not.exist')
});

When("I click Edit Project", () => {
    cy.get('[data-cy="detailActions"]').click()
    cy.get('[data-testid="rename-project"]').click()
});

Then("I should see the Edit Project modal", () => {
    cy.get('[data-cy="edit-project"]').should('be.visible')
});

And("I should not see the Edit Project modal when I click the close button", () => {
    cy.get('[data-cy="edit-project"]').within(() => {
        cy.get('button>svg').click()
    })
    cy.get('[data-cy="edit-project"]').should('not.exist')
});

When("I re-open the Edit Project modal", () => {
    cy.get('[data-cy="detailActions"]').click()
    cy.get('[data-testid="rename-project"]').click()
});

Then("I should see the current project's name", () => {
    cy.get('[data-cy="current-proj-name"]').should('have.value', 'Project Test')
});

When("I enter a new project name", () => {
    cy.get('[data-cy="new-proj-name"]').type('One more project')
});

And("I change the project's description", () => {
    cy.get('[data-cy="project-description"]').clear().type('Making a new project!')
})

And("I click save", () => {
    cy.get('[data-cy="edit-proj-bttn"]').click()
})

Then("I should see the new project name", () => {
    cy.get('[data-cy="project-name"]').should('have.text', 'One more project')
    cy.go('back')
    cy.get('tbody>tr>td').eq(49).should('have.text', 'One more project')
});

And("I should see the updated project description", () => {
    cy.get('tbody>tr>td').eq(50).should('have.text', 'Making a new project!')
})

When("I click Change Owner", () => {
    cy.go('forward')
    cy.get('[data-cy="detailActions"]').click()
    cy.get('[data-cy="change-owner"]').click()
})

Then("I should see the Change Owner modal", () => {
    cy.get('[data-cy="change-proj-owner"]').should('be.visible')
})

And("I should not see the Change Owner modal when I click the close button", () => {
    cy.get('[data-cy="change-proj-owner"]').within(() => {
        cy.get('button>svg').click()
    })
    cy.get('[data-cy="change-proj-owner"]').should('not.exist')
});

When("I change the project's owner", () => {
    cy.get('[data-cy="detailActions"]').click()
    cy.get('[data-cy="change-owner"]').click()
    cy.get('[data-cy="change-proj-owner"]').should('be.visible')
    cy.get('.css-1xc3v61-indicatorContainer').click()
    cy.get('#react-select-5-option-0').click()
})

And("I click the Save button", () => {
    cy.get('[data-cy="update-owner"]').click()
})

Then("I should see the new project owner", () => {
    cy.go('back')
    cy.get('tbody>tr>td').eq(54).should('have.text', 'Kamala Khan')
})

When("I see the API Keys component", () => {
    cy.go('forward')
    cy.get('[data-cy="keyComp"]').should('exist')
});

/* Then("I should see more details about API keys when I hover over the hint icon", () => {
    cy.get('[data-cy="keyHint"]').trigger('focus', {force: true} )
    cy.wait(5000)
    cy.get('[data-cy="keyInfo"]').should('exist')
});
 */
And("I should see the API key list table", () => {
    cy.get('[data-cy="keyTable"]').should('be.visible')
});

When("I click the + New Key button", () => {
    cy.get('[data-cy="addKey"]').should('be.visible').click()
});

Then("I should see the Generate API Key modal", () => {
    cy.get('[data-testid="keyModal"]').should('be.visible')
});

When("I click the close button I should not see the Generate API Key modal", () => {
    cy.get('[data-testid="keyModal"]').within(() => {
        cy.get('button>svg').click()
    })
    cy.get('[data-testid="keyModal"]').should('not.exist')
});

When('I re-open the Generate API Key modal', () => {
    cy.get('[data-cy="addKey"]').should('be.visible').click()
    cy.get('[data-testid="keyModal"]').should('be.visible')
});

Then("I should see an error if I try to create an API key without entering a Key Name", () => {
    cy.get('[data-testid="generateKey"]').click()
    cy.get('small').should('contain', 'The key name is required.')
});

When("I enter a Key Name", () => {
    cy.get('[data-testid="keyName"]').type('Test Key')
});

And("I click the Generate API Key button", () => {
    cy.get('[data-testid="generateKey"]').click()
});

Then("I should see the Your API Key modal", () => {
    cy.get('[data-testid="keyCreated"]').should('be.visible')
});

And("I should be able to copy the Client ID", () => {
    cy.get('[data-testid="clientId"]').should('be.visible')
    cy.get('[data-testid="copyID"]').should('be.visible')
});

And("I should be able to copy the Client Secret", () => {
    cy.get('[data-testid="clientSecret"]').should('be.visible')
    cy.get('[data-testid="copySecret"]').should('be.visible')
});

And("I should be able to download the API key details", () => {
    cy.get('[data-testid="download"]').should('be.visible')
});

When("I confirm that have read the info on the Your API Key modal", () => {
    cy.get('[data-testid="closeKey"]').click()
});

Then("I should see the new API key in the API key list table", () => {
    cy.get('[data-cy="keyTable"]').within(() => {
        cy.get('tbody>tr>td').eq(0).should('have.text', 'Test Key')
        cy.get('tbody>tr>td').eq(1).should('have.text', 'Full')
        cy.get('tbody>tr>td').eq(2).should('have.text', 'Unused')
        cy.get('tbody>tr>td').eq(3).should('have.text', 'N/A')
        cy.get('tbody>tr>td').eq(4).should('have.text', '05/12/2023')
    })
});

And("I should see that an API key has been created", () => {
    cy.get('[data-testid="api-key-created"]').should('be.visible')
});

When("I see the Topics component", () => {
    cy.get('[data-cy="topicComp"]').should('exist')
});

/* Then("I should see more details about topics when I hover over the hint icon", () => {
    cy.get('[data-cy="topicHint"]').trigger('focus', {force: true} )
    cy.wait(5000)
    cy.get('[data-cy="topicInfo"]').should('exist')
}); */

And("I should see the topic list table", () => {
    cy.get('[data-cy="topicTable"]').should('be.visible')
});

When("I click the + New Topic button", () => {
    cy.get('[data-cy="addTopic"]').should('be.visible').click()
});

Then("I should see the New Topic modal", () => {
    cy.get('[data-testid="topicModal"]').should('be.visible')
});

When("I click the close button I should not see the New Topic modal", () => {
    cy.get('[data-testid="topicModal"]').within(() => {
        cy.get('button>svg').click()
    })
    cy.get('[data-testid="topicModal"]').should('not.exist')
});

When('I re-open the New Topic modal', () => {
    cy.get('[data-cy="addTopic"]').should('be.visible').click()
    cy.get('[data-testid="topicModal"]').should('be.visible')
});

Then("I should see an error if I try to create a topic without entering a Topic Name", () => {
    cy.get('[data-cy="createTopic"]').click()
    cy.get('input').siblings('div').should('contain', 'Topic name is required.')
});

And("I should see an error if I type an invalid Topic Name", () => {
    cy.get('[data-cy="topicName"]').type('Test Topic')
    cy.get('[data-cy="createTopic"]').click()
    cy.get('input').siblings('div').should('contain', 'Topic name cannot include spaces.')
    cy.get('[data-cy="topicName"]').clear().type('0TestTopic')
    cy.get('[data-cy="createTopic"]').click()
    cy.get('input').siblings('div').should('contain', 'Topic name may only start with a letter and contain letters, numbers, underscores, and dashes.')
});

When("I enter a valid Topic Name", () => {
    cy.get('[data-cy="topicName"]').clear().type('Test-Topic-01')
});

And("I click the Create Topic button", () => {
    cy.get('[data-cy="createTopic"]').click()
});

Then("I should see the new topic in the topic list table", () => {
    cy.get('[data-cy="topicTable"]').within(() => {
        cy.get('tbody>tr>td').eq(0).should('have.text', 'Test-Topic-01')
        cy.get('tbody>tr>td>span').eq(0).should('have.text', 'active')
        cy.get('tbody>tr>td').eq(5).should('have.text', '05/12/2023')
    })
});

And("I should not see the project setup component", () => {
    cy.get('[data-testid="project-setup"]').should('not.exist')
});

When("I go back to the projects page", () => {
    cy.go('back')
    cy.location('pathname').should('include', 'app/projects')
});

Then("I should see the number of Topics increase by 1", () => {
    cy.contains('p', '7').should('exist')
});

And("I should see the number of API keys increase by 1", () => {
    cy.contains('p', '25').should('exist')
});

When("I go back to the main page", () => {
    cy.go('back')
    cy.location('pathname').should('include', 'app')
});

Then("I should see the number of Topics increase by 1", () => {
    cy.contains('p', '7').should('exist')
});

And("I should see the number of API keys increase by 1", () => {
    cy.contains('p', '25').should('exist')
});
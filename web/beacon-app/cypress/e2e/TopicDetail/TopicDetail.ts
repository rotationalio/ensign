import { And, Given, Then, When } from 'cypress-cucumber-preprocessor/steps';

Given("I'm logged into Beacon", () => {
    cy.fixture('user').then((user) => {
        cy.loginWith( {email: user.email, password: user.password });
    });
});

And("I navigate to the Topic Detail Page", () => {
    cy.contains('div', 'Projects').should('exist').click();
    cy.location('pathname').should('include', 'app/projects')
    cy.get('[data-cy="projectTable"]').within(() => {
        cy.get('tr>td').eq(0).click()
    })
    cy.fixture('user').then((user) => {
        cy.location('pathname').should('include', `app/projects/${user.projectID}`)
    });
    
    cy.get('[data-cy="topicTable"]').within(() => {
        cy.get('tr>td').eq(0).click()
    })
    
    cy.fixture('user').then((user) => {
        cy.location('pathname').should('include', `app/topics/${user.topicID}`)
    });
});

Then("I should see the topic name in the header component", () => {
    cy.get('[data-cy="topic-name"]').should('have.text', 'Test-Topic-01')
});

And("I should see the topic state tag", () => {
    cy.get('[data-cy="topic-status-tag"]')
    .should('exist')
    .and('have.text', 'Active');
});

And("I should see the cogwheel icon in the header component", () => {
    cy.get('[data-cy="topic-detail-actions"]').should('exist')
});

When("I click the cogwheel icon", () => {
    cy.get('[data-cy="topic-detail-actions"]').click()
});

Then("I should see a menu with menu items for Archive Topic, Delete Topic, and Clone Topic", () => {
    cy.contains('Archive Topic').should('exist')
    cy.contains('Delete Topic').should('exist')
    cy.contains('Clone Topic').should('exist')
});

When("I click Archive Topic", () => {
    cy.contains('Archive Topic').click()
});

Then("I should see the Archive Topic modal", () => {
    cy.get('[data-cy="archive-topic-modal"]').should('exist')
});

When("I click x to close the Archive Topic modal", () => {
    cy.get('[data-cy="archive-topic-modal"]').within(() => {
        cy.get('button').click()
    })
});

Then("I should not see the Archive Topic modal", () => {
    cy.get('[data-cy="archive-topic-modal"]').should('not.exist')
});

When("I click Delete Topic", () => {
    cy.get('[data-cy="topic-detail-actions"]').click()
    cy.contains('Delete Topic').click()
});

Then("I should see the Delete Topic modal", () => {
    cy.get('[data-cy="delete-topic-modal"]').should('exist')
});

When("I click x to close the Delete Topic modal", () => {
    cy.get('[data-cy="delete-topic-modal"]').within(() => {
        cy.get('button').click()
    })
});

Then("I should not see the Delete Topic modal", () => {
    cy.get('[data-cy="delete-topic-modal"]').should('not.exist')
});

When("I click Clone Topic", () => {
    cy.get('[data-cy="topic-detail-actions"]').click()
    cy.contains('Clone Topic').click()
});

Then("I should see the Clone Topic modal", () => {
    cy.get('[data-cy="clone-topic-modal"]').should('exist')
});

When("I click x to close the Clone Topic modal", () => {
    cy.get('[data-cy="clone-topic-modal"]').within(() => {
        cy.get('button').click()
    })
});

Then("I should not see the Clone Topic modal", () => {
    cy.get('[data-cy="clone-topic-modal"]').should('not.exist')
});

And("I should see 4 cards with metrics for the topic", () => {
    cy.get('[data-cy="quick-view-card-0"]').should('exist')
    cy.get('[data-cy="quick-view-card-1"]').should('exist')
    cy.get('[data-cy="quick-view-card-2"]').should('exist')
    cy.get('[data-cy="quick-view-card-3"]').should('exist')
});

When("I see the event detail table", () => {
    cy.get('[data-cy="event-detail-table"]').should('exist');
});

Then("I should see the event detail table headers", () => {
    cy.get('[data-cy="event-detail-table"]').within(() => {
        cy.get('th > div').eq(0).should('have.text', 'Event Type');
        cy.get('th > div').eq(1).should('have.text', 'Version');
        cy.get('th').eq(2).should('have.text', 'MIME Type');
        cy.get('th').eq(3).should('have.text', '# of Events');
        cy.get('th').eq(4).should('have.text', '% of Events');
        cy.get('th').eq(5).should('have.text', 'Storage Volume');
        cy.get('th').eq(6).should('have.text', '% of Volume');
    });
});



And("I should see the Topic Query compoent", () => {
    cy.get('[data-cy="topic-query-title"]').should('exist')
});

Then("I should see the Topic Query carat toggle is open by default and pointed down", () => {
    cy.get('[data-cy="topic-query-carat-down"]').should('exist')
})

And("I should see the Topic Query text instructions", () => {
    cy.get('[data-cy="topic-query-instructions"]').should('exist')
});

When("I click on the Topic Query title the carat toggle should be closed and pointed up", () => {
    cy.get('[data-cy="topic-query-heading"]').click()
    cy.get('[data-cy="topic-query-carat-up"]').should('exist')
    cy.get('[data-cy="topic-query-carat-down"]').should('not.exist')
});

And("I should not be able to see the Topic Query content", () => {
    cy.get('[data-cy="topic-query-instructions"]').should('not.exist')
});

When("I click on the Topic Query title again, the content should be visible", () => {
    cy.get('[data-cy="topic-query-heading"]').click()
    cy.get('[data-cy="topic-query-carat-down"]').should('exist')
    cy.get('[data-cy="topic-query-carat-up"]').should('not.exist')
    cy.get('[data-cy="topic-query-instructions"]').should('exist')
});

And("I should see the Advanced Topic Policy Management compoent", () => {
    cy.get('[data-cy="topic-mgmt-title"]').should('exist')
});

Then("I should see the Advanced Topic Policy Management carat toggle is open by default and pointed down", () => {
    cy.get('[data-cy="topic-mgmt-carat-down"]').should('exist')
    cy.get('[data-cy="topic-query-carat-up"]').should('not.exist')
});

And("I should see the Advanced Topic Policy Management content", () => {
    cy.get('[data-cy="topic-mgmt-content"]').should('exist')
});

When("I click on the Advanced Topic Policy Management title the carat toggle should be closed and pointed up", () => {
    cy.get('[data-cy="topic-mgmt-heading"]').click()
    cy.get('[data-cy="topic-mgmt-carat-up"]').should('exist')
    cy.get('[data-cy="topic-mgmt-carat-down"]').should('not.exist')
});

And("I should not see the Advanced Topic Policy Management content", () => {
    cy.get('[data-cy="topic-mgmt-content"]').should('not.exist')
});

When("I click on the Advanced Topic Policy Management title again, the content should be visible", () => {
    cy.get('[data-cy="topic-mgmt-heading"]').click()
    cy.get('[data-cy="topic-mgmt-carat-down"]').should('exist')
    cy.get('[data-cy="topic-mgmt-carat-up"]').should('not.exist')
    cy.get('[data-cy="topic-mgmt-content"]').should('exist')
});
import { And, Given, Then, When } from 'cypress-cucumber-preprocessor/steps';

Given("I'm logged into Beacon", () => {
    cy.fixture('user').then((testUser) => {
        cy.loginWith( {email: testUser.email, password: testUser.password });
    });
});

And("I navigate to the Topic Detail Page", () => {
    cy.contains('div', 'Projects').should('exist').click();
    cy.location('pathname').should('include', 'app/projects')
    cy.get('[data-cy="projectTable"]').within(() => {
        cy.get('tr>td').eq(0).click()
    })
    cy.location('pathname').should('include', 'app/projects/01H0924MJW1QB8JKT30PP815ET')
    cy.get('[data-cy="topicTable"]').within(() => {
        cy.get('tr>td').eq(0).click()
    })
    cy.location('pathname').should('include', 'app/topics/01H66SDYV5AXMDGH13EQC1BZ9Z')
});

Then("I should see the topic name in the header component", () => {
    cy.get('[data-cy="topic-name"]').should('have.text', 'Test-Topic-01')
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

And("I should see the Topic Query compoent", () => {
    cy.get('[data-cy="topic-query"]').should('exist')
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
    cy.get('[data-cy="topic-mgmt"]').should('exist')
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
import { And, Given, Then, When } from 'cypress-cucumber-preprocessor/steps';

Given("I'm logged into Beacon", () => {
    cy.fixture('user').then((testUser) => {
        cy.loginWith( {email: testUser.email, password: testUser.password });
    });
});

When("I see the topic query section of the Topic Detail page", () => {
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
    cy.get('[data-cy="topic-query-title"]').should('exist');
});

Then("I should see the EnSQL link in the topic query instructions", () => {
    cy.get('[data-cy="topic-query-enql-link"]')
    .should('exist')
    .should('have.text', 'EnSQL')
    .should('have.attr', 'href', 'https://ensign.rotational.dev/ensql/')
    .should('have.attr', 'target', '_blank');
});

And("I should see the SDKs link in the topic query instructions", () => {
    cy.get('[data-cy="topic-query-sdks-link"]')
    .should('exist')
    .should('have.text', 'SDKs')
    .should('have.attr', 'href', 'https://ensign.rotational.dev/sdk/')
    .should('have.attr', 'target', '_blank');
});

When("I see the topic query input field", () => {
    cy.get('[data-cy="topic-query-input"]').should('exist');
});

Then("I should see the default topic query", () => {
    cy.get('[data-cy="topic-query-input"]').should('have.value', 'SELECT * FROM Test-Topic-01 LIMIT 1');
});

And("I should see the query button", () => {
    cy.get('[data-cy="submit-topic-query-bttn"]').should('exist');
});

And("I should see the clear button", () => {
    cy.get('[data-cy="clear-topic-query-bttn"]').should('exist');
});

And("I should see no query results", () => {
    cy.get('[data-cy="query-result-count"]')
    .should('exist')
    .should('have.text', '0 results of 0 total');
});

And("I should see no viewing event results", () => {
    cy.get('[data-cy="viewing-event-result-count"]')
    .should('exist')
    .should('have.text', 'Viewing Event: 0 of 0');
});

And("I should see an empty meta data table", () => {
    cy.get('[data-cy="query-meta-data-table"]').should('exist');
    cy.get('[data-cy="query-meta-data-table"]').within(() => {
        cy.get('p').should('have.text', 'No data available');
    });
});

And("I should see NA listed for the mime type and event type", () => {
    cy.get('[data-cy="result-mime-type"]')
    .should('exist')
    .and('have.text', 'N/A');

  cy.get('[data-cy="result-event-type"]')
    .should('exist')
    .and('have.text', 'N/A');
});

When("I see the topic query results view", () => {
    cy.get('[data-cy="topic-query-result"]').should('exist')
});

Then("I should see the no query result message", () => {
    cy.get('[data-cy="topic-query-result"]')
    .and('have.text', 'No query result. Try the default query or enter your own query. See EnSQL documentation for example queries.');
});

And("I should see disabled pagination buttons", () => {
    cy.get('[data-cy="prev-query-bttn"]').should('be.disabled');
    cy.get('[data-cy="next-query-bttn"]').should('be.disabled');
});

When("I click the query button", () => {
    cy.get('[data-cy="submit-topic-query-bttn"]').click();
});

Then("I should see 10 query results out of 11 total results", () => {
    cy.get('[data-cy="query-result-count"]')
    .should('have.text', '10 results of 11 total');
});

And("I should view event 1 of 10", () => {
    cy.get('[data-cy="viewing-event-result-count"]')
    .should('have.text', 'Viewing Event: 1 of 10');
});

And("I should see the mime type for event 1", () => {
    cy.get('[data-cy="result-mime-type"]')
    .should('have.text', 'text/plain');
});

And("I should see the event type for event 1", () => {
    cy.get('[data-cy="result-event-type"]')
    .should('have.text', 'Message v1.0.0');
});

And("I should see the topic query result for event 1", () => {
    cy.get('[data-cy="topic-query-result"]')
    .should('have.text', 'hello world');
});

And("I should see that the previous button is disabled", () => {
    cy.get('[data-cy="prev-query-bttn"]').should('be.disabled');
});

And("I should see that the next button is enabled", () => {
    cy.get('[data-cy="next-query-bttn"]').should('be.enabled');
});

When("I click the next button", () => {
    cy.get('[data-cy="next-query-bttn"]').click();
});

Then("I should view event 2 of 10", () => {
    cy.get('[data-cy="viewing-event-result-count"]')
    .should('have.text', 'Viewing Event: 2 of 10');
});

And("I should see the mime type for event 2", () => {
    cy.get('[data-cy="result-mime-type"]')
    .should('have.text', 'text/csv');
});

And("I should see the event type for event 2", () => {
    cy.get('[data-cy="result-event-type"]')
    .should('have.text', 'Spreadsheet v1.1.0');
});

And("I should see that the previous button is enabled", () => {
    cy.get('[data-cy="prev-query-bttn"]').should('be.enabled');
});

When("I click to view a result that could not be parsed", () => {
    cy.get('[data-cy="next-query-bttn"]')
    .click()
    .click()
    .click()
    .click()
    .click()
    .click();
});

Then("I should see the could not parse message next to the mime type", () => {
    cy.get('[data-cy="result-mime-type"]')
    .should('have.text', 'Could not parse. Rendered as base-64 encoded data.');
});

And("I should see base 64 encoded data in the results view", () => {
    cy.get('[data-cy="topic-query-result"]')
    .should('have.text', 'gaRuYW1lo0JvYqNhZ2UYHg==');
});

When("I click to the last result", () => {
    cy.get('[data-cy="next-query-bttn"]')
    .click()
    .click();

    cy.get('[data-cy="viewing-event-result-count"]')
    .should('have.text', 'Viewing Event: 10 of 10');
});

Then("I should see that the next button is disabled", () => {
    cy.get('[data-cy="next-query-bttn"]').should('be.disabled');
});

And("I should see that the previous button is still enabled", () => {
    cy.get('[data-cy="prev-query-bttn"]').should('be.enabled');
});

When("I click the previous button", () => {
    cy.get('[data-cy="prev-query-bttn"]').click();
});

Then("I should view event 9 of 10", () => {
    cy.get('[data-cy="viewing-event-result-count"]')
    .should('have.text', 'Viewing Event: 9 of 10');
});

When("I click the clear button", () => {
    cy.get('[data-cy="clear-topic-query-bttn"]').click();
});

Then("I should see the default result view with the no query result message", () => {
    cy.get('[data-cy="topic-query-result"]')
    .and('have.text', 'No query result. Try the default query or enter your own query. See EnSQL documentation for example queries.');
});

And("I should not see a value in the input field", () => {
    cy.get('[data-cy="topic-query-input"]').should('not.have.value');
});

And("I should see that the pagination buttons are disabled", () => {
    cy.get('[data-cy="prev-query-bttn"]').should('be.disabled');
    cy.get('[data-cy="next-query-bttn"]').should('be.disabled');
});

When("I click query without typing a query", () => {
    cy.get('[data-cy="submit-topic-query-bttn"]').click();
});

Then("I should see the validation error message", () => {
    cy.get('[data-cy="topic-query-form"]').within(() => {
        cy.get('div > div > div')
        .should('have.text', 'Please enter a query. See EnSQL documentation for examples of valid queries.');
    });
});

When("I type a query into the input field", () => {
    cy.get('[data-cy="topic-query-input"]')
    .type('SELECT * FROM Test-Topic-01 LIMIT 5');
});

Then("I should not see the validation error message", () => {
    cy.get('[data-cy="topic-query-form"]').within(() => {
        cy.get('div > div > div')
        .should('not.exist');
    });
});
import { And, Given, Then, When } from 'cypress-cucumber-preprocessor/steps';

beforeEach(function () {
    cy.fixture('user').then((user) => {
        this.user = user;
    });
})
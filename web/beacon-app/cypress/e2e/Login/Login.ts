import { Given, Then, When, And } from 'cypress-cucumber-preprocessor/steps';

// Given('I fill my credentials and submit', () => {

//    cy.loginWith({ email: "masskoder+toto@gmail.com", password: "Tototo2022@!" });
   
// });

Given('I open the login page', () => {
    cy.visit('/');
});

When('I fill my credentials and submit the login form', () => {

    cy.loginWith({ email: "masskoder+toto@gmail.com", password: "Tototo2022@!" });

});

// And('I Submit the Login Form', () => {
//     cy.get('[data-testid="login-button"]').click();
// });


Then('I\'m Logged In', () => {
    cy.url().should('include', 'app')
    cy.getCookies().should('exist')
    cy.getCookie('bc_atk').should('exist')
    
});

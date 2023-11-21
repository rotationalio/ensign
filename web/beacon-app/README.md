# Beacon App

The official user UI for Ensign.

This project was bootstrapped with [ViteJS](https://vitejs.dev/).

## Getting Started

Change into the `beacon-app` directory of the repository.

In order to install or update dependencies:

```
$ yarn install
```

To view the app in development mode:

```
$ yarn dev
```

Once the local server is running, you can navigate to [http://localhost:3000/](http://localhost:3000/) to view the site in the browser. This will allow you to view pages that do not require a user to be logged in. 

In order to register an account or view pages that may only be viewed after logging in, open a separate terminal window and change into the root directory.

Build the docker compose images:

```
$ ./containers/local.sh -p backend build
```

After the images are built run the docker containers:

```
$ ./containers/local.sh -p backend up
```

## Registering & Verifying A Development Account

With the local server running, you can navigate to [http://localhost:3000/register](http://localhost:3000/register) to create an account to be used in development mode.

After completing the registration process, you will need to verify your email address. To complete the verification process, open the email located in the following directory:

`containers/quarterdeck/emails/your-email-address`

You should see a .mim file with the registration email asking you to confirm your email address. 

After clicking on verification link or copying and pasting it into a browser, your email address should be verified. 

You should then be able to log into Beacon. If you are creating an account for the first time, you will need to complete the onboarding process before you are navigated to the dashboard.

## Linting Errors

After making changes, the page should automatically reload. If you're running the server locally, linting errors should appear in the terminal window and the browser.

If you see a linting error message that is `potentially fixable with the --fix option`, you may need to open a new terminal. 

In the new terminal, change into the `beacon-app` directory and run the following command to resolve the issue:

```
$ yarn lint
```

## Beacon Design System

The Beacon app is composed of components created in the Beacon Design System. For more details, visit the [Beacon Design System repository](https://github.com/rotationalio/beacon-ds).


## Tailwind CSS

The [Tailwind CSS](https://tailwindcss.com/) library is used to add CSS styles to the site.


## Building the App

To build the app for production to the `build` folder:

```
$ yarn build
```

This will bundle React in production mode and optimize the build for the best performance.

import browserify from '@cypress/browserify-preprocessor';
import cucumber from 'cypress-cucumber-preprocessor';

// eslint-disable-next-line
type EventListener = (eventName: string | symbol, listener: (...args: any[]) => void) => void;

module.exports = (on: EventListener) => {
  const options = browserify.defaultOptions;

  options.browserifyOptions.plugin.unshift(['tsify', { project: '../../tsconfig.json' }]);

  on('file:preprocessor', cucumber(options));
};

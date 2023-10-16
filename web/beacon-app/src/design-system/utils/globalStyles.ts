import { createGlobalStyle } from 'styled-components';

const global = {
  '*': {
    'box-sizing': 'border-box',
  },
  '*:before': {
    'box-sizing': 'border-box',
  },
  '*:after': {
    'box-sizing': 'border-box',
  },
};

const GlobalStyles = createGlobalStyle({
  ...require('@rotational/beacon-foundation/lib/tokens/css/tokens.css'),
  ...global,
});

export default GlobalStyles;

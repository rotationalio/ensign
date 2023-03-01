import { MenuItemUnstyled, menuItemUnstyledClasses, PopperUnstyled } from '@mui/base';
import styled from 'styled-components';

const grey = {
  50: '#f6f8fa',
  100: '#eaeef2',
  200: '#d0d7de',
  300: '#afb8c1',
  400: '#8c959f',
  500: '#6e7781',
  600: '#57606a',
  700: '#424a53',
  800: '#32383f',
  900: '#24292f',
};

export const StyledListbox = styled('ul')(
  () => `
  font-family: Montserrat, Inter, -apple-system, sans-serif;
  font-size: 0.875rem;
  box-sizing: border-box;
  margin: 12px 0;
  padding: 10px 0px;
  min-width: 200px;
  border-radius: 12px;
  overflow: auto;
  outline: 0px;
  background: ${grey[900]} ;
  border: 1px solid ${grey[200]};
  color: #fff;
  box-shadow: 0px 4px 30px ${grey[200]};
  `
);

export const StyledMenuItem = styled(MenuItemUnstyled)(
  () => `
  list-style: none;
  padding: 8px 14px;
  cursor: default;
  user-select: none;
  outline: none;

  &:last-of-type {
    border-bottom: none;
  }

  &.${menuItemUnstyledClasses.disabled} {
    color: ${grey[400]};
  }

  &:hover:not(.${menuItemUnstyledClasses.disabled}) {
    background-color: ${grey[100]};
    color: ${grey[900]};
  }
  `
);

export const Popper = styled(PopperUnstyled)`
  z-index: 1;
`;

export const MenuSectionRoot = styled('li')`
  list-style: none;

  & > ul {
    padding-left: 1em;
  }
`;

export const MenuSectionLabel = styled('span')`
  display: block;
  padding: 10px 0 5px 10px;
  font-size: 0.75em;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05rem;
  color: ${grey[600]};
`;

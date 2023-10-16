import styled from 'styled-components';

export const StyledLabel = styled.label(() => {
  return {
    '&[aria-required=true]::after': {
      content: "'*'",
      marginLeft: 'var(--spacings-1)',
      color: 'var(--colors-danger-500)',
    },
  };
});

export default StyledLabel;

import { AriaButton } from '@rotational/beacon-core';
import React from 'react';
import styled from 'styled-components';

type ButtonProps = React.ComponentProps<typeof AriaButton>;

function Button(props: ButtonProps) {
  return <StyledButton {...props} />;
}

const StyledButton = styled(AriaButton)((props) => ({
  fontSize: '0.875rem',
  lineHeight: '1.25rem',
  '&:focus': {
    outline: 'none',
  },
  ...(props.isDisabled && {
    background: 'rgb(233 236 239)',
    color: 'rgb(206 212 218)',
  }),
  ...(props.variant === 'primary' && {
    '&:hover': {
      background: 'rgba(29,101,166, 0.8)!important',
      borderColor: 'rgba(29,101,166, 0.8)!important',
    },
    '&:active': {
      background: 'rgba(29,101,166, 0.8)!important',
      borderColor: 'rgba(29,101,166, 0.8)!important',
    },
  }),
  ...(props.variant === 'secondary' && {
    '&:hover': {
      background: 'rgba(230,104,9, 0.8)!important',
      borderColor: 'rgba(230,104,9, 0.8)!important',
    },
    '&:active': {
      background: 'rgba(230,104,9, 0.8)!important',
      borderColor: 'rgba(230,104,9, 0.8)!important',
    },
  }),
  '&:disabled': {
    background: 'rgb(233 236 239)',
    color: 'rgb(206 212 218)',
    pointer: 'not-allowed',
  },
}));

export default Button;

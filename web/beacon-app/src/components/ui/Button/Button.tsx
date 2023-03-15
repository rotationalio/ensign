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
  ...(props.isDisabled && {
    background: 'rgb(233 236 239)',
    color: 'rgb(206 212 218)',
  }),
  '&:hover': {
    ...(!props.isDisabled && {
      background: 'rgba(29,101,166, 0.8)!important',
    }),
    ...(props.variant === 'secondary' &&
      !props.isDisabled && {
        background: 'rgba(230,104,9, 0.8)!important',
      }),
  },
}));

export default Button;

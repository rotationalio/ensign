import { Button } from '@rotational/beacon-core';
import styled from 'styled-components';

const StyledButton = styled(Button)((props) => ({
  fontSize: '0.875rem',
  lineHeight: '1.25rem',
  '&:focus': {
    outline: 'none',
  },
  ...(props.disabled && {
    background: 'rgb(233 236 239)',
    color: 'rgb(206 212 218)',
    pointer: 'not-allowed',
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
  ...(props.variant === 'ghost' && {
    backgroundColor: 'white!important',
    color: 'rgba(52 58 64)!important',
    border: 'none!important',
    height: 'auto!important',
    width: 'auto!important',
    '&:hover': {
      background: 'rgba(255,255,255, 0.8)!important',
      borderColor: 'rgba(255,255,255, 0.8)!important',
    },
    '&:active': {
      background: 'rgba(255,255,255, 0.8)!important',
      borderColor: 'rgba(255,255,255, 0.8)!important',
    },
  }),
  '&:disabled': {
    background: 'rgb(233 236 239)',
    color: 'rgb(206 212 218)',
    pointer: 'not-allowed',
  },
}));

export default StyledButton;

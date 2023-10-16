import ModalUnstyled from '@mui/base/ModalUnstyled';
import clsx from 'clsx';
import React from 'react';
import styled from 'styled-components';

import CloseIcon from './CloseIcon';
import { ModalContainerProps } from './Modal.types';

export const BackdropUnstyled = React.forwardRef<
  HTMLDivElement,
  { open?: boolean; className: string }
>((props, ref) => {
  const { open, className, ...other } = props;
  return <div className={clsx({ 'MuiBackdrop-open': open }, className)} ref={ref} {...other} />;
});

export const StyledModal = styled(ModalUnstyled)`
  position: fixed;
  z-index: 1300;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
`;

export const Backdrop = styled(BackdropUnstyled)`
  z-index: -1;
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.5);
  -webkit-tap-highlight-color: transparent;
`;

export const Container = styled.div<ModalContainerProps>(({ fullScreen, size }) => ({
  minWidth: '20rem',
  borderRadius: '1rem',
  background: '#fff',
  padding: '1rem 2rem',
  outline: 'none',
  position: 'relative',
  overflow: 'scroll',

  ...(size === 'small' && {
    width: '25%',
  }),

  ...(size === 'medium' && {
    width: '50%',
  }),

  ...(size === 'large' && {
    width: '80%',
    padding: 20,
  }),

  ...(fullScreen && {
    width: '100%',
    height: '100%',
    borderRadius: 0,
  }),
}));

export const Title = styled.h1((props) => ({
  fontSize: '1.25rem' /* 20px */,
  fontWeight: 600,
  fontFamily: 'Inter, Roboto, sans-serif',
  lineHeight: '1.75rem' /* 28px */,
  textAlign: 'center',
  margin: '1rem 0',
}));

export const StyledCloseButton = styled.button((props) => ({
  position: 'absolute',
  top: 15,
  right: 20,
  color: 'gray',
  borderRadius: '50%',
  border: '2px solid gray',
  padding: 1,
  ':hover': {},
}));

export const CloseButton = () => {
  return (
    <StyledCloseButton>
      <CloseIcon />
    </StyledCloseButton>
  );
};

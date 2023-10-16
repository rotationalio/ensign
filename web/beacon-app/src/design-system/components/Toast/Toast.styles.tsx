import styled from 'styled-components';

import { ToastWithRadixProps } from './Toast.types';

const StyledToast = styled.div<Pick<ToastWithRadixProps, 'variant'>>((props) => ({
  ...(props.variant === 'default' && {
    backgroundColor: 'var(--colors-gray-100)',
    color: 'var(--colors-gray-900)',
  }),
  ...(props.variant === 'primary' && {
    backgroundColor: 'var(--colors-primary)',
    color: 'var(--colors-white)',
  }),
  ...(props.variant === 'secondary' && {
    backgroundColor: 'var(--colors-secondary)',
    color: 'var(--colors-white)',
  }),
  ...(props.variant === 'success' && {
    backgroundColor: 'var(--colors-green)',
    color: 'var(--colors-white)',
  }),
  ...(props.variant === 'danger' && {
    backgroundColor: 'var(--colors-red)',
    color: 'var(--colors-white)',
  }),
  ...(props.variant === 'warning' && {
    backgroundColor: 'var(--colors-yellow)',
    color: 'var(--colors-gray-900)',
  }),
  ...(props.variant === 'info' && {
    backgroundColor: 'var(--colors-blue)',
    color: 'var(--colors-white)',
  }),
  borderRadius: 'var(--radius-1)',
  padding: 'var(--space-2)',
  display: 'grid',
  gridTemplateAreas: '"title action" "description action"',
  gridTemplateColumns: 'auto max-content',
  columnGap: 'var(--spacings-5)',
  alignItems: 'center',
  boxShadow: 'hsl(206 22% 7% / 35%) 0px 10px 38px -10px, hsl(206 22% 7% / 20%) 0px 10px 20px -15px',
  ':data-state="open"': {
    animation: 'toast-in 0.3s ease-out',
  },
  ':data-state="closed"': {
    animation: 'toast-out 0.3s ease-out',
  },
  '@keyframes toast-in': {
    '0%': {
      opacity: 0,
      transform: 'translateY(10px)',
    },
    '100%': {
      opacity: 1,
      transform: 'translateY(0)',
    },
  },
  '@keyframes toast-out': {
    '0%': {
      opacity: 1,
      transform: 'translateY(0)',
    },
    '100%': {
      opacity: 0,
      transform: 'translateY(10px)',
    },
  },
}));

StyledToast.defaultProps = {
  variant: 'default',
};

const StyledToastTitle = styled.div({
  fontWeight: 'bold',
  fontSize: 'var(--fontSizes-2)',
  marginBottom: 'var(--space-1)',
});

const StyledToastDescription = styled.div({
  fontSize: 'var(--fontSizes-1)',
});

const StyledToastCloseButton = styled.button({
  position: 'absolute',
  top: 'var(--space-1)',
  right: 'var(--space-1)',
  border: 'none',
  backgroundColor: 'transparent',
  color: 'var(--colors-gray-900)',
  cursor: 'pointer',
  '&:hover': {
    color: 'var(--colors-gray-800)',
  },
});

export { StyledToast, StyledToastCloseButton, StyledToastDescription, StyledToastTitle };

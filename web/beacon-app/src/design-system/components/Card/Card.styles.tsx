import styled from 'styled-components';

import Box from '../Box';
export const CardWrapper = styled(Box)({
  display: 'flex',
  borderRadius: 'var(--border-radius-1)',
  borderColor: '#0000004d',
  borderStyle: 'solid',
  overflow: 'hidden',
});

export const CardHeader = styled(Box)({
  display: 'flex',
  alignItems: 'center',
  minHeight: 'var(--spacings-8)',
  paddingBottom: 'var(--spacings-2)',
  fontWeight: 500,
});

CardHeader.displayName = 'CardHeader';

export const CardBody = styled(Box)({
  overflow: 'hidden',
  position: 'relative',
});

CardBody.displayName = 'CardBody';

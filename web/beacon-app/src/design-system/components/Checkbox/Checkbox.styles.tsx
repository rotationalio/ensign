import styled from 'styled-components';

export const Label = styled.label((props) => ({
  display: 'flex',
  alignItems: 'center',
}));

export const Input = styled.input((props) => ({
  marginRight: 8,
}));

export const Span = styled.span((props: { isDisabled?: boolean }) => ({
  ...(props?.isDisabled && {
    color: 'gray',
  }),
}));

import styled from 'styled-components';

type LoaderProps = {
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl';
  color?: string;
  className?: string;
  [key: string]: any;
};
const Loader = (props: LoaderProps) => (
  <StyledSpinner viewBox="0 0 50 50" {...props}>
    <circle className="path" cx="25" cy="25" r="20" fill="none" strokeWidth="2" />
  </StyledSpinner>
);
const StyledSpinner = styled.svg<LoaderProps>((props) => ({
  animation: 'rotate 2s linear infinite',
  ...(props.size === 'xs' && {
    width: '1.5rem',
    height: '1.5rem',
  }),
  ...(props.size === 'sm' && {
    width: '2.5rem',
    height: '2.5rem',
  }),
  ...(props.size === 'md' && {
    width: '3.5rem',
    height: '3.5rem',
  }),
  ...(props.size === 'lg' && {
    width: '4.5rem',
    height: '4.5rem',
  }),
  ...(props.size === 'xl' && {
    width: '6.5rem',
    height: '6.5rem',
  }),

  '& .path': {
    stroke: props.color || 'currentColor',
    strokeLinecap: 'round',
    animation: 'dash 1.5s ease-in-out infinite',
  },
  '@keyframes rotate': {
    '100%': {
      transform: 'rotate(360deg)',
    },
  },
  '@keyframes dash': {
    '0%': {
      strokeDasharray: '1, 150',
      strokeDashoffset: '0',
    },
    '50%': {
      strokeDasharray: '90, 150',
      strokeDashoffset: '-35',
    },
    '100%': {
      strokeDasharray: '90, 150',
      strokeDashoffset: '-124',
    },
  },
}));

StyledSpinner.defaultProps = {
  size: 'md',
};

export default Loader;

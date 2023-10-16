import { cva } from 'class-variance-authority';

export const getTextFieldInputStyle = cva(['px-4 py-2 font-sm rounded-md text-sm'], {
  variants: {
    size: {
      xsmall: ['text-xs', 'py-1', 'px-2'],
      small: ['text-sm', 'py-1', 'px-2'],
      medium: ['text-md', 'py-3', 'px-4'],
      large: ['text-base', 'py-4', 'px-7'],
    },
  },
  defaultVariants: {
    size: 'medium',
  },
});

export const setSizeVariantStyle = (size = 'medium') => {
  switch (size) {
    case 'xsmall':
      return 'text-xs py-1 px-2';
    case 'small':
      return 'text-sm py-1 px-2';
    case 'medium':
      return 'text-md py-3 px-4';
    case 'large':
      return 'text-base py-4 px-7';
    default:
      return 'text-md py-3 px-4';
  }
};

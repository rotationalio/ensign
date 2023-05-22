import { mergeClassnames } from '@rotational/beacon-core';
import styled from 'styled-components';
const StyledTextArea = styled.textarea({
  outline: 'none',
  '&:focus': {
    borderColor: 'var(--colors-primary-900)',
  },
  '&[aria-invalid=true]': {
    borderColor: 'var(--colors-danger-500)',
  },
});

export type TextAreaProps = {
  className?: string;
  placeholder?: any;
  label?: string;
  labelClassName?: string;
  errorMessage?: any;
  fullWidth?: boolean;
  disabled?: boolean;
  rows?: number;
  cols?: number;
  value?: string;
  onChange?: (event: React.ChangeEvent<HTMLTextAreaElement>) => void;
  onBlur?: (event: React.FocusEvent<HTMLTextAreaElement>) => void;
  onFocus?: (event: React.FocusEvent<HTMLTextAreaElement>) => void;
  'data-cy'?: string;
  [key: string]: any;
};

function TextArea({
  className,
  placeholder,
  label,
  labelClassName,
  errorMessage,
  fullWidth,
  disabled,
  rows,
  cols,
  value,
  onChange,
  onBlur,
  onFocus,
  'data-cy': dataCy,
  ...rest
}: TextAreaProps) {
  return (
    <div className={mergeClassnames('flex flex-col', fullWidth && 'w-full')}>
      {label && (
        <label htmlFor={dataCy} className={mergeClassnames(' mb-2', labelClassName)}>
          {label}
        </label>
      )}
      <StyledTextArea
        id={dataCy}
        className={mergeClassnames(
          'border-none p-2',
          errorMessage && 'border-danger-500',
          className
        )}
        placeholder={placeholder}
        aria-label={label}
        aria-invalid={!!errorMessage}
        disabled={disabled}
        rows={rows}
        cols={cols}
        value={value}
        onChange={onChange}
        onBlur={onBlur}
        onFocus={onFocus}
        data-cy={dataCy}
        {...rest}
      />
      {errorMessage && (
        <div className="text-sm text-danger-500" role="alert">
          {errorMessage}
        </div>
      )}
    </div>
  );
}

export default TextArea;

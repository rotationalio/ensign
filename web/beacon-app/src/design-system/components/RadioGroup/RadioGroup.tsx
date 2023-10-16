import { useRadioGroup } from 'react-aria';
import { useRadioGroupState } from 'react-stately';

import { RadioContext } from './RadioGroup.context';
import { RadioGroupProps } from './RadioGroup.type';

export default function RadioGroup(props: RadioGroupProps) {
  const { children, label, description, errorMessage, validationState } = props;
  const state = useRadioGroupState(props);
  const { radioGroupProps, labelProps, descriptionProps, errorMessageProps } = useRadioGroup(
    props,
    state
  );

  return (
    <div {...radioGroupProps}>
      <span {...labelProps}>{label}</span>
      <RadioContext.Provider value={state}>{children}</RadioContext.Provider>
      {description && (
        <div {...descriptionProps} style={{ fontSize: 12 }}>
          {description}
        </div>
      )}
      {errorMessage && validationState === 'invalid' && (
        <div {...errorMessageProps} style={{ color: 'red', fontSize: 12 }}>
          {errorMessage}
        </div>
      )}
    </div>
  );
}

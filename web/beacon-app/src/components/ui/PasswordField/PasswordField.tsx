import { TextField } from '@rotational/beacon-core';

import { PasswordFieldProps } from './PasswordField.type';
import { usePasswordField } from './usePasswordField';

const PasswordField = (props: PasswordFieldProps) => {
  const _props = usePasswordField(props);

  return <TextField {..._props} />;
};

export default PasswordField;

// this class will be used to validate the form that are repeated in the application
import { t } from '@lingui/macro';
import * as Yup from 'yup';

class FormValidation {
  static passwordValidation = Yup.string()
    .required(t`The password is required.`)
    .matches(/^(?=.*[a-z])/, t`The password must contain at least one lowercase letter.`)
    .matches(/^(?=.*[A-Z])/, t`The password must contain at least one uppercase letter.`)
    .matches(/^(?=.*[0-9])/, t`The password must contain at least one number.`)
    .matches(
      /^(?=.*[!/[@#$%^&*+,-./:;<=>?^_`{|}~])/,
      t`The password must contain at least one special character.`
    )
    .matches(/^(?=.{12,})/, t`The password must be at least 12 characters long.`);

  static ConfirmPassword = (ref: string) =>
    Yup.string()
      .oneOf([Yup.ref(`${ref}`), null], t`The passwords must match.`)
      .required(t`Please re-enter your password to confirm.`);
}

export default FormValidation;

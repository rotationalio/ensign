import { t } from '@lingui/macro';
import * as Yup from 'yup';

import { NewUserAccount } from '../../../types/RegisterService';

const validationSchema = Yup.object().shape({
  email: Yup.string()
    .email(t`Email is invalid.`)
    .required(t`The email address is required.`),
  password: Yup.string()
    .required(t`The password is required.`)
    .matches(/^(?=.*[a-z])/, t`The password must contain at least one lowercase letter.`)
    .matches(/^(?=.*[A-Z])/, t`The password must contain at least one uppercase letter.`)
    .matches(/^(?=.*[0-9])/, t`The password must contain at least one number.`)
    .matches(
      /^(?=.*[!/[@#$%^&*+,-./:;<=>?^_`{|}~])/,
      t`The password must contain at least one special character.`
    )
    .matches(/^(?=.{12,})/, t`The password must be at least 12 characters long.`),

  pwcheck: Yup.string()
    .oneOf([Yup.ref('password'), null], t`The passwords must match.`)
    .required(t`Please re-enter your password to confirm.`),
  invite_token: Yup.string().notRequired(),
  privacy_agreement: Yup.boolean().notRequired(),
  terms_agreement: Yup.boolean().notRequired(),
}) satisfies Yup.SchemaOf<NewUserAccount>;

export default validationSchema;

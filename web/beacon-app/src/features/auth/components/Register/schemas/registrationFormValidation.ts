import * as Yup from 'yup';

import { NewUserAccount } from '../../../types/RegisterService';

const validationSchema = Yup.object().shape({
  name: Yup.string().required('The name is required.'),
  email: Yup.string().email('Email is invalid.').required('The email address is required.'),
  password: Yup.string()
    .required('The password is required.')
    .matches(/^(?=.*[a-z])/, 'The password must contain at least one lowercase letter.')
    .matches(/^(?=.*[A-Z])/, 'The password must contain at least one uppercase letter.')
    .matches(/^(?=.*[0-9])/, 'The password must contain at least one number.')
    .matches(
      /^(?=.*[!/[@#$%^&*+,-./:;<=>?^_`{|}~])/,
      'The password must contain at least one special character.'
    )
    .matches(/^(?=.{12,})/, 'The password must be at least 12 characters long.'),

  pwcheck: Yup.string()
    .oneOf([Yup.ref('password'), null], 'The passwords must match.')
    .required('Please re-enter your password to confirm.'),
  organization: Yup.string().required('The organization is required.'),
  domain: Yup.string().required('The domain is required.'),
  terms_agreement: Yup.boolean().required('The agreement is required.'),
  invite_token: Yup.string().notRequired(),
}) satisfies Yup.SchemaOf<Omit<NewUserAccount, 'privacy_agreement'>>;

export default validationSchema;

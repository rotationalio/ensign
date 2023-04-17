import * as Yup from 'yup';

import { NewUserAccount } from '../../../types/RegisterService';

const validationSchema = Yup.object().shape({
  name: Yup.string().required('The name is required.'),
  email: Yup.string().email().required('The email address is required.'),
  password: Yup.string().required('The password is required.'),
  pwcheck: Yup.string()
    .oneOf([Yup.ref('password'), null], 'The passwords must match.')
    .required('Please re-enter your password to confirm.'),
  organization: Yup.string().required('The organization is required.'),
  domain: Yup.string().required('The domain is required.'),
  terms_agreement: Yup.boolean().required('The agreement is required.'),
  invite_token: Yup.string().notRequired(),
}) satisfies Yup.SchemaOf<Omit<NewUserAccount, 'privacy_agreement'>>;

export default validationSchema;

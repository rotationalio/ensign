import * as Yup from 'yup';

import { NewInvitedUserAccount } from '@/features/auth';

const validationSchema = Yup.object().shape({
  name: Yup.string().required('The name is required.'),
  email: Yup.string().email().required('The email address is required.'),
  password: Yup.string().required('The password is required.'),
  pwcheck: Yup.string()
    .oneOf([Yup.ref('password'), null], 'The passwords must match.')
    .required('The confirm password is required.'),
  terms_agreement: Yup.boolean().required('The agreement is required.'),
}) satisfies Yup.SchemaOf<Omit<NewInvitedUserAccount, 'privacy_agreement'>>;

export default validationSchema;

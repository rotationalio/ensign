import * as Yup from 'yup';

import { NewInvitedUserAccount } from '@/features/auth';

const validationSchema = Yup.object().shape({
  email: Yup.string().email().required('The email address is required.'),
  password: Yup.string().required('The password is required.'),
  pwcheck: Yup.string()
    .oneOf([Yup.ref('password'), null], 'The passwords must match.')
    .required('Please re-enter your password to confirm.'),
  invite_token: Yup.string().notRequired(),
}) satisfies Yup.SchemaOf<NewInvitedUserAccount>;

export default validationSchema;

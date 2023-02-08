import yup from 'yup';

const LoginFormValidation = yup.object().shape({
  email: yup.string().email().required(),
  password: yup.string().required(),
});

export default LoginFormValidation;

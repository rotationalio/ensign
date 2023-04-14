/* eslint-disable prettier/prettier */
import { useFormik } from 'formik';
import { object, string } from 'yup';

import { AuthUser, InviteAuthUser } from '../types/LoginService';

export const FORM_INITIAL_VALUES = {
  email: '',
  password: '',
  invite_token: '',
} satisfies AuthUser | InviteAuthUser;

export const FORM_VALIDATION_SCHEMA = object({
  email: string().required('Email is required').email('Email is invalid'),
  password: string().required('Password is required'),
});
export const LOGIN_FORM_OPTIONS = (onSubmit: any, initialValues?: any) => ({
  initialValues: initialValues || FORM_INITIAL_VALUES,
  validationSchema: FORM_VALIDATION_SCHEMA,
  onSubmit,
});

export const useLoginForm = (onSubmit: any, initialValues?: any) =>
  useFormik(LOGIN_FORM_OPTIONS(onSubmit, initialValues));

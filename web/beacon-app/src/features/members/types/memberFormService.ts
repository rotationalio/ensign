/* eslint-disable prettier/prettier */
import { useFormik } from 'formik';
import { object, string } from 'yup';

import { MEMBER_ROLE } from '@/constants/rolesAndStatus';

import { MemberRole, NewMemberDTO } from '../types/memberServices';

export const FORM_INITIAL_VALUES = {
  email: '',
  name: '',
  role: MEMBER_ROLE.ADMIN as MemberRole,
} satisfies NewMemberDTO;

export const FORM_VALIDATION_SCHEMA = object({
  email: string().required('Email is required').email('Email is invalid'),
  role: string().required('role is required'),
});
export const FORM_OPTIONS = (onSubmit: any) => ({
  initialValues: FORM_INITIAL_VALUES,
  validationSchema: FORM_VALIDATION_SCHEMA,
  onSubmit,
});

export const ROLE_OPTIONS = [
  { value: 'Owner', label: 'Owner' },
  { value: 'Admin', label: 'Admin' },
  { value: 'Member', label: 'Member' },
  { value: 'Observer', label: 'Observer' },
];

export const useNewMemberForm = (onSubmit: any) => useFormik(FORM_OPTIONS(onSubmit));

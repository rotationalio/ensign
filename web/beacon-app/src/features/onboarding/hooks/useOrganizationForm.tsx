import { t } from '@lingui/macro';
import { useFormik } from 'formik';
import { object, string } from 'yup';

import { MemberResponse } from '@/features/members/types/memberServices';

export type OrganizationFormValues = Pick<MemberResponse, 'organization'>;

export const FORM_INITIAL_VALUES = {
  organization: '',
} as OrganizationFormValues;

export const FORM_VALIDATION_SCHEMA = object({
  organization: string()
    .trim()
    .required(t`Team or organization name is required.`),
});

export const FORM_OPTIONS = (onSubmit: any, initialValues: OrganizationFormValues) => ({
  initialValues: {
    ...FORM_INITIAL_VALUES,
    ...initialValues,
  },
  validationSchema: FORM_VALIDATION_SCHEMA,
  onSubmit,
});

export const useOrganizationForm = (onSubmit: any, initialValues: OrganizationFormValues) =>
  useFormik(FORM_OPTIONS(onSubmit, initialValues));

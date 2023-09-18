/* eslint-disable prettier/prettier */
import { t } from '@lingui/macro';
import { useFormik } from 'formik';
import { object, string } from 'yup';

import { MemberResponse } from '@/features/members/types/memberServices';

export type WorkspaceFormValues = Pick<MemberResponse, 'workspace'>;
export const FORM_INITIAL_VALUES = {
  workspace: '',
} as WorkspaceFormValues;

export const FORM_VALIDATION_SCHEMA = object({
  workspace: string()
    .trim()
    .required(t`Workspace name is required.`), 
});
export const FORM_OPTIONS = (onSubmit: any, initialValues: WorkspaceFormValues) => ({
  initialValues: {
    ...FORM_INITIAL_VALUES,
    ...initialValues,
  },
  validationSchema: FORM_VALIDATION_SCHEMA,
  onSubmit,
});

export const useWorkspaceForm = (onSubmit: any, initialValues: WorkspaceFormValues) =>
  useFormik(FORM_OPTIONS(onSubmit, initialValues));

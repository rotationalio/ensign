/* eslint-disable prettier/prettier */
import { useFormik } from 'formik';
import { boolean, object } from 'yup';

import { DeleteMemberDTO } from './memberServices';

export type DeleteMemberFormValue = {
  name: string;
  delete_agreement: boolean;
} & DeleteMemberDTO;

export const FORM_OPTIONS = (onSubmit: any, values: DeleteMemberFormValue) => ({
  initialValues: values,
  validationSchema: FORM_VALIDATION_SCHEMA,
  onSubmit,
});

export const FORM_VALIDATION_SCHEMA = object({
  delete_agreement: boolean().oneOf([true], 'Please confirm deletion'),
});

export const useDeleteMemberForm = (onSubmit: any, initialValues: DeleteMemberFormValue) =>
  useFormik(FORM_OPTIONS(onSubmit, initialValues));

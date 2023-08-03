import { t } from '@lingui/macro';
import { useFormik } from 'formik';
import * as Yup from 'yup';

export const topicQuerySchema = Yup.object().shape({
  query: Yup.string()
    .trim()
    .required(t`Topic query is required.`),
});

export const QUERY_INPUT_FORM_OPTIONS = (onSubmit: any, initialValues?: any) => ({
  initialValues: initialValues || '',
  validationSchema: topicQuerySchema,
  onSubmit,
});

export const useTopicQueryInputForm = (onSubmit: any, initialValues?: any) =>
  useFormik(QUERY_INPUT_FORM_OPTIONS(onSubmit, initialValues));

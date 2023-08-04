import { t } from '@lingui/macro';
import { useFormik } from 'formik';
import * as Yup from 'yup';

import { ProjectQueryDTO } from '@/features/projects/types/projectQueryService';

export const topicQuerySchema = Yup.object().shape({
  query: Yup.string()
    .trim()
    .required(t`Topic query is required.`),
});

const TOPIC_QUERY_INITIAL_VALUES = {
  query: '',
} satisfies Omit<ProjectQueryDTO, 'projectID'>;

export const QUERY_INPUT_FORM_OPTIONS = (onSubmit: any) => ({
  initialValues: TOPIC_QUERY_INITIAL_VALUES,
  validationSchema: topicQuerySchema,
  onSubmit,
});

export const useTopicQueryInputForm = (onSubmit: any) =>
  useFormik(QUERY_INPUT_FORM_OPTIONS(onSubmit));

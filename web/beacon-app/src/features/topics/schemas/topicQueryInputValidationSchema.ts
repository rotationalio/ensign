import { t } from '@lingui/macro';
import * as Yup from 'yup';

import { ProjectQueryDTO } from '@/features/projects/types/projectQueryService';

export const topicQuerySchema = Yup.object().shape({
  query: Yup.string()
    .trim()
    .required(t`Please enter a query. See EnSQL documentation for examples of valid queries.`),
});

const TOPIC_QUERY_INITIAL_VALUES = {
  query: '',
} satisfies Omit<ProjectQueryDTO, 'projectID'>;

export const QUERY_INPUT_FORM_OPTIONS = (
  onSubmit: any,
  defaultValue: Omit<ProjectQueryDTO, 'projectID'>
) => ({
  initialValues: {
    ...TOPIC_QUERY_INITIAL_VALUES,
    ...defaultValue,
  },
  validationSchema: topicQuerySchema,
  onSubmit,
});

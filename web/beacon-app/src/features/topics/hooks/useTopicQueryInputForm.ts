import { useFormik } from 'formik';

import { ProjectQueryDTO } from '@/features/projects/types/projectQueryService';
import { QUERY_INPUT_FORM_OPTIONS } from '@/features/topics/schemas/topicQueryInputValidationSchema';
export const useTopicQueryInputForm = (
  onSubmit: any,
  initialValues: Omit<ProjectQueryDTO, 'projectID'>
) => useFormik(QUERY_INPUT_FORM_OPTIONS(onSubmit, initialValues));

import { t } from '@lingui/macro';
import { useFormik } from 'formik';
import * as Yup from 'yup';

import { NewTopicDTO } from '../types/createTopicService';

export const createProjectTopicSchema = Yup.object().shape({
  topic_name: Yup.string()
    .trim()
    .required(t`Topic name is required.`)
    .matches(/^[^\s]*$/, t`Topic name cannot include spaces.`)
    .matches(
      /^[a-zA-Z][a-zA-Z0-9._-]*$/,
      t`Topic name may only start with a letter and contain letters, numbers, underscores, and dashes.`
    )
    .max(512, t`Topic name must be less than 512 characters.`),
});

export const FORM_INITIAL_VALUES = {
  topic_name: '',
} satisfies Omit<NewTopicDTO, 'projectID'>;

export const FORM_OPTIONS = (onSubmit: any) => ({
  initialValues: FORM_INITIAL_VALUES,
  validationSchema: createProjectTopicSchema,
  onSubmit,
});

export const useNewTopicForm = (onSubmit: any) => useFormik(FORM_OPTIONS(onSubmit));

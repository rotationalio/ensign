/* eslint-disable prettier/prettier */
import { t } from '@lingui/macro';
import { useFormik } from 'formik';
import { object, string } from 'yup';

import type { Project } from '../types/Project';

export type UpdateProjectOwnerFormDTO = {
  // select options types here
  current_owner: {
    label: string;
    value: string;
  };
  new_owner: {
    label: string;
    value: string;
  };
};

export const FORM_INITIAL_VALUES = {
  current_owner: {
    label: '',
    value: '',
  },
  new_owner: {
    label: '',
    value: '',
  },
} satisfies UpdateProjectOwnerFormDTO;

export const FORM_VALIDATION_SCHEMA = object({
  new_owner: string()
    .required(t`Please select a new owner.`)
    .test('new_owner', t`New owner must be different from current owner`, function (value) {
      const { current_owner } = this.parent;
      return value !== current_owner.value;
    }),
});

export const FORM_OPTIONS = (onSubmit: any, initialValues: Project) => ({
  initialValues: {
    ...FORM_INITIAL_VALUES,
    current_owner: {
      label: initialValues?.owner?.name,
      value: initialValues?.owner?.id,
    },
  },
  validationSchema: FORM_VALIDATION_SCHEMA,
  onSubmit,
});

export const useUpdateProjectOwnerForm = (onSubmit: any, initialValues: Project) =>
  useFormik(FORM_OPTIONS(onSubmit, initialValues));

import { t, Trans } from '@lingui/macro';
import { Button, TextField } from '@rotational/beacon-core';
import { ErrorMessage, Form, FormikHelpers, FormikProvider } from 'formik';
import { useEffect, useState } from 'react';

import { TextArea } from '@/components/ui/TextArea';

import { Project } from '../../types/Project';
import { useUpdateProjectForm } from '../../types/updateProjectFormService';
import type { UpdateProjectDTO } from '../../types/updateProjectService';
type RenameProjectModalFormProps = {
  handleSubmit: (values: UpdateProjectDTO, helpers: FormikHelpers<UpdateProjectDTO>) => void;
  project: Project;
};

function RenameProjectForm({ handleSubmit, project }: RenameProjectModalFormProps) {
  const formik = useUpdateProjectForm(handleSubmit, project);

  const { getFieldProps, isSubmitting, values, touched, errors } = formik;
  const MAX_DESCRIPTION_LENGTH = 2000;
  const [char, setChar] = useState(0);
  const [maxChar, setMaxChar] = useState(MAX_DESCRIPTION_LENGTH);
  useEffect(() => {
    setChar(formik.values?.description?.length || 0);
  }, [formik.values.description]);

  useEffect(() => {
    setMaxChar(MAX_DESCRIPTION_LENGTH - char);
  }, [char]);
  return (
    <FormikProvider value={formik}>
      <Form className="space-y-3">
        <TextField label={t`Current Project Name`} {...getFieldProps('project')} isDisabled />
        <TextField label={t`New Project Name`} {...getFieldProps('name')} />
        <ErrorMessage name="name" component="small" className="text-xs text-danger-500" />
        <TextArea
          label={t`Description (optional)`}
          className="border-transparent bg-[#F7F9FB]"
          rows={5}
          maxLength={2000}
          errorMessage={touched.description && errors.description}
          data-cy="project-description"
          {...getFieldProps('description')}
        />
        {values?.description && values.description.length > 0 && (
          <div className="text-right">
            <span className="text-sm text-gray-600">
              <Trans>Max Length: {maxChar}</Trans>
            </span>
          </div>
        )}
        <div className="pt-3 text-center">
          <Button type="submit" isLoading={isSubmitting} disabled={isSubmitting}>
            {isSubmitting ? t`Renaming project...` : t`Save`}
          </Button>
        </div>
      </Form>
    </FormikProvider>
  );
}

export default RenameProjectForm;

import { t } from '@lingui/macro';
import { Button, TextField } from '@rotational/beacon-core';
import { ErrorMessage, Form, FormikHelpers, FormikProvider } from 'formik';

import { Project } from '../../types/Project';
import { useUpdateProjectForm } from '../../types/updateProjectFormService';
import type { UpdateProjectDTO } from '../../types/updateProjectService';

type RenameProjectModalFormProps = {
  handleSubmit: (values: UpdateProjectDTO, helpers: FormikHelpers<UpdateProjectDTO>) => void;
  project: Project;
};

function RenameProjectForm({ handleSubmit, project }: RenameProjectModalFormProps) {
  const formik = useUpdateProjectForm(handleSubmit, project);

  const { getFieldProps, isSubmitting } = formik;
  return (
    <FormikProvider value={formik}>
      <Form className="space-y-3">
        <TextField label={t`Current Project Name`} {...getFieldProps('project')} isDisabled />
        <TextField label={t`New Project Name`} {...getFieldProps('name')} />
        <ErrorMessage name="name" component="small" className="text-xs text-danger-500" />
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

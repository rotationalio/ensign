import { t } from '@lingui/macro';
import { Button, TextField } from '@rotational/beacon-core';
import { Form, Formik } from 'formik';

import { Project } from '../../types/Project';
type RenameProjectModalFormProps = {
  handleSubmit: (values: any) => void;
  project: Project | null;
};

function RenameProjectModalForm({ handleSubmit, project }: RenameProjectModalFormProps) {
  const initialValues = {
    project: project?.name || '',
  };

  return (
    <Formik onSubmit={handleSubmit} initialValues={initialValues} enableReinitialize>
      {({ getFieldProps, isSubmitting }: any) => (
        <Form className="space-y-3">
          <TextField label={t`Current Project Name`} {...getFieldProps('project')} isDisabled />
          <TextField label={t`New Project Name`} {...getFieldProps('new-name')} />
          <div className="pt-3 text-center">
            <Button type="submit" isLoading={isSubmitting} disabled={isSubmitting}>
              {isSubmitting ? t`Renaming project...` : t`Save`}
            </Button>
          </div>
        </Form>
      )}
    </Formik>
  );
}

export default RenameProjectModalForm;

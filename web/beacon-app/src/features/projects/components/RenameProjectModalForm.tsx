import { TextField } from '@rotational/beacon-core';
import { Form, Formik } from 'formik';

import Button from '@/components/ui/Button/Button';

import { Project } from '../types/Project';

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
      {({ getFieldProps, isSubmitting }) => (
        <Form className="space-y-3">
          <TextField label="Project" {...getFieldProps('project')} isDisabled />
          <TextField label="New project name" {...getFieldProps('new-name')} />
          <div className="pt-3 text-center">
            <Button type="submit" isLoading={isSubmitting} isDisabled={isSubmitting}>
              Save
            </Button>
          </div>
        </Form>
      )}
    </Formik>
  );
}

export default RenameProjectModalForm;

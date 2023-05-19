import { t, Trans } from '@lingui/macro';
import { Button, TextField } from '@rotational/beacon-core';
import { ErrorMessage, Form, FormikHelpers, FormikProvider } from 'formik';
import { useEffect, useState } from 'react';

import { TextArea } from '@/components/ui/TextArea';

import { Project } from '../../types/Project';
import { useUpdateProjectForm } from '../../types/updateProjectFormService';
import type { UpdateProjectDTO } from '../../types/updateProjectService';
type EditProjectModalFormProps = {
  handleSubmit: (values: UpdateProjectDTO, helpers: FormikHelpers<UpdateProjectDTO>) => void;
  project: Project;
};

function EditProjectForm({ handleSubmit, project }: EditProjectModalFormProps) {
  // console.log('[] edit project form', project);
  const formik = useUpdateProjectForm(handleSubmit, project);

  const { getFieldProps, isSubmitting, values, touched, errors } = formik;
  const MAX_DESCRIPTION_LENGTH = 500;
  const [char, setChar] = useState(0);
  const [maxChar, setMaxChar] = useState(MAX_DESCRIPTION_LENGTH);
  const [isDisabled, setIsDisabled] = useState<boolean>(false);
  useEffect(() => {
    setChar(formik.values?.description?.length || 0);
  }, [formik.values.description]);

  useEffect(() => {
    setMaxChar(MAX_DESCRIPTION_LENGTH - char);
  }, [char]);

  useEffect(() => {
    // if any changes are not made to the form, we should disable the submit button
    if (
      formik.values?.name === project?.name &&
      formik.values?.description === project?.description
    ) {
      setIsDisabled(true);
    }
    // disable the submit button if the name is empty and the description is the same
    if (!formik.values?.name && formik.values?.description === project?.description) {
      setIsDisabled(true);
    }
    return () => {
      setIsDisabled(false);
    };
  }, [formik?.values?.name, formik.values?.description, project?.description, project?.name]);
  return (
    <FormikProvider value={formik}>
      <Form className="space-y-3" data-testid="update-project-form">
        <TextField
          label={t`Current Project Name`}
          {...getFieldProps('project')}
          isDisabled
          data-testid="prj-current-name"
        />
        <TextField
          label={t`New Project Name (optional)`}
          data-testid="prj-new-name"
          {...getFieldProps('name')}
        />
        <ErrorMessage name="name" component="small" className="text-xs text-danger-500" />
        <TextArea
          label={t`Description (optional)`}
          labelClassName="text-[12px]"
          className="border-transparent bg-[#F7F9FB]"
          rows={5}
          data-testid="prj-description"
          maxLength={500}
          errorMessage={touched.description && errors.description}
          data-cy="project-description"
          {...getFieldProps('description')}
        />
        {values?.description && values?.description?.length > 0 && (
          <div className="text-right">
            <span className="text-sm text-gray-600">
              <Trans>Max Length: {maxChar}</Trans>
            </span>
          </div>
        )}
        <div className="pt-3 text-center">
          <Button
            type="submit"
            data-testid="update-project-submit"
            isLoading={isSubmitting}
            disabled={isSubmitting || isDisabled}
          >
            {isSubmitting ? t`Renaming project...` : t`Save`}
          </Button>
        </div>
      </Form>
    </FormikProvider>
  );
}

export default EditProjectForm;

import { t, Trans } from '@lingui/macro';
import { Form, FormikHelpers, FormikProvider } from 'formik';
import { useEffect, useState } from 'react';
import styled from 'styled-components';

import Button from '@/components/ui/Button';
import { TextArea } from '@/components/ui/TextArea';
import TextField from '@/components/ui/TextField';

import type { NewProjectDTO } from '../../types/createProjectService';
import { useNewProjectForm } from '../../types/newProjectFormService';

type NewProjectFormProps = {
  onSubmit: (values: NewProjectDTO, helpers: FormikHelpers<NewProjectDTO>) => void;
  isDisabled?: boolean;
  isSubmitting?: boolean;
};

function NewProjectForm({ onSubmit, isSubmitting }: NewProjectFormProps) {
  const formik = useNewProjectForm(onSubmit);

  const [char, setChar] = useState(0);
  const [maxChar, setMaxChar] = useState(512);

  useEffect(() => {
    setChar(formik.values.description.length);
  }, [formik.values.description]);

  useEffect(() => {
    setMaxChar(512 - char);
  }, [char]);

  const { touched, errors, getFieldProps, values } = formik;
  return (
    <FormikProvider value={formik}>
      <Form className="mt-3 mb-2 space-y-2">
        <StyledTextField
          label={t`Project Name (required)`}
          placeholder={t`Enter project name`}
          labelClassName="font-semibold mb-2 text-md"
          className="bg-[#F7F9FB]"
          errorMessage={touched.name && errors.name}
          data-cy="project-name"
          fullWidth
          {...getFieldProps('name')}
        />
        <TextArea
          label={t`Description (optional)`}
          placeholder={t`Enter project description such as the purpose and outcome e.g. to set up events streams for our machine learning model in development.`}
          labelClassName="font-semibold"
          className="border-transparent bg-[#F7F9FB]"
          rows={5}
          maxLength={512}
          errorMessage={touched.description && errors.description}
          {...getFieldProps('description')}
        />
        {values?.description?.length > 0 && (
          <div className="text-right">
            <span className="text-sm text-gray-500">
              <Trans>Length: {maxChar}</Trans>
            </span>
          </div>
        )}
        <div className="pt-3 text-center">
          <Button
            type="submit"
            isLoading={isSubmitting}
            isDisabled={isSubmitting}
            data-cy="NewProjectButton"
          >
            <Trans>Create Project</Trans>
          </Button>
        </div>
      </Form>
    </FormikProvider>
  );
}

const StyledTextField = styled(TextField)`
  border: 'none';
`;

export default NewProjectForm;

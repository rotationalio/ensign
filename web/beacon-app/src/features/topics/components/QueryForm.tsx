import { Trans } from '@lingui/macro';
import { Button, TextField } from '@rotational/beacon-core';
import { Form, FormikHelpers, FormikProvider } from 'formik';

import { ProjectQueryDTO } from '@/features/projects/types/projectQueryService';

import { useTopicQueryInputForm } from '../hooks/useTopicQueryInputForm';

type QueryFormProps = {
  defaultEnSQL: string;
  isSubmitting?: boolean;
  onSubmit: (values: ProjectQueryDTO, formikHelpers: FormikHelpers<ProjectQueryDTO>) => void;
};

const QueryForm = ({ defaultEnSQL, isSubmitting, onSubmit }: QueryFormProps) => {
  const formik = useTopicQueryInputForm(onSubmit, { query: defaultEnSQL });

  const { handleSubmit, getFieldProps, touched, errors } = formik;

  // Using handleReset would reset the form to the defaultEnSQL value when the
  // user clicks clear so we need to manually clear the query field.
  const handleClearQuery = () => {
    formik.setFieldValue('query', '');
  };

  return (
    <FormikProvider value={formik}>
      <Form>
        <div className="mt-4 flex space-x-2">
          <TextField
            fullWidth
            errorMessage={touched.query && errors.query}
            {...getFieldProps('query')}
          />
          <div className="flex max-h-[44px] space-x-2">
            <Button
              variant="secondary"
              onclick={handleSubmit}
              isLoading={isSubmitting}
              disabled={isSubmitting}
            >
              <Trans>Query</Trans>
            </Button>
            <Button onClick={handleClearQuery} disabled={isSubmitting}>
              <Trans>Clear</Trans>
            </Button>
          </div>
        </div>
      </Form>
    </FormikProvider>
  );
};

export default QueryForm;

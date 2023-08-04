import { Trans } from '@lingui/macro';
import { Button, TextField } from '@rotational/beacon-core';
import { Form, FormikHelpers, FormikProvider } from 'formik';

import { ProjectQueryDTO } from '@/features/projects/types/projectQueryService';

import { useTopicQueryInputForm } from '../schemas/topicQueryInputValidationSchema';

type QueryFormProps = {
  defaultEnSQL: string;
  isSubmitting?: boolean;
  onSubmit: (values: ProjectQueryDTO, formikHelpers: FormikHelpers<ProjectQueryDTO>) => void;
};

const QueryForm = ({ defaultEnSQL, isSubmitting, onSubmit }: QueryFormProps) => {
  const formik = useTopicQueryInputForm(onSubmit);

  const { handleReset, handleSubmit, getFieldProps, touched, errors } = formik;

  return (
    <FormikProvider value={formik}>
      <Form>
        <div className="mt-4 flex space-x-2">
          <TextField
            placeholder={defaultEnSQL}
            fullWidth
            errorMessage={touched.query && errors.query}
            {...getFieldProps('query')}
            name="query"
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
            <Button onClick={handleReset} disabled={isSubmitting}>
              <Trans>Clear</Trans>
            </Button>
          </div>
        </div>
      </Form>
    </FormikProvider>
  );
};

export default QueryForm;

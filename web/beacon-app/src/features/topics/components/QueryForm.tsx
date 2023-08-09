import { Trans } from '@lingui/macro';
import { Button, TextField } from '@rotational/beacon-core';
import { Form, FormikHelpers, FormikProvider } from 'formik';

import { ProjectQueryDTO } from '@/features/projects/types/projectQueryService';

import { useTopicQueryInputForm } from '../hooks/useTopicQueryInputForm';

type QueryFormProps = {
  defaultEnSQL: string;
  isSubmitting?: boolean;
  onSubmit: (values: ProjectQueryDTO, formikHelpers: FormikHelpers<ProjectQueryDTO>) => void;
  onReset: () => void;
};

const QueryForm = ({ defaultEnSQL, isSubmitting, onSubmit, onReset }: QueryFormProps) => {
  const formik = useTopicQueryInputForm(onSubmit, { query: defaultEnSQL });

  const { handleSubmit, getFieldProps, touched, errors, setFieldValue } = formik;

  const handleClearQuery = () => {
    setFieldValue('query', '');
    onReset();
  };

  return (
    <FormikProvider value={formik}>
      <Form data-cy="topic-query-form">
        <div className="mt-4 flex space-x-2">
          <TextField
            fullWidth
            errorMessage={touched.query && errors.query}
            {...getFieldProps('query')}
            data-cy="topic-query-input"
          />
          <div className="flex max-h-[44px] space-x-2">
            <Button
              variant="secondary"
              onclick={handleSubmit}
              isLoading={isSubmitting}
              disabled={isSubmitting}
              data-cy="submit-topic-query-bttn"
            >
              <Trans>Query</Trans>
            </Button>
            <Button
              onClick={handleClearQuery}
              disabled={isSubmitting}
              data-cy="clear-topic-query-bttn"
            >
              <Trans>Clear</Trans>
            </Button>
          </div>
        </div>
      </Form>
    </FormikProvider>
  );
};

export default QueryForm;

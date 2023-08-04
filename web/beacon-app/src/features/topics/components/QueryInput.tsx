import { Trans } from '@lingui/macro';
import { Button, TextField } from '@rotational/beacon-core';
import { Form, FormikHelpers, FormikProvider } from 'formik';
import styled from 'styled-components';

import { useTopicQueryInputForm } from '../schemas/topicQueryInputValidationSchema';

type TopicQueryInputProps = {
  name: string;
  onSubmit: (values: any, helpers: FormikHelpers<string>) => void;
  initialValues?: any;
};

const QueryInput = ({ name, onSubmit, initialValues }: TopicQueryInputProps) => {
  const formik = useTopicQueryInputForm(onSubmit, initialValues);

  const { touched, errors, getFieldProps } = formik;

  return (
    <FormikProvider value={formik}>
      <Form>
        <div className="mt-4 flex space-x-2">
          <StyledTextField
            errorMessage={touched.query && errors.query}
            {...getFieldProps('query')}
            type="search"
            fullWidth
          />
          <div className="flex max-h-[44px] space-x-2">
            <Button variant="secondary" type="submit">
              <Trans>Query</Trans>
            </Button>
            <Button type="button">
              <Trans>Clear</Trans>
            </Button>
          </div>
        </div>
      </Form>
    </FormikProvider>
  );
};

const StyledTextField = styled(TextField)`
  input {
    margin-bottom: 1.5rem;
  }
`;

export default QueryInput;

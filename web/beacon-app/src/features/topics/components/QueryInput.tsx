import { Trans } from '@lingui/macro';
import { Button, TextField } from '@rotational/beacon-core';
import { Form, FormikHelpers, FormikProvider } from 'formik';
import { useEffect, useState } from 'react';
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

  const [topicQuery, setTopicQuery] = useState('');

  useEffect(() => {
    setTopicQuery(`SELECT * FROM ${name} LIMIT 10`);
  }, [name]);

  const handleTopicQueryChange = (e: any) => {
    setTopicQuery(e.target.value);
  };

  const handleClearTopicQuery = () => {
    setTopicQuery('');
  };

  return (
    <FormikProvider value={formik}>
      <Form>
        <div className="mt-4 flex space-x-2">
          <StyledTextField
            errorMessage={touched.query && errors.query}
            {...getFieldProps('query')}
            type="search"
            value={topicQuery}
            fullWidth
            onChange={handleTopicQueryChange}
          />
          <div className="flex max-h-[44px] space-x-2">
            <Button variant="secondary" type="submit">
              <Trans>Query</Trans>
            </Button>
            <Button onClick={handleClearTopicQuery} type="button">
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

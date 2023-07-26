import { t, Trans } from '@lingui/macro';
import { Button, TextField } from '@rotational/beacon-core';
import { Form, FormikHelpers, FormikProvider } from 'formik';
import styled from 'styled-components';

import { useNewTopicForm } from '../schemas/createProjectTopicSchema';
import { NewTopicDTO } from '../types/createTopicService';

type NewTopicModalFormProps = {
  onSubmit: (values: NewTopicDTO, helpers: FormikHelpers<NewTopicDTO>) => void;
  isDisabled?: boolean;
  isSubmitting?: boolean;
};

function NewTopicModalForm({ onSubmit, isSubmitting }: NewTopicModalFormProps) {
  const formik = useNewTopicForm(onSubmit);

  const { touched, errors, getFieldProps } = formik;

  return (
    <FormikProvider value={formik}>
      <Form className="mt-3 mb-2 space-y-2">
        <StyledTextField
          label={t`Topic Name (required)`}
          labelClassName="font-semibold mb-2"
          placeholder={t`Enter topic name`}
          fullWidth
          errorMessage={touched.topic_name && errors.topic_name}
          data-cy="topicName"
          {...getFieldProps('topic_name')}
        />
        <div className="text-center">
          <Button
            type="submit"
            isLoading={isSubmitting}
            disabled={isSubmitting}
            data-cy="createTopic"
          >
            <Trans>Create Topic</Trans>
          </Button>
        </div>
      </Form>
    </FormikProvider>
  );
}
// add margin to the text field
const StyledTextField = styled(TextField)`
  margin-bottom: 1rem;
`;

export default NewTopicModalForm;

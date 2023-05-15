import { t, Trans } from '@lingui/macro';
import { Button, TextField } from '@rotational/beacon-core';
import { Form, FormikHelpers, FormikProvider } from 'formik';

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
        <TextField
          label={t`Topic Name (required)`}
          labelClassName="font-semibold"
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

export default NewTopicModalForm;

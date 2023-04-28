import { t, Trans } from '@lingui/macro';
import { Button, TextField } from '@rotational/beacon-core';
import { Form, FormikHelpers, FormikProvider, useFormik } from 'formik';

import { NewTopic, NewTopicDTO } from '@/features/topics/types/topicService';
import { useOrgStore } from '@/store';

import { useCreateTopic } from '../hooks/useCreateTopic';
import createProjectTopicSchema from '../schemas/createProjectTopicSchema';

type NewTopicModalFormProps = {
  onSubmit: (values: NewTopicDTO, helpers: FormikHelpers<NewTopicDTO>) => void;
  isDisabled?: boolean;
  isSubmitting?: boolean;
};

function NewTopicModalForm({ isSubmitting }: NewTopicModalFormProps) {
  const org = useOrgStore.getState() as any;
  const { createTopic } = useCreateTopic();

  const handleCreateTopic = ({ name }: any) => {
    const payload = {
      projectID: org.projectID,
      name,
    } satisfies NewTopicDTO;

    createTopic(payload);
  };

  const formik = useFormik<NewTopic>({
    initialValues: {
      name: '',
    },
    validationSchema: createProjectTopicSchema,
    onSubmit: (values) => {
      handleCreateTopic(values as NewTopicDTO);
    },
  });

  const { touched, errors, getFieldProps } = formik;

  return (
    <FormikProvider value={formik}>
      <Form className="mt-3 mb-2 space-y-2">
        <TextField
          label={t`Topic Name (required)`}
          labelClassName="font-semibold"
          placeholder={t`Enter topic name`}
          fullWidth
          errorMessage={touched.name && errors.name}
          {...getFieldProps('name')}
        />
        {/* TODO: Make button color match Figma design */}
        <div className="text-center">
          <Button
            className="bg-[#6DD19C]"
            type="submit"
            isLoading={isSubmitting}
            disabled={isSubmitting}
          >
            <Trans>Create Topic</Trans>
          </Button>
        </div>
      </Form>
    </FormikProvider>
  );
}

export default NewTopicModalForm;

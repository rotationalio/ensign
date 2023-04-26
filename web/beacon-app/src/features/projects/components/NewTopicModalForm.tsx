import { t, Trans } from '@lingui/macro';
import { Button, TextField } from '@rotational/beacon-core';
import { Form, Formik } from 'formik';

function NewTopicModalForm({ handleSubmit }: { handleSubmit: () => void }) {
  return (
    <Formik onSubmit={handleSubmit} initialValues={{}}>
      {() => (
        <Form className="mt-3 mb-2 space-y-2">
          <TextField
            label={t`Topic Name (required)`}
            labelClassName="font-semibold"
            placeholder={t`Enter topic name`}
            fullWidth
          />
          <div className="text-center">
            <Button>
              <Trans>Create Topic</Trans>
            </Button>
          </div>
        </Form>
      )}
    </Formik>
  );
}

export default NewTopicModalForm;

import { t } from '@lingui/macro';
import { ErrorMessage, Form, FormikHelpers, FormikProvider } from 'formik';

// import { useEffect } from 'react';
import StyledTextField from '@/components/ui/TextField/TextField';
import { useNameForm } from '@/features/onboarding/hooks/useNameForm';

import StepButtons from '../../StepButtons';

type NameFormProps = {
  onSubmit: (values: any, helpers: FormikHelpers<any>) => void;
  isDisabled?: boolean;
  isSubmitting?: boolean;
  initialValues?: any;
  shouldDisableInput?: boolean;
};

const NameForm = ({ onSubmit, isSubmitting, isDisabled, initialValues }: NameFormProps) => {
  const formik = useNameForm(onSubmit, initialValues);
  const { getFieldProps, values, errors } = formik;

  return (
    <FormikProvider value={formik}>
      <Form className="mt-5 space-y-3">
        <StyledTextField
          fullWidth
          placeholder="Ex. Haley Smith"
          label={t`Name`}
          labelClassName="sr-only"
          className="rounded-lg"
          data-cy="user-name"
          {...getFieldProps('name')}
        />
        <ErrorMessage
          name="name"
          component={'p'}
          className="text-error-900 py-2 text-xs text-danger-700"
          data-cy="user-name-error"
        />
        <StepButtons
          isSubmitting={isSubmitting}
          isDisabled={isDisabled || isSubmitting}
          formValues={values}
          hasErrored={!!errors.name}
        />
      </Form>
    </FormikProvider>
  );
};

export default NameForm;

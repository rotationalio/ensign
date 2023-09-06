import { t } from '@lingui/macro';
import { Form, FormikHelpers, FormikProvider } from 'formik';
import { useEffect } from 'react';

import StyledTextField from '@/components/ui/TextField/TextField';
import { useNameForm } from '@/features/onboarding/hooks/useNameForm';

import StepButtons from '../../StepButtons';

type NameFormProps = {
  onSubmit: (values: any, helpers: FormikHelpers<any>) => void;
  isDisabled?: boolean;
  isSubmitting?: boolean;
};

const NameForm = ({ onSubmit, isSubmitting, isDisabled }: NameFormProps) => {
  const formik = useNameForm(onSubmit);
  const { getFieldProps, touched, values, setFieldValue, errors } = formik;

  useEffect(() => {
    if (touched.name && values.name) {
      setFieldValue('name', values.name);
    }
    return () => {
      touched.name = false;
    };
  }, [touched.name, setFieldValue, values, touched]);
  return (
    <FormikProvider value={formik}>
      <Form>
        <StyledTextField
          fullWidth
          placeholder="Ex. Haley Smith"
          label={t`Name`}
          labelClassName="sr-only"
          className="rounded-lg"
          errorMessage={touched.name && errors.name}
          {...getFieldProps('name')}
        />
        <StepButtons isSubmitting={isSubmitting} isDisabled={isDisabled || isSubmitting} />
      </Form>
    </FormikProvider>
  );
};

export default NameForm;

import { t } from '@lingui/macro';
import { Form, FormikHelpers, FormikProvider } from 'formik';

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
  const { getFieldProps, touched, errors } = formik;
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

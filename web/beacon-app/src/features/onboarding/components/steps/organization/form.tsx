import { t } from '@lingui/macro';
import { Form, FormikHelpers, FormikProvider } from 'formik';
import { useEffect } from 'react';

import StyledTextField from '@/components/ui/TextField/TextField';
import { useOrganizationForm } from '@/features/onboarding/hooks/useOrganizationForm';
import { UpdateMemberDTO } from '@/features/onboarding/types/onboardingServices';

import StepButtons from '../../StepButtons';

type OrganizationFormProps = {
  onSubmit: (values: UpdateMemberDTO, helpers: FormikHelpers<any>) => void;
  isDisabled?: boolean;
  isSubmitting?: boolean;
};

const OrganizationForm = ({ onSubmit, isSubmitting, isDisabled }: OrganizationFormProps) => {
  const formik = useOrganizationForm(onSubmit);
  const { getFieldProps, touched, setFieldValue, values, errors } = formik;

  useEffect(() => {
    if (touched.organization && values.organization) {
      setFieldValue('organization', values.organization);
    }
    return () => {
      touched.organization = false;
    };
  }, [touched.organization, setFieldValue, values, touched]);
  return (
    <FormikProvider value={formik}>
      <Form>
        <StyledTextField
          fullWidth
          placeholder="Ex. Rotational Labs"
          label={t`Team or Organization Name`}
          labelClassName="sr-only"
          className="rounded-lg"
          errorMessage={touched.organization && errors.organization}
          {...getFieldProps('organization')}
        />
        <StepButtons isSubmitting={isSubmitting} isDisabled={isDisabled || isSubmitting} />
      </Form>
    </FormikProvider>
  );
};

export default OrganizationForm;

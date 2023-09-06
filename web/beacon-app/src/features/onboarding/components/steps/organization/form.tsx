import { t } from '@lingui/macro';
import { Form, FormikHelpers, FormikProvider } from 'formik';

import StyledTextField from '@/components/ui/TextField/TextField';
import { UpdateMemberDTO } from '@/features/onboarding/types/onboardingServices';
import { useOrganizationForm } from '@/features/onboarding/hooks/useOrganizationForm';

import StepButtons from '../../StepButtons';

type OrganizationFormProps = {
  onSubmit: (values: UpdateMemberDTO, helpers: FormikHelpers<any>) => void;
  isDisabled?: boolean;
  isSubmitting?: boolean;
};

const OrganizationForm = ({ onSubmit, isSubmitting, isDisabled }: OrganizationFormProps) => {
  const formik = useOrganizationForm(onSubmit);
  const { getFieldProps, touched, errors } = formik;
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

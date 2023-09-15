import { t } from '@lingui/macro';
import { ErrorMessage, Form, FormikHelpers, FormikProvider } from 'formik';
import { useEffect } from 'react';

import StyledTextField from '@/components/ui/TextField/TextField';

import { OrganizationFormValues, useOrganizationForm } from '../../../hooks/useOrganizationForm';
import StepButtons from '../../StepButtons';

type OrganizationFormProps = {
  onSubmit: (values: any, helpers: FormikHelpers<any>) => void;
  isDisabled?: boolean;
  isSubmitting?: boolean;
  initialValues?: OrganizationFormValues | any;
  shouldDisableInput?: boolean;
  hasError?: boolean;
};

const OrganizationForm = ({
  onSubmit,
  isSubmitting,
  isDisabled,
  initialValues,
  shouldDisableInput = false,
  hasError,
}: OrganizationFormProps) => {
  const formik = useOrganizationForm(onSubmit, initialValues);
  const { getFieldProps, setFieldError, values } = formik;

  useEffect(() => {
    if (hasError) {
      setFieldError(
        'organization',
        t`The workspace name is taken! Enter a new workspace name or request access from the owner of the ${values.organization} workspace.`
      );
    }
  }, [hasError, setFieldError, values]);

  return (
    <FormikProvider value={formik}>
      <Form>
        <StyledTextField
          fullWidth
          placeholder="Ex. Rotational Labs"
          label={t`Team or Organization Name`}
          labelClassName="sr-only"
          className="rounded-lg"
          disabled={shouldDisableInput}
          {...getFieldProps('organization')}
        />
        <ErrorMessage
          name="organization"
          component={'p'}
          className="text-error-900 py-2 text-xs text-danger-700"
        />
        <StepButtons
          isSubmitting={isSubmitting}
          isDisabled={isDisabled || isSubmitting}
          formValues={values}
        />
      </Form>
    </FormikProvider>
  );
};

export default OrganizationForm;

import { t } from '@lingui/macro';
import { ErrorMessage, Form, FormikHelpers, FormikProvider } from 'formik';
import React, { useEffect } from 'react';

import Select from '@/components/ui/Select';

import { usePreferenceForm } from '../../../hooks/usePreferenceForm';
import { getDeveloperOptions } from '../../../shared/utils';
import StepButtons from '../../StepButtons';
import DeveloperSegment from './DeveloperSegment';
import { ProfessionSegment } from './profession';
type PreferenceFormProps = {
  onSubmit: (values: any, helpers: FormikHelpers<any>) => void;
  isDisabled?: boolean;
  isSubmitting?: boolean;
  initialValues?: any;
  hasError?: boolean;
  error?: any;
};
const UserPreferenceForm = ({
  onSubmit,
  isSubmitting,
  hasError,
  error,
  initialValues,
}: PreferenceFormProps) => {
  const formik = usePreferenceForm(onSubmit, initialValues);
  const [selectedOptions, setSelectedOptions] = React.useState<any[]>([]);
  const { values, setFieldValue, getFieldProps, setFieldError } = formik;
  const ROLE_OPTIONS = getDeveloperOptions();

  const getDefaultOption = () => null;
  const getOptionsAvailable = () =>
    ROLE_OPTIONS?.filter((opt) => opt?.value !== values?.developer_segment?.value);
  const shouldDisableOption = () => selectedOptions?.length >= 3;

  useEffect(() => {
    if (hasError) {
      // error contain developer_segment field error
      if (error?.response?.data?.developer_segment) {
        setFieldError('developer_segment', error?.response?.data?.developer_segment);
      }
      if (error?.response?.data?.profession_segment) {
        setFieldError('profession_segment', error?.response?.data?.profession_segment);
      }
    }
  }, [hasError, error, setFieldError]);

  return (
    <FormikProvider value={formik}>
      <Form className="mt-5 space-y-3">
        <ProfessionSegment
          onChange={(value: any) => {
            setFieldValue('profession_segment', value.value);
          }}
          selectedValue={values?.profession_segment}
        />
        <ErrorMessage
          name={'profession_segment'}
          component={'div'}
          className="text-error-900 py-2 text-xs text-danger-700"
        />
        <fieldset className="my-5">
          <DeveloperSegment />
          <Select
            id="developer_segment"
            inputId="developer_segment"
            placeholder={t`Select one or more options...`}
            className="mt-5"
            defaultValue={getDefaultOption()}
            options={getOptionsAvailable()}
            isOptionDisabled={() => shouldDisableOption()}
            {...getFieldProps('developer_segment')}
            onChange={(value: any) => {
              setFieldValue('developer_segment', value);
              setSelectedOptions(value);
            }}
            isMulti={true}
          />
          <ErrorMessage
            name={'developer_segment'}
            component={'div'}
            className="text-error-900 py-2 text-xs text-danger-700"
          />
        </fieldset>
        <StepButtons isSubmitting={isSubmitting} formValues={values} />
      </Form>
    </FormikProvider>
  );
};

export default UserPreferenceForm;

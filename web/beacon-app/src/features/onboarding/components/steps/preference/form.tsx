import { t } from '@lingui/macro';
import { Form, FormikHelpers, FormikProvider } from 'formik';
import React from 'react';

import Select from '@/components/ui/Select';

import { usePreferenceForm } from '../../../usePreferenceForm';
import { getWorksOptions } from '../../../utils';
import StepButtons from '../../StepButtons';
import DeveloperSegment from './DeveloperSegment';
import { ProfessionSegment } from './profession';
type PreferenceFormProps = {
  onSubmit: (values: any, helpers: FormikHelpers<any>) => void;
  isDisabled?: boolean;
  isSubmitting?: boolean;
  initialValues?: any;
};
const UserPreferenceForm = ({ onSubmit, isSubmitting, initialValues }: PreferenceFormProps) => {
  const formik = usePreferenceForm(onSubmit, initialValues);
  const { values, setFieldValue } = formik;
  const ROLE_OPTIONS = getWorksOptions();
  console.log('ROLE_OPTIONS', ROLE_OPTIONS);

  const getDefaultOption = () => {
    if (ROLE_OPTIONS?.length === 0) {
      return null;
    }
    // select the first option
    return null; // for now lets return null
  };

  const getOptionsAvailable = () =>
    ROLE_OPTIONS?.filter((opt) => opt?.value !== values?.developer_segment?.value);

  return (
    <FormikProvider value={formik}>
      <Form className="mt-5 space-y-3">
        <ProfessionSegment />
        <fieldset className="my-5">
          <DeveloperSegment />
          <Select
            id="developer_segment"
            inputId="developer_segment"
            placeholder={t`Select one or more options...`}
            className="mt-5"
            isDisabled={isSubmitting}
            defaultValue={getDefaultOption()}
            options={getOptionsAvailable()}
            name="developer_segment"
            onChange={(value: any) => setFieldValue('developer_segment', value.value)}
            isMulti={true}
          />
        </fieldset>
        <StepButtons isSubmitting={isSubmitting} />
      </Form>
    </FormikProvider>
  );
};

export default UserPreferenceForm;

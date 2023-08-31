import { t } from '@lingui/macro';
import { ErrorMessage, Form, FormikHelpers, FormikProvider } from 'formik';
import React, { useEffect } from 'react';
import styled from 'styled-components';

import { stringify_org } from '@/utils/slugifyDomain';

import { useNewWorkspaceForm } from '../../../useNewWorkspaceForm';
import StepButtons from '../../StepButtons';

const DOMAIN_BASE = 'https://rotational.io/';
type WorkspaceFormProps = {
  onSubmit: (values: any, helpers: FormikHelpers<any>) => void;
  isDisabled?: boolean;
  isSubmitting?: boolean;
};
const WorkspaceForm = ({ onSubmit, isSubmitting }: WorkspaceFormProps) => {
  const formik = useNewWorkspaceForm(onSubmit);
  const { getFieldProps } = formik;

  // watch workspace field for changes and slugify it
  useEffect(() => {
    const workspace = getFieldProps('workspace').value;
    if (workspace) {
      const slug = stringify_org(workspace);
      formik.setFieldValue('workspace', slug);
    }
  }, [getFieldProps('workspace').value]);

  return (
    <FormikProvider value={formik}>
      <Form className="mt-5 space-y-3">
        <Fieldset>
          <Span className="mt-[3px] font-medium">{DOMAIN_BASE}</Span>

          <StyledTextField placeholder={t`your-workspace-url`} {...getFieldProps('workspace')} />

          <div>
            <ErrorMessage
              name={'workspace'}
              component={'div'}
              className="text-error-900 py-2 text-xs text-danger-700"
            />
          </div>
        </Fieldset>

        <StepButtons isSubmitting={isSubmitting} />
      </Form>
    </FormikProvider>
  );
};

const Fieldset = styled.fieldset`
  display: flex;
  position: relative;
  border-radius: 0.5rem;
  padding-top: 25px;

  & div label {
    position: absolute;
    top: 0;
    left: 0;
  }
  & input {
    border-top-left-radius: 0px;
    border-bottom-left-radius: 0px;
    border-left: none;
    padding-left: 0;
    margin-top: 3px !important;
  }
  & input:focus {
    border: none;
    outline: none;
  }
  & div {
    position: static;
    flex-grow: 1;
  }
  & div > div {
    position: absolute;
    bottom: -13px;
    left: 145px;
    padding: 0 0.2rem;
  }
`;

const Span = styled.span`
  display: flex;
  align-items: center;
  border-right: none;
  color: gray;
  border-top-left-radius: 0.375rem /* 6px */;
  border-bottom-left-radius: 0.375rem /* 6px */;
  padding-left: 1rem;
  white-space: nowrap;
`;

const StyledTextField = styled.input`
  border: none;
  padding: 0.5rem;
  border-radius: 8px;
  outline: none;
  width: 100%;

  &:focus,
  &:active {
    outline: none;
    border: none;
    background-color: #fff;
  }
`;

export default WorkspaceForm;

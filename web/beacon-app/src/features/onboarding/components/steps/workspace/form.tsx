import { t } from '@lingui/macro';
import { ErrorMessage, Form, FormikHelpers, FormikProvider } from 'formik';
import React, { useEffect } from 'react';
import styled from 'styled-components';

import { stringify_org } from '@/utils/slugifyDomain';

import { useWorkspaceForm } from '../../../hooks/useWorkspaceForm';
import StepButtons from '../../StepButtons';

const DOMAIN_BASE = 'https://rotational.app/';
type WorkspaceFormProps = {
  onSubmit: (values: any, helpers: FormikHelpers<any>) => void;
  isDisabled?: boolean;
  isSubmitting?: boolean;
  initialValues?: any;
  hasError?: boolean;
};
const WorkspaceForm = ({ onSubmit, isSubmitting, initialValues, hasError }: WorkspaceFormProps) => {
  const formik = useWorkspaceForm(onSubmit, initialValues);
  const { getFieldProps, touched, setFieldValue, values, setFieldError } = formik;

  useEffect(() => {
    if (touched.workspace && values.workspace) {
      setFieldValue('workspace', stringify_org(values.workspace));
    }
    return () => {
      touched.workspace = false;
    };
  }, [touched.workspace, setFieldValue, values, touched]);

  // set error if workspace is already taken
  useEffect(() => {
    if (hasError) {
      setFieldError(
        'workspace',
        t`The workspace URL is taken by another team. Try a variation or another slug.
 `
      );
    }
  }, [hasError, setFieldError]);

  return (
    <FormikProvider value={formik}>
      <Form className="mt-5 space-y-3">
        <Fieldset>
          <Span className="mt-[3px] font-medium">{DOMAIN_BASE}</Span>

          <StyledTextField placeholder={'rotational-labs'} {...getFieldProps('workspace')} />

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
  padding: 5px;
  border: 1px solid #e5e7eb;

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
    bottom: -17px;
    left: 160px;
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
  font-weight: 500;

  &:focus,
  &:active {
    outline: none;
    border: none;
    background-color: #fff;
  }
`;

export default WorkspaceForm;

/* eslint-disable simple-import-sort/imports */
import { AriaButton as Button, Checkbox, TextField } from '@rotational/beacon-core';
import Tooltip from '@rotational/beacon-core/lib/components/Tooltip';
import { Form, FormikHelpers, FormikProvider, useFormik } from 'formik';
import { Link } from 'react-router-dom';
import styled from 'styled-components';

import HelpIcon from '@/components/icons/help-icon';
import { PasswordStrength } from '@/components/PasswordStrength';
import { stringify_org } from '@/utils/slugifyDomain';
import registrationFormValidationSchema from './schemas/registrationFormValidation';
import { NewUserAccount } from '../../types/RegisterService';
import { useState } from 'react';
import { OpenEyeIcon } from '@/components/icons/openEyeIcon';
import { CloseEyeIcon } from '@/components/icons/closeEyeIcon';
import { EXTRENAL_LINKS } from '@/application/routes/paths';
const initialValues = {
  name: '',
  email: '',
  password: '',
  pwcheck: '',
  organization: '',
  domain: '',
  terms_agreement: false,
  privacy_agreement: false,
} satisfies NewUserAccount;

type RegistrationFormProps = {
  onSubmit: (values: NewUserAccount, helpers: FormikHelpers<NewUserAccount>) => void;
};

function RegistrationForm({ onSubmit }: RegistrationFormProps) {
  const formik = useFormik<NewUserAccount>({
    initialValues,
    onSubmit,
    validationSchema: registrationFormValidationSchema,
  });
  const { touched, errors, values, getFieldProps, setFieldValue, isSubmitting } = formik;

  const handlePasswordMatch = (_result: boolean) => {
    // console.log('result', result)
  };
  console.log('values', values);

  const [openEyeIcon, setOpenEyeIcon] = useState(false);

  const toggleEyeIcon = () => {
    setOpenEyeIcon(!openEyeIcon);
  };

  return (
    <FormikProvider value={formik}>
      <Form>
        <div className="mb-4 space-y-4">
          <TextField
            label={`Name (required)`}
            placeholder="Holly Golightly"
            data-testid="name"
            fullWidth
            errorMessage={touched.name && errors.name}
            {...getFieldProps('name')}
          />
          <TextField
            label={`Email address (required)`}
            placeholder="holly@golight.ly"
            fullWidth
            data-testid="email"
            errorMessage={touched.email && errors.email}
            {...getFieldProps('email')}
          />
          <div className="relative">
            <TextField
              label={`Password`}
              placeholder={`Password`}
              type={!openEyeIcon ? 'password' : 'text'}
              data-testid="password"
              errorMessage={touched.password && errors.password}
              fullWidth
              {...getFieldProps('password')}
            />
            <button
              onClick={toggleEyeIcon}
              className="absolute right-2 top-8 h-8 pb-2"
              data-testid="button"
            >
              {openEyeIcon ? <OpenEyeIcon /> : <CloseEyeIcon />}
              <span className="sr-only" data-testid="screenReadText">
                {openEyeIcon ? 'Hide Password' : 'Show Password'}
              </span>
            </button>
          </div>
          {touched.password && values.password ? (
            <PasswordStrength string={values.password} onMatch={handlePasswordMatch} />
          ) : null}
          <TextField
            label={`Confirm Password`}
            placeholder={`Password`}
            type="password"
            fullWidth
            data-testid="pwcheck"
            errorMessage={touched.pwcheck && errors.pwcheck}
            {...getFieldProps('pwcheck')}
          />
          <TextField
            label={
              <span className="flex items-center gap-2">
                <span>Organization (required)</span>
                <Tooltip
                  title={
                    <span>
                      Your organization allows you to collaborate with teammates and set up multiple
                      tenants and projects.
                    </span>
                  }
                >
                  <HelpIcon className="w-4" />
                </Tooltip>
              </span>
            }
            placeholder="Team Diamonds"
            fullWidth
            data-testid="organization"
            errorMessage={touched.organization && errors.organization}
            {...getFieldProps('organization')}
          />
          <Fieldset>
            <Span>
              https://rotational.app/
              {stringify_org(values.organization) || 'organization Inc'}/
            </Span>
            <TextField
              label={
                <span className="flex items-center gap-2">
                  <span>Domain</span>
                  <Tooltip
                    title={
                      <span>
                        Your domain is a universal resource locator for use across the Ensign
                        ecosystem.
                      </span>
                    }
                  >
                    <HelpIcon className="w-4" />
                  </Tooltip>
                </span>
              }
              placeholder="breakfast.tiffany.io"
              fullWidth
              errorMessage={touched.domain && errors.domain}
              {...getFieldProps('domain')}
              data-testid="domain"
            />
          </Fieldset>
        </div>
        <CheckboxFieldset>
          <Checkbox
            name="terms_agreement"
            onChange={(isSelected) => {
              setFieldValue('terms_agreement', isSelected);
              setFieldValue('privacy_agreement', isSelected);
            }}
            data-testid="terms_agreement"
          >
            I agree to the Rotational Labs Inc.{' '}
            <Link to={EXTRENAL_LINKS.TERMS} className="font-bold underline" target="_blank">
              Terms of Service
            </Link>{' '}
            and{' '}
            <Link to={EXTRENAL_LINKS.PRIVACY} className="font-bold underline" target="_blank">
              Privacy Policy
            </Link>
            .
          </Checkbox>
          <div>{touched.terms_agreement && errors.terms_agreement}</div>
        </CheckboxFieldset>
        <Button
          type="submit"
          color="secondary"
          className="mt-4"
          isDisabled={isSubmitting}
          aria-label="Create Starter account"
        >
          Create Starter Account
        </Button>
        <p className="mt-2">No cost. No credit card required.</p>
      </Form>
    </FormikProvider>
  );
}

const Fieldset = styled.fieldset`
  display: flex;
  position: relative;
  border-radius: 0.5rem;
  padding-top: 25px;
  padding-bottom: 17px;
  overflow: hidden;
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
  }
  & div {
    position: static;
  }
  & div > div {
    position: absolute;
    bottom: 0;
    left: 0;
  }
`;

const Span = styled.span`
  display: flex;
  align-items: center;
  border: 1px solid black;
  border-right: none;
  color: gray;
  border-top-left-radius: 0.375rem /* 6px */;
  border-bottom-left-radius: 0.375rem /* 6px */;
  padding-left: 1rem;
  width: 60%;
  white-space: nowrap;
`;

// TODO: fix it in the design system
const CheckboxFieldset = styled.fieldset`
  label svg {
    min-width: 23px;
  }
`;

export default RegistrationForm;

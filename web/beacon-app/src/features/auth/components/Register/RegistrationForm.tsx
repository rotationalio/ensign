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
import { Trans, t } from '@lingui/macro';

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

  const handlePasswordMatch = (_result: boolean) => {};

  const [openEyeIcon, setOpenEyeIcon] = useState(false);

  const toggleEyeIcon = () => {
    setOpenEyeIcon(!openEyeIcon);
  };

  return (
    <FormikProvider value={formik}>
      <Form>
        <div className="mb-4 space-y-4">
          <TextField
            label={t`Name (required)`}
            placeholder="Holly Golightly"
            data-testid="name"
            fullWidth
            errorMessage={touched.name && errors.name}
            {...getFieldProps('name')}
          />
          <TextField
            label={t`Email address (required)`}
            placeholder="holly@golight.ly"
            fullWidth
            data-testid="email"
            errorMessage={touched.email && errors.email}
            {...getFieldProps('email')}
          />
          <div className="relative">
            <TextField
              label={t`Password`}
              placeholder={t`Password`}
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
                {openEyeIcon ? <Trans>Hide Password</Trans> : <Trans>Show Password</Trans>}
              </span>
            </button>
          </div>
          {touched.password && values.password ? (
            <PasswordStrength string={values.password} onMatch={handlePasswordMatch} />
          ) : null}
          <TextField
            label={t`Confirm Password`}
            placeholder={t`Password`}
            type="password"
            fullWidth
            data-testid="pwcheck"
            errorMessage={touched.pwcheck && errors.pwcheck}
            {...getFieldProps('pwcheck')}
          />
          <TextField
            label={
              <span className="flex items-center gap-2">
                <span>
                  <Trans>Organization (required)</Trans>
                </span>
                <Tooltip
                  title={
                    <span>
                      <Trans>
                        Your organization allows you to collaborate with teammates and set up
                        multiple tenants and projects.
                      </Trans>
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
                  <span>
                    <Trans>Domain</Trans>
                  </span>
                  <Tooltip
                    title={
                      <span>
                        <Trans>
                          Your domain is a universal resource locator for use across the Ensign
                          ecosystem.
                        </Trans>
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
            <Trans>
              I agree to the Rotational Labs Inc.{' '}
              <Link to="/#" className="font-bold underline">
                Terms of Service
              </Link>{' '}
              and{' '}
              <Link to="/#" className="font-bold underline">
                Privacy Policy
              </Link>
              .
            </Trans>
          </Checkbox>
          <div>{touched.terms_agreement && errors.terms_agreement}</div>
        </CheckboxFieldset>
        <Button
          type="submit"
          color="secondary"
          className="mt-4"
          isDisabled={isSubmitting}
          aria-label={t`Create Starter account`}
        >
          <Trans>Create Starter Account</Trans>
        </Button>
        <p className="mt-2">
          <Trans>No cost. No credit card required.</Trans>
        </p>
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

import * as RadixTooltip from '@radix-ui/react-tooltip';
import { Checkbox } from '@rotational/beacon-core';
import Tooltip from '@rotational/beacon-core/lib/components/Tooltip';
import { ErrorMessage, Form, FormikHelpers, FormikProvider, useFormik } from 'formik';
import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import styled from 'styled-components';

import { EXTRENAL_LINKS } from '@/application/routes/paths';
import HelpIcon from '@/components/icons/help-icon';
import { PasswordStrength } from '@/components/PasswordStrength';
import Button from '@/components/ui/Button';
import PasswordField from '@/components/ui/PasswordField/PasswordField';
import TextField from '@/components/ui/TextField';
import useFocus from '@/hooks/useFocus';
import { stringify_org } from '@/utils/slugifyDomain';

import { NewUserAccount } from '../../types/RegisterService';
import registrationFormValidationSchema from './schemas/registrationFormValidation';

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

const DOMAIN_BASE = 'https://rotational.app/';

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
  const [isFocused, { onBlur, onFocus }] = useFocus();
  // eslint-disable-next-line unused-imports/no-unused-vars
  const [isPasswordMatchOpen, setIsPasswordMatchOpen] = useState<boolean | undefined>(
    !!values.password
  );

  const handlePasswordMatch = (_result: boolean) => {
    // console.log('result', result)
  };

  useEffect(() => {
    setIsPasswordMatchOpen(!!values.password);
    setTimeout(() => {
      setIsPasswordMatchOpen(undefined);
    }, 10000);
  }, [values.password]);

  // if organization name is set then set domain to the slugified version of the organization name
  useEffect(() => {
    if (touched.organization && !touched.domain) {
      setFieldValue('domain', stringify_org(values.organization));
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [touched.organization, touched.domain]);

  useEffect(() => {
    if (values.domain) {
      setFieldValue('domain', stringify_org(values.domain));
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [values.domain]);

  return (
    <FormikProvider value={formik}>
      <Form>
        <div className="mb-1 space-y-2">
          <TextField
            label={`Name (required)`}
            placeholder="Holly Golightly"
            data-testid="name"
            fullWidth
            errorMessage={touched.name && errors.name}
            errorMessageClassName="py-2"
            {...getFieldProps('name')}
          />
          <TextField
            label={`Email address (required)`}
            placeholder="holly@golight.ly"
            fullWidth
            data-testid="email"
            errorMessage={touched.email && errors.email}
            errorMessageClassName="py-2"
            {...getFieldProps('email')}
          />
          <div className="relative">
            <PasswordField
              label={
                <RadixTooltip.Provider>
                  <RadixTooltip.Root open={isFocused}>
                    <span className="-my-1 flex items-center gap-2">
                      Password
                      <RadixTooltip.Trigger asChild>
                        <button className="flex">
                          <HelpIcon className="w-4" />
                        </button>
                      </RadixTooltip.Trigger>
                    </span>
                    <RadixTooltip.Portal>
                      <RadixTooltip.Content
                        className="text-violet11 select-none rounded-[4px] bg-white px-[15px] py-[10px] text-[15px] leading-none shadow-[hsl(206_22%_7%_/_35%)_0px_10px_38px_-10px,_hsl(206_22%_7%_/_20%)_0px_10px_20px_-15px] will-change-[transform,opacity] data-[state=delayed-open]:data-[side=top]:animate-slideDownAndFade data-[state=delayed-open]:data-[side=right]:animate-slideLeftAndFade data-[state=delayed-open]:data-[side=left]:animate-slideRightAndFade data-[state=delayed-open]:data-[side=bottom]:animate-slideUpAndFade"
                        sideOffset={5}
                      >
                        <PasswordStrength string={values.password} onMatch={handlePasswordMatch} />
                        <RadixTooltip.Arrow className="fill-white" />
                      </RadixTooltip.Content>
                    </RadixTooltip.Portal>
                  </RadixTooltip.Root>
                </RadixTooltip.Provider>
              }
              placeholder={`Password`}
              data-testid="password"
              errorMessage={touched.password && errors.password}
              errorMessageClassName="py-2"
              fullWidth
              {...getFieldProps('password')}
              onFocus={onFocus}
              onBlur={onBlur}
            />
          </div>
          <TextField
            label={`Confirm Password`}
            placeholder={`Password`}
            type="password"
            fullWidth
            data-testid="pwcheck"
            errorMessage={touched.pwcheck && errors.pwcheck}
            errorMessageClassName="py-2"
            {...getFieldProps('pwcheck')}
          />
          <TextField
            label={
              <span className="-my-1 flex items-center gap-2">
                <span>Organization (required)</span>
                <TooltipSpan>
                  <Tooltip
                    title={
                      <>
                        Your organization allows you to collaborate with teammates and set up
                        multiple tenants and projects.
                      </>
                    }
                  >
                    <HelpIcon className="w-4" />
                  </Tooltip>
                </TooltipSpan>
              </span>
            }
            placeholder="Team Diamonds"
            fullWidth
            data-testid="organization"
            errorMessage={touched.organization && errors.organization}
            errorMessageClassName="py-2"
            {...getFieldProps('organization')}
          />
          <Fieldset>
            <Span className="mt-[3px]">{DOMAIN_BASE}</Span>
            <TextField
              label={
                <span className=" flex items-center gap-2">
                  <span>Domain</span>
                  <Tooltip
                    title={
                      <span className="text-sm">
                        Your domain is a universal resource locator for use across the Ensign
                        ecosystem.
                      </span>
                    }
                  >
                    <HelpIcon className="w-4" />
                  </Tooltip>
                </span>
              }
              placeholder="organization name"
              fullWidth
              data-testid="domain"
              errorMessage={touched.domain && errors.domain}
              errorMessageClassName="-pt-2"
              className="mt-0"
              {...getFieldProps('domain')}
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
          <ErrorMessage component="p" name="terms_agreement" className="text-xs text-danger-500" />
        </CheckboxFieldset>
        <div>
          <TextField type="hidden" {...getFieldProps('domain')} data-testid="domain" />
        </div>
        <Button
          type="submit"
          variant="secondary"
          size="large"
          className="mt-4"
          isLoading={isSubmitting}
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
    bottom: -7px;
    left: 0;
  }
`;

const Span = styled.span`
  display: flex;
  align-items: center;
  background-color: #f5f5f5;
  border: 1px solid black;
  border-right: none;
  color: gray;
  border-top-left-radius: 0.375rem /* 6px */;
  border-bottom-left-radius: 0.375rem /* 6px */;
  padding-left: 1rem;
  width: 200px;
  white-space: nowrap;
`;

// TODO: fix it in the design system
const CheckboxFieldset = styled.fieldset`
  margin-top: 1rem;
  label svg {
    min-width: 23px;
  }
`;

const TooltipSpan = styled.span`
  & span {
    display: flex;
    align-items: center;
  }
`;

export default RegistrationForm;

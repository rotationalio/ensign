import { t, Trans } from '@lingui/macro';
import * as RadixTooltip from '@radix-ui/react-tooltip';
import { Button, Checkbox } from '@rotational/beacon-core';
import Tooltip from '@rotational/beacon-core/lib/components/Tooltip';
import { ErrorMessage, Form, FormikHelpers, FormikProvider, useFormik } from 'formik';
import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import useMedia from 'react-use/lib/useMedia';
import styled from 'styled-components';

import { EXTRENAL_LINKS } from '@/application/routes/paths';
import HelpIcon from '@/components/icons/help-icon';
import { PasswordStrength } from '@/components/PasswordStrength';
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
  const isMobile = useMedia('(max-width: 860px)');

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
            <RadixTooltip.Provider>
              <RadixTooltip.Root open={isFocused}>
                <RadixTooltip.Trigger asChild>
                  <div>
                    <PasswordField
                      placeholder={t`Password`}
                      data-testid="password"
                      errorMessage={touched.password && errors.password}
                      fullWidth
                      {...getFieldProps('password')}
                      onFocus={onFocus}
                      onBlur={onBlur}
                    />
                  </div>
                </RadixTooltip.Trigger>
                <RadixTooltip.Portal>
                  <RadixTooltip.Content
                    className="select-none rounded-[4px] bg-white px-[15px] py-[10px] text-xs text-[15px] leading-none shadow-[hsl(206_22%_7%_/_35%)_0px_10px_38px_-10px,_hsl(206_22%_7%_/_20%)_0px_10px_20px_-15px] will-change-[transform,opacity] data-[state=delayed-open]:data-[side=top]:animate-slideDownAndFade data-[state=delayed-open]:data-[side=right]:animate-slideLeftAndFade data-[state=delayed-open]:data-[side=left]:animate-slideRightAndFade data-[state=delayed-open]:data-[side=bottom]:animate-slideUpAndFade"
                    sideOffset={2}
                    side={isMobile ? 'bottom' : 'right'}
                  >
                    <PasswordStrength string={values.password} onMatch={handlePasswordMatch} />
                    <RadixTooltip.Arrow className="fill-white" />
                  </RadixTooltip.Content>
                </RadixTooltip.Portal>
              </RadixTooltip.Root>
            </RadixTooltip.Provider>
          </div>
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
              <span className="-my-1 flex items-center gap-2">
                <span>
                  <Trans>Organization (required)</Trans>
                </span>
                <TooltipSpan>
                  <Tooltip
                    title={
                      <span className="text-xs">
                        <Trans>
                          Your organization allows you to collaborate with teammates and set up
                          multiple tenants and projects.
                        </Trans>
                      </span>
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
            {...getFieldProps('organization')}
          />
          <Fieldset>
            <Span className="mt-[3px]">{DOMAIN_BASE}</Span>
            <TextField
              label={
                <span className="flex items-center gap-2">
                  <span>
                    <Trans>Domain</Trans>
                  </span>
                  <Tooltip
                    title={
                      <span className="text-xs">
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
              placeholder="organization name"
              fullWidth
              data-testid="domain"
              {...getFieldProps('domain')}
            />
          </Fieldset>
          <ErrorMessage name={'domain'} component={'small'} className="text-xs text-danger-700" />
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
            <Trans>I agree to the Rotational Labs Inc.</Trans>{' '}
            <Link to={EXTRENAL_LINKS.TERMS} className="font-bold underline" target="_blank">
              <Trans>Terms of Service</Trans>
            </Link>{' '}
            <Trans>and</Trans>{' '}
            <Link to={EXTRENAL_LINKS.PRIVACY} className="font-bold underline" target="_blank">
              <Trans>Privacy Policy</Trans>
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
          size="medium"
          className="mt-4"
          isLoading={isSubmitting}
          disabled={isSubmitting}
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
  & div {
    position: static;
    flex-grow: 1;
    width: 0;
  }
  & div > div {
    position: absolute;
    bottom: -13px;
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

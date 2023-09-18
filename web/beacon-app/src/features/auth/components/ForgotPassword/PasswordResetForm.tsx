import { t, Trans } from '@lingui/macro';
import * as RadixTooltip from '@radix-ui/react-tooltip';
import { Button } from '@rotational/beacon-core';
import { ErrorMessage, Form, FormikHelpers, FormikProvider } from 'formik';
import { useEffect, useState } from 'react';
import useMedia from 'react-use/lib/useMedia';

import { PasswordStrength } from '@/components/PasswordStrength';
import PasswordField from '@/components/ui/PasswordField';
import StyledTextField from '@/components/ui/TextField/TextField';
import useFocus from '@/hooks/useFocus';

import { usePasswordResetForm } from './hooks/usePasswordResetForm';

type PasswordResetFormProps = {
  onSubmit: (values: any, helpers: FormikHelpers<any>) => void;
};

const PasswordResetForm = ({ onSubmit }: PasswordResetFormProps) => {
  const formik = usePasswordResetForm(onSubmit);
  const { getFieldProps, values } = formik;
  const isMobile = useMedia('(max-width: 860px)');

  const [isFocused, { onBlur, onFocus }] = useFocus();
  // eslint-disable-next-line unused-imports/no-unused-vars
  const [isPasswordMatchOpen, setIsPasswordMatchOpen] = useState<boolean | undefined>(
    !!values.password
  );

  useEffect(() => {
    setIsPasswordMatchOpen(!!values.password);
    setTimeout(() => {
      setIsPasswordMatchOpen(undefined);
    }, 10000);
  }, [values.password]);

  return (
    <FormikProvider value={formik}>
      <Form>
        <div className="relative">
          <RadixTooltip.Provider>
            <RadixTooltip.Root open={isFocused}>
              <RadixTooltip.Trigger asChild>
                <div>
                  <PasswordField
                    placeholder={t`Password`}
                    data-testid="password"
                    data-cy="password"
                    fullWidth
                    className="mb-2"
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
                  <PasswordStrength string={values.password} />
                  <RadixTooltip.Arrow className="fill-white" />
                </RadixTooltip.Content>
              </RadixTooltip.Portal>
            </RadixTooltip.Root>
          </RadixTooltip.Provider>
          <ErrorMessage name="password" component={'p'} className="text-xs text-danger-700" />
        </div>
        <PasswordField
          label={t`Confirm Password`}
          placeholder={t`Password`}
          type="password"
          fullWidth
          className="mb-2"
          data-testid="pwcheck"
          data-cy="pwcheck"
          {...getFieldProps('pwcheck')}
        />
        <ErrorMessage name="pwcheck" component={'p'} className="text-xs text-danger-700" />
        <StyledTextField className="hidden" {...getFieldProps('reset_token')} />
        <div className="mt-3 flex justify-between">
          <div></div>
          <Button type="submit" variant="secondary" className="mt-2" data-cy="reset-password-bttn">
            <Trans>Submit</Trans>
          </Button>
        </div>
      </Form>
    </FormikProvider>
  );
};

export default PasswordResetForm;

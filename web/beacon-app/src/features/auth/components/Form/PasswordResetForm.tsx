import { t, Trans } from '@lingui/macro';
import { Button } from '@rotational/beacon-core';
import { ErrorMessage, Form, FormikHelpers, FormikProvider } from 'formik';
import useMedia from 'react-use/lib/useMedia';

import { PasswordStrength } from '@/components/PasswordStrength';
import PasswordField from '@/components/ui/PasswordField';
import PasswordTooltip from '@/components/ui/PasswordTooltip/PasswordTooltip';
import StyledTextField from '@/components/ui/TextField/TextField';
import useFocus from '@/hooks/useFocus';

import { usePasswordResetForm } from '../../hooks/usePasswordResetForm';

type PasswordResetFormProps = {
  onSubmit: (values: any, helpers: FormikHelpers<any>) => void;
};

const PasswordResetForm = ({ onSubmit }: PasswordResetFormProps) => {
  const formik = usePasswordResetForm(onSubmit);
  const { getFieldProps, values } = formik;
  const isMobile = useMedia('(max-width: 860px)');

  const [isFocused, { onBlur, onFocus }] = useFocus();

  return (
    <FormikProvider value={formik}>
      <Form>
        <div className="relative">
          <PasswordTooltip
            isFocused={isFocused}
            isMobile={isMobile}
            triggerNode={
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
            }
            contentNode={<PasswordStrength string={values.password} />}
          />
          <ErrorMessage
            name="password"
            component={'p'}
            className="text-xs text-danger-700"
            data-cy="password-error-msg"
          />
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
        <ErrorMessage
          name="pwcheck"
          component={'p'}
          className="text-xs text-danger-700"
          data-cy="pwcheck-error-msg"
        />
        <StyledTextField className="hidden" {...getFieldProps('reset_token')} />
        <div className="mt-3 flex justify-between">
          <div></div>
          <Button
            type="submit"
            variant="secondary"
            className="mt-2"
            data-cy="reset-password-submit-bttn"
          >
            <Trans>Submit</Trans>
          </Button>
        </div>
      </Form>
    </FormikProvider>
  );
};

export default PasswordResetForm;

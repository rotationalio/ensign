import { Trans } from '@lingui/macro';
import React, { memo } from 'react';
import { Link } from 'react-router-dom';
type Props = {
  ButtonElement: any;
  isAuthenticating: boolean;
};
const LoginFooter = ({ ButtonElement, isAuthenticating }: Props) => {
  return (
    <>
      <div className="space-y-4 rounded-md border border-[#1D65A6] bg-[#1D65A6] p-4 text-white sm:p-8 md:w-[402px]">
        <h1 className="text-center font-bold">
          <Trans>Need an Account?</Trans>
        </h1>

        <ul className="ml-5 list-disc">
          <li>
            <Trans>Set up your first event stream in minutes</Trans>
          </li>
          <li>
            <Trans>No DevOps foo needed</Trans>
          </li>
          <li>
            <Trans>Goodbye YAML!</Trans>
          </li>
          <li>
            <Trans>We 🤍 SDKs</Trans>
          </li>
          <li>
            <Trans>Learn from beginner-friendly examples</Trans>
          </li>
          <li>
            <Trans>No credit card required</Trans>
          </li>
          <li>
            <Trans>Cancel any time</Trans>
          </li>
        </ul>

        <div className="flex justify-center">
          <Link to="/register">
            <ButtonElement
              variant="ghost"
              disabled={isAuthenticating}
              className="mt-4"
              data-testid="get__started"
            >
              <Trans>Get Started</Trans>
            </ButtonElement>
          </Link>
        </div>
      </div>
    </>
  );
};

export default memo(LoginFooter);

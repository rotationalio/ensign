import { Trans } from '@lingui/macro';
import { Container } from '@rotational/beacon-core';

import OtterFloating from '@/assets/images/otter-floating.svg';

const ResetVerificationPage = () => {
  return (
    <>
      <Container className="my-20">
        <div className="mx-auto min-h-min max-w-2xl rounded-lg border border-solid border-primary-800">
          <div className="pt-10">
            <div className="ml-8 text-left">
              <Trans>
                <p>
                  Thank you. We have sent you instructions to reset your password. The link to reset
                  your password will expire in 1 hour.
                </p>
              </Trans>
              <div className="text-right">
                <img src={OtterFloating} alt="" className="inline scale-75" />
              </div>
            </div>
          </div>
        </div>
      </Container>
    </>
  );
};

export default ResetVerificationPage;

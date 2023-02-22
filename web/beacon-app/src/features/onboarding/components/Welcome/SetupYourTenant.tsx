import { Trans } from '@lingui/macro';
import { AriaButton } from '@rotational/beacon-core';
import { Link } from 'react-router-dom';

import { ROUTES } from '@/application';
import HostedDataIcon from '@/components/icons/hosted-data-icon';

function SetupYourTenant() {
  return (
    <section className="mx-auto grid max-w-4xl grid-cols-3 rounded-lg border border-solid border-primary-800 py-6 text-2xl">
      <HostedDataIcon className="mx-auto mt-6" />
      <div>
        <h2 className="mt-8 font-bold">
          <Trans>
            Set Up Your Tenant <span className="font-normal">(required)</span>
          </Trans>
        </h2>
        <p className="mt-8">
          <Trans>
            Your tenant is your team&apos;s control panel for all projects and topics. Specify
            preferences around encryption, privacy, and locality (e.g. for GDPR, CCPA, etc).
          </Trans>
        </p>
      </div>
      <div className="mx-auto mt-36">
        <Link to={ROUTES.SETUP}>
          <AriaButton color="secondary" size="large">
            <Trans>Set Up</Trans>
          </AriaButton>
        </Link>
      </div>
    </section>
  );
}

export default SetupYourTenant;

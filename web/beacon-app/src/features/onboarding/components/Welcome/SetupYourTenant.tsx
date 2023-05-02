import { Button, Card, Heading } from '@rotational/beacon-core';
import { Link } from 'react-router-dom';

import { ROUTES } from '@/application';
import HostedDataIcon from '@/components/icons/hosted-data-icon';

function SetupYourTenant() {
  return (
    <Card className="mx-auto max-w-5xl py-6 text-center text-xl">
      <div className="grid-cols-3  md:grid">
        <HostedDataIcon className="mx-auto mt-6" />
        <div>
          <Heading as="h2">
            Set Up Your Tenant <span className="font-normal">(required)</span>
          </Heading>
          <p className="mt-6">
            Your tenant is your team&apos;s control panel for all projects and topics. Specify
            preferences around encryption, privacy, and locality (e.g. for GDPR, CCPA, etc).
          </p>
        </div>
        <div className="m-10 mx-auto md:mt-16">
          <Link to={ROUTES.SETUP}>
            <Button variant="secondary" size="large">
              Set Up
            </Button>
          </Link>
        </div>
      </div>
    </Card>
  );
}

export default SetupYourTenant;

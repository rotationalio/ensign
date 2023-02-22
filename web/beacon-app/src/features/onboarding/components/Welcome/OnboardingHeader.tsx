import { Trans } from '@lingui/macro';
import { memo } from 'react';

function OnboardingHeader() {
  return (
    <main className="mx-auto mt-16 max-w-4xl">
      <h1 className="text-center text-5xl font-bold">
        <Trans>Let&apos;s Get Eventing</Trans>
      </h1>
      <div className="text-2xl">
        <p className="mx-auto mt-8">
          <Trans>
            You did it! What&apos;s next? Set up your tenant. Invite team members. Create your first
            project and get streaming. Update your resume to let everyone know you&apos;re a
            microservices expert now.
          </Trans>
        </p>
        <div className="mt-8">
          <p className="font-bold">
            <Trans>With the Starter Plan, you get:</Trans>
          </p>
          <ul className="mt-2 list-inside list-disc">
            <li>
              <Trans>2 default projects</Trans>
            </li>
            <li>
              <Trans>3 topics per project</Trans>
            </li>
            <li>
              <Trans>5 GB of data storage</Trans>
            </li>
          </ul>
        </div>
        <p className="pt-8">
          <Trans>Upgrade any time you&apos;re ready to kick things up a notch.</Trans>
        </p>
      </div>
    </main>
  );
}

export default memo(OnboardingHeader);

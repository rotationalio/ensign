import { memo } from 'react';

function OnboardingHeader() {
  return (
    <main className="mx-auto mt-16 max-w-4xl">
      <h1 className="text-center text-6xl font-bold">Let&apos;s Get Eventing</h1>
      <div className="text-2xl">
        <p className="mx-auto mt-8">
          You did it! What&apos;s next? Set up your tenant. Invite team members. Create your first
          project and get streaming. Update your resume to let everyone know you&apos;re a
          microservices expert now.
        </p>
        <div className="mt-8">
          <p>With the Starter Plan, you get:</p>
          <ul className="mt-2 list-inside list-disc">
            <li>3 projects</li>
            <li>1 topic per project</li>
            <li>10 GB of data storage per project</li>
          </ul>
        </div>
        <p className="pt-8">Upgrade any time you&apos;re ready to kick things up a notch.</p>
      </div>
    </main>
  );
}

export default memo(OnboardingHeader);

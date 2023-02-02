import { memo } from 'react'

function OnboardingHeader() {
  return (
    <main className="mx-auto max-w-4xl mt-16">
        <h1 className="text-6xl text-center font-bold">Let's Get Eventing</h1>
        <div className="text-2xl">
          <p className="mt-8 mx-auto">
            You did it! What's next? Set up your tenant. Invite team members. Create your first project and get streaming. Update your resume to let everyone know you're a microservices expert now.
          </p>
          <div className="mt-8">
            <p>With the Starter Plan, you get:</p>
            <ul className="mt-2 list-disc list-inside">
              <li>3 projects</li>
              <li>1 topic per project</li>
              <li>10 GB of data storage per project</li>
            </ul>
          </div>
          <p className="pt-8">Upgrade any time you're ready to kick things up a notch.</p>
        </div>
      </main>
  );
}

export default memo(OnboardingHeader);

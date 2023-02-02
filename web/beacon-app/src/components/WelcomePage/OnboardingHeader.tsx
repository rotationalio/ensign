import { memo } from 'react'

import SetupTenant from './SetupTenant';
import ViewTutorials from './ViewTutorials';

function OnboardingHeader() {
  return (
    <main className="mx-auto max-w-4xl">
        <h1 className="text-6xl text-center font-bold">Let's Get Eventing</h1>
        <p className="mt-8 text-2xl mx-auto">
          You did it! What's next? Set up your tenant. Invite team members. Create your first project and get streaming. Update your resume to let everyone know you're a microsrevices expert now.
        </p>
        <div className="mt-8">
          <p>With the Starter Plan, you get:</p>
          <ul className="mt-2 list-disc list-inside">
            <li>3 projects</li>
            <li>1 topic per project</li>
            <li>10 GB of data storage per project</li>
          </ul>
        </div>
        <p className="pt-4">Upgrade any time you're ready to kick things up a notch.</p>
    </main>
  );
}

export default memo(OnboardingHeader);

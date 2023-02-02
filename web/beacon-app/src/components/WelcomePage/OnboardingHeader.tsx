import { memo } from 'react'

import SetupTenant from './SetupTenant';
import ViewTutorials from './ViewTutorials';

function OnboardingHeader() {
  return (
    <main>
      <section>
        <h1 className="text-6xl font-bold">Let's Get Eventing</h1>
        <p className="mt-4 text-2xl">
          You did it! What's next? Set up your tenant. Invite team members. Create your first project and get streaming. Update your resume to let everyone know you're a microsrevices expert now.
        </p>
        <div>
          <p>With the Starter Plan, you get:</p>
          <ul>
            <li>3 projects</li>
            <li>1 topic per project</li>
            <li>10 GB of data storage per project</li>
          </ul>
        </div>
        <p>Upgrade any time you're ready to kick things up a notch.</p>
      </section>
      <div className="mt-6">
        <SetupTenant />
      </div>
      <div className="mt-6">
        <ViewTutorials />
      </div>
    </main>
  );
}

export default memo(OnboardingHeader);

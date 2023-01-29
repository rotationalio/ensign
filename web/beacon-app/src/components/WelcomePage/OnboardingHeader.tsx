import { memo } from 'react'

import SetupTenant from './SetupTenant';
import ViewTutorials from './ViewTutorials';

function OnboardingHeader() {
  return (
    <main>
      <section>
        <h1 className="text-6xl font-bold">Welcome to Ensign</h1>
        <p className="mt-4 text-2xl">
          Get started now! Setup your tenant. Then create your first project and topic to start
          streaming and delivering in real-time.
        </p>
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

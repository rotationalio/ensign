import AccessDashboard from './AccessDashboard';
import SetupTenant from './SetupTenant';
import ViewTutorials from './ViewTutorials';

function OnboardingHeader() {
  return (
    <main>
      <section>
        <h1>Welcome to Ensign</h1>
        <p>
          Get started now! Setup your tenant. Then create your first project and topic to start
          streaming and delivering in real-time.
        </p>
      </section>
      <SetupTenant />
      <ViewTutorials />
      <AccessDashboard />
    </main>
  );
}

export default OnboardingHeader;

import OnboardingHeader from './OnboardingHeader';
import SetupYourTenant from './SetupYourTenant';
import ViewTutorials from './ViewTutorials';

export default function WelcomePage() {
  return (
    <div>
      <OnboardingHeader />
      <div className="mt-6">
        <SetupYourTenant />
        <div className="mt-12">
          <ViewTutorials />
        </div>
      </div>
    </div>
  );
}

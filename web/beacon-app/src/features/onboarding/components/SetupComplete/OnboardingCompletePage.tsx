import OnboardingHeader from '../Welcome/OnboardingHeader';
import ViewTutorials from '../Welcome/ViewTutorials';
import SetupTenantComplete from './SetupTenantComplete';

export default function OnboardingCompletePage() {
  return (
    <div>
      <OnboardingHeader />
      <div className="mt-6">
        <SetupTenantComplete />
        <div className="mt-12">
          <ViewTutorials />
        </div>
      </div>
    </div>
  );
}

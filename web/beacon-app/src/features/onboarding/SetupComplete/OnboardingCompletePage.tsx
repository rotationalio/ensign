import OnboardingHeader from '../components/Welcome/OnboardingHeader';
import ViewTutorials from '../components/Welcome/ViewTutorials';
import SetupTenantComplete from './SetupTenantComplete';

export default function OnboardingCompletePage() {
  return (
    <div className="bg-hexagon bg-contain">
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

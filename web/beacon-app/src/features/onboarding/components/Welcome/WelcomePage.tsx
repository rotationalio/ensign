import OnboardingHeader from "./OnboardingHeader";
import SetupYourTenant from "./SetupYourTenant";
import ViewTutorials from "./ViewTutorials";
import { t } from '@lingui/macro';

export default function WelcomePage() {
    return(
    <div>
      <OnboardingHeader />
      <div className="mt-6">
        <SetupYourTenant />
        <div className="mt-12">
          <ViewTutorials />
        </div>
      </div>
    </div>
    )
}
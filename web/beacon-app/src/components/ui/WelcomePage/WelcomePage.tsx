import { LandingFooter } from "../../auth/LandingFooter";
import { LandingHeader } from "../../auth/LandingHeader";
import OnboardingHeader from "./OnboardingHeader";
import SetupYourTenant from "./SetupYourTenant";
import ViewTutorials from "./ViewTutorials";
import { t } from '@lingui/macro';

export default function WelcomePage() {
    return(
    <div>
      <LandingHeader />
      <OnboardingHeader />
      <div className="mt-6">
        <SetupYourTenant />
        <div className="mt-12">
          <ViewTutorials />
        </div>
      </div>
      <LandingFooter />
    </div>
    )
}
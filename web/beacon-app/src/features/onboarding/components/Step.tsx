import { useOrgStore } from '@/store';

import OnboardingFormLayout from '../layout';
import { ONBOARDING_STEPS } from '../shared/constants';
import { OrganizationStep, UserNameStep, UserPreferenceStep, WorkspaceStep } from './steps';
const Step = () => {
  const onboarding = useOrgStore((state: any) => state.onboarding) as any;
  const { currentStep } = onboarding as any;
  let stepContent;
  switch (currentStep) {
    case ONBOARDING_STEPS.ORGANIZATION:
      stepContent = <OrganizationStep />;
      break;
    case ONBOARDING_STEPS.WORKSPACE:
      stepContent = <WorkspaceStep />;
      break;
    case ONBOARDING_STEPS.NAME:
      stepContent = <UserNameStep />;
      break;
    case ONBOARDING_STEPS.PREFERENCE:
      stepContent = <UserPreferenceStep />;
      break;
    default:
      stepContent = <OrganizationStep />;
      break;
  }

  return <OnboardingFormLayout>{stepContent}</OnboardingFormLayout>;
};

export default Step;

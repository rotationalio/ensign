import { useOrgStore } from '@/store';

import { OrganizationStep, UserNameStep, UserPreferenceStep, WorkspaceStep } from './steps';

const Step = () => {
  const orgDataState = useOrgStore.getState() as any;
  const { currentStep } = orgDataState?.onboarding || null;
  let stepContent;
  switch (currentStep) {
    case 0:
      stepContent = <OrganizationStep />;
      break;
    case 1:
      stepContent = <WorkspaceStep />;
      break;
    case 2:
      stepContent = <UserNameStep />;
      break;
    case 3:
      stepContent = <UserPreferenceStep />;
      break;
    default:
      stepContent = <OrganizationStep />;
      break;
  }

  return stepContent;
};

export default Step;

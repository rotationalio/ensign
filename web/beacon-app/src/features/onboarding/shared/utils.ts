import { t } from '@lingui/macro';

import { MemberResponse } from '@/features/members/types/memberServices';

import { ONBOARDING_STEPS } from './constants';
export const getDeveloperOptions = () => {
  return [
    { value: 'Application development', label: t`Application development` },
    { value: 'Data science', label: t`Data science` },
    { value: 'Data engineering', label: t`Data engineering` },
    { value: 'Developer experience', label: t`Developer experience` },
    { value: 'Cybersecurity (blue or purple team)', label: t`Cybersecurity (blue or purple team)` },
    { value: 'DevOps and observability', label: t`DevOps and observability` },

    { value: 'Something else', label: 'Something else' },
  ];
};

export const getProfessionOptions = () => [
  {
    id: 'work_segment',
    value: 'work',
    label: t`Work`,
  },
  {
    id: 'education_segment',
    value: 'education',
    label: t`Education`,
  },
  {
    id: 'personal_segment',
    value: 'personal',
    label: t`Personal`,
  },
];

export const getCurrentStepFromMember = (member: any) => {
  let current = 1;
  const hasOrganization = member?.organization?.length > 0;
  const hasWorkspace = member?.workspace?.length > 0;
  const hasName = member?.name?.length > 0;
  const hasProfessionSegment = member?.profession_segment?.length > 0;
  const hasDeveloperSegment = member?.developer_segment?.length > 0;
  const hasReachStep = (step: number) => {
    switch (step) {
      case 1:
        return !hasOrganization;
      case 2:
        return hasOrganization && !hasWorkspace;
      case 3:
        return hasOrganization && hasWorkspace && !hasName;
      case 4:
        return hasOrganization && hasWorkspace && hasName;
      default:
        return false;
    }
  };

  if (hasReachStep(ONBOARDING_STEPS.ORGANIZATION)) {
    current = 1;
  }

  if (hasReachStep(ONBOARDING_STEPS.WORKSPACE)) {
    current = 2;
  }

  if (hasReachStep(ONBOARDING_STEPS.NAME)) {
    current = 3;
  }

  if (hasReachStep(ONBOARDING_STEPS.PREFERENCE)) {
    current = 4;
  }

  if (
    !hasOrganization &&
    !hasWorkspace &&
    !hasName &&
    (!hasProfessionSegment || !hasDeveloperSegment)
  ) {
    current = 1;
  }

  return current;
};

export const getOnboardingStepsData = (member: Partial<MemberResponse>) => {
  const { name, organization, workspace, profession_segment, developer_segment } = member;

  return {
    name,
    organization,
    workspace,
    developer_segment,
    profession_segment,
  };
};

export const hasCompletedOnboarding = (member: MemberResponse) => {
  const { name, organization, workspace, profession_segment, developer_segment } = member;

  return name && organization && workspace && profession_segment.length > 0 && developer_segment;
};

export const isInvitedUser = (member: Pick<MemberResponse, 'invited'>) => {
  return member?.invited === true;
};

export const stepperContents = [
  {
    title: t`Step 1 of 4`,
    description: t`Your Team Name`,
  },
  {
    title: t`Step 2 of 4`,
    description: t`Your Workspace URL`,
  },
  {
    title: t`Step 3 of 4`,
    description: t`Your Name`,
  },
  {
    title: t`Step 4 of 4`,
    description: t`Your Goals`,
  },
];

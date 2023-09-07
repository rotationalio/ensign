import { t } from '@lingui/macro';

import { MemberResponse } from '@/features/members/types/memberServices';
export const getDeveloperOptions = () => {
  return [
    { value: 'Application development', label: t`Application development` },
    { value: 'Data science', label: t`Data science` },
    { value: 'Data engineering', label: t`Data engineering` },
    { value: 'Developer experience', label: t`Developer experience` },
    { value: 'Cybersecurity', label: t`Cybersecurity (blue or purple team)` },
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

  if (hasOrganization) {
    current = 2;
  }
  if (hasOrganization && hasWorkspace) {
    current = 3;
  }
  if (hasOrganization && hasWorkspace && hasName) {
    current = 4;
  }
  if (hasDeveloperSegment && hasProfessionSegment) {
    current = 4;
  }
  if (!hasOrganization && !hasWorkspace && !hasName && !hasProfessionSegment) {
    current = 1;
  }

  return current;
};

export const getOnboardingStepsData = (member: Partial<MemberResponse>) => {
  const { name, organization, workspace, profession_segment } = member;

  return {
    name,
    organization,
    workspace,
    profession_segment,
  };
};

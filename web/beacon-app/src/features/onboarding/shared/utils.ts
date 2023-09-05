import { t } from '@lingui/macro';
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
  let current = 0;
  const hasOrganization = member?.organization?.length > 0;
  const hasWorkspace = member?.workspace?.length > 0;
  const hasName = member?.name?.length > 0;
  const hasProfessionSegment = member?.profession_segment?.length > 0;

  if (hasOrganization) {
    current = 2;
  }
  if (hasWorkspace) {
    current = 3;
  }
  if (hasName) {
    current = 4;
  }
  if (hasProfessionSegment) {
    current = 5; // this step 5 doesnt exist so that will trigger the onboarding complete modal
  }
  if (!hasOrganization && !hasWorkspace && !hasName && !hasProfessionSegment) {
    current = 1;
  }

  return current;
};

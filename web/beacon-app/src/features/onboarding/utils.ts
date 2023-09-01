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

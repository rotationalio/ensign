import { t } from '@lingui/macro';

export const getDefaultTopicStats = () => {
  return [
    {
      name: t`Online Publishers`,
      value: 0,
      units: '',
    },
    {
      name: t`Online Subscribers`,
      value: 0,
      units: '',
    },
    {
      name: t`Avg Events/ Second`,
      value: 0,
      units: 'eps',
    },
    {
      name: t`Data Storage`,
      value: 0,
      units: 'GB',
    },
  ];
};

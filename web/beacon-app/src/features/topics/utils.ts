import { t } from '@lingui/macro';
import DOMPurify from 'dompurify';

import { formatDate } from '@/utils/formatDate';

import type { Topic } from '../topics/types/topicService';

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
      name: t`Total Events`,
      value: 0,
    },
    {
      name: t`Data Storage`,
      value: 0,
      units: 'GB',
    },
  ];
};

export const getTopicStatsHeaders = () => {
  return [t`Online Publishers`, t`Online Subscribers`, t`Total Events`, t`Data Storage`];
};

export const getFormattedTopicData = (topic: Topic) => {
  return [
    {
      label: t`Topic ID`,
      value: topic?.id,
    },
    {
      label: t`Status`,
      value: topic?.status,
    },
    {
      label: t`Created`,
      value: formatDate(new Date(topic?.created as string)),
    },
    {
      label: t`Modified`,
      value: formatDate(new Date(topic?.modified as string)),
    },
  ];
};

// this abstraction will sanitize the topic query params

export const inputSanitizer = (input: string) => {
  //  prevent XSS attacks
  const sanitizedInput = DOMPurify.sanitize(input);
  const sanitizedSqlInjection = sanitizedInput.replace(/'/g, "\\'");
  const jsInjectionSafeInput = sanitizedSqlInjection.replace(/</g, '&lt;').replace(/>/g, '&gt;');
  const finalSanitizedInput = jsInjectionSafeInput.trim();

  return finalSanitizedInput;
};

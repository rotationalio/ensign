import { Trans } from '@lingui/macro';

import { Tag } from '@/components/ui/Tag';
import { TOPIC_STATE } from '@/constants/rolesAndStatus';

interface TopicStateTagProps {
  status: string;
}

const topicStateMap = {
  [TOPIC_STATE.ACTIVE]: 'success',
  [TOPIC_STATE.PENDING]: 'secondary',
  [TOPIC_STATE.ARCHIVED]: 'warning',
  [TOPIC_STATE.DELETTING]: 'error',
} as const;

const TopicStateTag = ({ status }: TopicStateTagProps) => {
  return (
    <Tag size="medium" variant={topicStateMap[status]} data-cy="topic-status-tag">
      <Trans>{status}</Trans>
    </Tag>
  );
};

export default TopicStateTag;

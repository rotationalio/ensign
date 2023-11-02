import { Trans } from '@lingui/macro';

import { Tag } from '@/components/ui/Tag';
import { TOPIC_STATE } from '@/constants/rolesAndStatus';

interface TopicStateTagProps {
  status: string;
}

const topicStateMap = {
  [TOPIC_STATE.ACTIVE]: 'success',
  [TOPIC_STATE.ARCHIVED]: 'warning',
  [TOPIC_STATE.DELETING]: 'error',
} as const;

const TopicStateTag = ({ status }: TopicStateTagProps) => {
  return (
    <Tag size="medium" variant={topicStateMap[status]}>
      <span data-cy="topic-status-tag">
        <Trans>{status}</Trans>
      </span>
    </Tag>
  );
};

export default TopicStateTag;

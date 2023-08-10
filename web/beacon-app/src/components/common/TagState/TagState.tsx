import { Trans } from '@lingui/macro';

import { Tag } from '@/components/ui/Tag';
import { TOPIC_STATE } from '@/constants/rolesAndStatus';

interface TopicStateTagProps {
  status: string;
  colorScheme?: string;
}

// TODO: list all possible values of status
const StateMap = {
  [TOPIC_STATE.ACTIVE]: 'success',
  [TOPIC_STATE.PENDING]: 'secondary',
  [TOPIC_STATE.ARCHIVED]: 'warning',
  [TOPIC_STATE.DELETTING]: 'error',
} as const;

const TagState = ({ status }: TopicStateTagProps) => {
  return (
    <Tag size="medium" variant={StateMap[status]}>
      <Trans>{status}</Trans>
    </Tag>
  );
};

export default TagState;

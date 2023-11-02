import { Trans } from '@lingui/macro';

import { Tag } from '@/components/ui/Tag';
import { PROJECT_STATE } from '@/constants/rolesAndStatus';

interface TopicStateTagProps {
  status: string;
  colorScheme?: string;
}

// TODO: list all possible values of status
const StateMap = {
  [PROJECT_STATE.ACTIVE]: 'success',
  [PROJECT_STATE.ARCHIVED]: 'error',
  [PROJECT_STATE.INCOMPLETE]: 'warning',
} as const;

const TagState = ({ status }: TopicStateTagProps) => {
  return (
    <Tag size="medium" variant={StateMap[status]}>
      <Trans>{status}</Trans>
    </Tag>
  );
};

export default TagState;

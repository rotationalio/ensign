import { Trans } from '@lingui/macro';

import { Tag } from '@/components/ui/Tag';

import { getTopicTagVariant } from '../utils';

interface TopicStateTagProps {
  status: any;
}

const TopicStateTag = ({ status }: TopicStateTagProps) => {
  return (
    <Tag size="medium" variant={getTopicTagVariant(status)}>
      <Trans>{status}</Trans>
    </Tag>
  );
};

export default TopicStateTag;

import { Button } from '@rotational/beacon-core';

import SettingIcon from '@/components/icons/setting';

const TopicSettings = () => {
  return (
    <>
      <Button
        variant="ghost"
        size="custom"
        className="flex-end bg-inherit hover:bg-transparent border-none"
        data-cy="topicDetailActions"
      >
        <SettingIcon />
      </Button>
    </>
  );
};

export default TopicSettings;

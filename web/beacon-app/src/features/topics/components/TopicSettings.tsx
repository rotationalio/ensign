import { Trans } from '@lingui/macro';
import { Button, Menu, useMenu } from '@rotational/beacon-core';

import SettingIcon from '@/components/icons/setting';

const TopicSettings = () => {
  const { isOpen, close, open, anchorEl } = useMenu({ id: 'topic-menu-action' });
  return (
    <>
      <Button
        variant="ghost"
        size="custom"
        className="flex-end bg-inherit hover:bg-transparent border-none"
        onClick={open}
        data-cy="topicDetailActions"
      >
        <SettingIcon />
      </Button>
      <Menu open={isOpen} onClose={close} anchorEl={anchorEl}>
        <Menu.Item>
          <Trans>Archive Topic</Trans>
        </Menu.Item>
        <Menu.Item>
          <Trans>Delete Topic</Trans>
        </Menu.Item>
        <Menu.Item>
          <Trans>Clone Topic</Trans>
        </Menu.Item>
      </Menu>
    </>
  );
};

export default TopicSettings;

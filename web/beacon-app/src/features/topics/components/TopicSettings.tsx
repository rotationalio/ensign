import { Trans } from '@lingui/macro';
import { Button, Menu, useMenu } from '@rotational/beacon-core';
import { useState } from 'react';

import SettingIcon from '@/components/icons/setting';

import ArchiveTopicModal from './ArchiveTopicModal';
import DeleteTopicModal from './DeleteTopicModal';

const TopicSettings = () => {
  const { isOpen, close, open, anchorEl } = useMenu({ id: 'topic-menu-action' });
  const [isArchiveModalOpen, setIsArchiveModalOpen] = useState<boolean>(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState<boolean>(false);

  const openArchiveModal = () => {
    setIsArchiveModalOpen(true);
  };

  const onCloseArchiveModal = () => {
    setIsArchiveModalOpen(false);
  };

  const openDeleteModal = () => {
    setIsDeleteModalOpen(true);
  };

  const onCloseDeleteModal = () => {
    setIsDeleteModalOpen(false);
  };

  return (
    <>
      <div>
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
          <Menu.Item onClick={openArchiveModal}>
            <Trans>Archive Topic</Trans>
          </Menu.Item>
          <Menu.Item onClick={openDeleteModal}>
            <Trans>Delete Topic</Trans>
          </Menu.Item>
          <Menu.Item>
            <Trans>Clone Topic</Trans>
          </Menu.Item>
        </Menu>
      </div>
      <ArchiveTopicModal isOpen={isArchiveModalOpen} close={onCloseArchiveModal} />
      <DeleteTopicModal isOpen={isDeleteModalOpen} close={onCloseDeleteModal} />
    </>
  );
};

export default TopicSettings;

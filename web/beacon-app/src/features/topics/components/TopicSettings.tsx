import { Trans } from '@lingui/macro';
import { Button, Menu, useMenu } from '@rotational/beacon-core';
import { useState } from 'react';

import SettingIcon from '@/components/icons/setting';

import ArchiveTopicModal from './Modal/ArchiveTopicModal';
import CloneTopicModal from './Modal/CloneTopicModal';
import DeleteTopicModal from './Modal/DeleteTopicModal';

const TopicSettings = () => {
  const { isOpen, close, open, anchorEl } = useMenu({ id: 'topic-menu-action' });
  const [isArchiveTopicModalOpen, setIsArchiveTopicModalOpen] = useState<boolean>(false);
  const [isDeleteTopicModalOpen, setIsDeleteTopicModalOpen] = useState<boolean>(false);
  const [isCloneTopicModalOpen, setIsCloneTopicModalOpen] = useState<boolean>(false);

  const openArchiveTopicModal = () => {
    setIsArchiveTopicModalOpen(true);
  };

  const onCloseArchiveTopicModal = () => {
    setIsArchiveTopicModalOpen(false);
  };

  const openDeleteTopicModal = () => {
    setIsDeleteTopicModalOpen(true);
  };

  const onCloseDeleteTopicModal = () => {
    setIsDeleteTopicModalOpen(false);
  };

  const openCloneTopicModal = () => {
    setIsCloneTopicModalOpen(true);
  };

  const onCloseCloneTopicModal = () => {
    setIsCloneTopicModalOpen(false);
  };

  return (
    <>
      <div>
        <Button
          variant="ghost"
          size="custom"
          className="flex-end bg-inherit hover:bg-transparent border-none"
          onClick={open}
          data-cy="topic-detail-actions"
        >
          <SettingIcon />
        </Button>
        <Menu open={isOpen} onClose={close} anchorEl={anchorEl}>
          <Menu.Item onClick={openArchiveTopicModal}>
            <Trans>Archive Topic</Trans>
          </Menu.Item>
          <Menu.Item onClick={openDeleteTopicModal}>
            <Trans>Delete Topic</Trans>
          </Menu.Item>
          <Menu.Item onClick={openCloneTopicModal}>
            <Trans>Clone Topic</Trans>
          </Menu.Item>
        </Menu>
      </div>
      <ArchiveTopicModal isOpen={isArchiveTopicModalOpen} close={onCloseArchiveTopicModal} />
      <DeleteTopicModal isOpen={isDeleteTopicModalOpen} close={onCloseDeleteTopicModal} />
      <CloneTopicModal isOpen={isCloneTopicModalOpen} close={onCloseCloneTopicModal} />
    </>
  );
};

export default TopicSettings;

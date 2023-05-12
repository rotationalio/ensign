import { Trans } from '@lingui/macro';
import { Button, Menu, useMenu } from '@rotational/beacon-core';
import { useState } from 'react';

import SettingIcon from '@/components/icons/setting';

import DeleteProjectModal from './DeleteProjectModal';
import RenameProjectModal from './RenameProject/RenameProjectModal';

interface ProjectSettingsProps {
  data: any;
}
const ProjectSettings = ({ data }: ProjectSettingsProps) => {
  const { isOpen, close, open, anchorEl } = useMenu({ id: 'prj-menu-action' });
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [isRenameModalOpen, setIsRenameModalOpen] = useState(false);

  const openRenameModal = () => {
    setIsRenameModalOpen(true);
  };

  const onCloseRenameModal = () => {
    setIsRenameModalOpen(false);
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
          data-cy="detailActions"
        >
          <SettingIcon />
        </Button>
        <Menu open={isOpen} onClose={close} anchorEl={anchorEl}>
          <Menu.Item onClick={openDeleteModal} data-testid="cancelButton">
            <Trans>Delete Project</Trans>
          </Menu.Item>
          <Menu.Item onClick={openRenameModal} data-testid="rename-project">
            <Trans>Rename Project</Trans>
          </Menu.Item>
        </Menu>
      </div>
      <DeleteProjectModal isOpen={isDeleteModalOpen} close={onCloseDeleteModal} />
      <RenameProjectModal
        open={isRenameModalOpen}
        handleModalClose={onCloseRenameModal}
        project={data}
      />
    </>
  );
};

export default ProjectSettings;

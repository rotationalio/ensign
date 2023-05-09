import { Trans } from '@lingui/macro';
import { Button, Menu, useMenu } from '@rotational/beacon-core';
import { useState } from 'react';

import SettingIcon from '@/components/icons/setting';

import DeleteProjectModal from './DeleteProjectModal';

const ProjectSettings = () => {
  const { isOpen, close, open, anchorEl } = useMenu({ id: 'prj-menu-action' });
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);

  const openDeleteModal = () => {
    setIsDeleteModalOpen(true);
  };

  const onCloseDeleteModal = () => {
    setIsDeleteModalOpen(false);
  };

  return (
    <>
      <div>
        <Button variant="ghost" className="flex-end bg-inherit w-8 border-none" onClick={open}>
          <div className="ml-[4.5rem]">
            <SettingIcon />
          </div>
        </Button>
        <Menu open={isOpen} onClose={close} anchorEl={anchorEl}>
          <Menu.Item onClick={openDeleteModal} data-testid="cancelButton">
            <Trans>Delete Project</Trans>
          </Menu.Item>
        </Menu>
      </div>
      <DeleteProjectModal isOpen={isDeleteModalOpen} close={onCloseDeleteModal} />
    </>
  );
};

export default ProjectSettings;

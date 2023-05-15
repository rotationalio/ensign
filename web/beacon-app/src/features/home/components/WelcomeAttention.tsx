import { Trans } from '@lingui/macro';
import { Button } from '@rotational/beacon-core';
import React, { useState } from 'react';

import NewProjectModal from '@/features/projects/components/NewProject/NewProjectModal';

const WelcomeAttention = () => {
  const [isOpenNewProjectModal, setIsOpenNewProjectModal] = useState<boolean>(false);

  const onOpenNewProjectModal = () => {
    setIsOpenNewProjectModal(true);
  };

  const onCloseNewProjectModal = () => {
    setIsOpenNewProjectModal(false);
  };
  return (
    <>
      <div
        className="px-auto mb-8 mt-4 flex flex-row items-center justify-between space-x-4 rounded-md border border-neutral-500 bg-[#F7F9FB] p-2 px-5 text-justify"
        data-cy="projWelcome"
      >
        <p className="text-md">
          <Trans>
            Welcome to Ensign! Get started on your first project. We’ll guide you along the way!
          </Trans>
        </p>

        <Button
          //variant="ghost"
          size="small"
          className="!bg-green text-white hover:!bg-green/[0.8]"
          onClick={onOpenNewProjectModal}
          data-cy="startSetupBttn"
        >
          <Trans>Start</Trans>
        </Button>
      </div>
      <NewProjectModal isOpened={isOpenNewProjectModal} onClose={onCloseNewProjectModal} />
    </>
  );
};

export default WelcomeAttention;

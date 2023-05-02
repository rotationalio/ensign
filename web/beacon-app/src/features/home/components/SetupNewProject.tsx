import { t, Trans } from '@lingui/macro';
import { Button } from '@rotational/beacon-core';
import React, { useState } from 'react';

import { CardListItem } from '@/components/common/CardListItem';
import NewProjectModal from '@/features/projects/components/NewProject/NewProjectModal';
function SetupNewProject() {
  const [isOpenNewProjectModal, setIsOpenNewProjectModal] = useState<boolean>(false);

  const onOpenNewProjectModal = () => {
    setIsOpenNewProjectModal(true);
  };

  const onCloseNewProjectModal = () => {
    setIsOpenNewProjectModal(false);
  };

  return (
    <>
      <CardListItem
        title={t`Set Up A New Project`}
        titleClassName="text-lg"
        className="min-h-[130px]"
        contentClassName="my-2"
      >
        <div className="mt-2 flex flex-col gap-8 px-3 xl:flex-row">
          <p className="text-md  w-full sm:w-4/5">
            <Trans>
              Set up a project to customize data flows. A project is a collection of topics. Topics
              are event streams that your services, applications, or models can publish or subscribe
              to for real-time data flows. Control access by generating API keys.
            </Trans>
          </p>
          <div className="sm:w-1/5">
            <Button
              className="text-md"
              size="small"
              onClick={onOpenNewProjectModal}
              data-testid="set-new-project"
            >
              <Trans>Create</Trans>
            </Button>
          </div>

          <NewProjectModal isOpened={isOpenNewProjectModal} onClose={onCloseNewProjectModal} />
        </div>
      </CardListItem>
    </>
  );
}

export default SetupNewProject;

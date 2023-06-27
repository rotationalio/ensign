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
        <div className="mt-2 flex flex-col justify-between gap-4 px-3 xl:flex-row ">
          <p className="text-md  w-full sm:w-4/5">
            <Trans>
              A project is <span className="font-bold"> a database for events </span> — a collection
              of datasets related by use case. However, it stores all updates to each object over
              time, so you can observe changes and activity in your data feeds, applications, and
              models. Use SDKs to connect sources to publish data to your project or subscribe to
              read updates in real-time. Control access by generating API keys.
            </Trans>
          </p>
          <div className="item-center place-items-center">
            <Button size="medium" onClick={onOpenNewProjectModal} data-testid="set-new-project">
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

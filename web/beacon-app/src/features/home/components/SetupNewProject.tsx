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
        <div
          className="mt-2 flex flex-col justify-between gap-4 px-3 xl:flex-row"
          data-cy="setup-new-project"
        >
          <p className="text-md  w-full sm:w-4/5">
            <Trans>
              A project is <span className="font-bold">your use case</span> for real-time data
              management. A project contains one or many <span className="font-bold">topics</span>.
              Topics are immutable data stores that capture change over time. Control access to your
              projects with API keys. Using SDKs, set up your data sources to publish data to your
              topics in real-time. Then, set up downstream subscribers to read updates from your
              topics in real-time.
            </Trans>
          </p>
          <div className="item-center place-items-center">
            <Button
              size="medium"
              onClick={onOpenNewProjectModal}
              data-testid="set-new-project"
              data-cy="create-project-btn"
            >
              <Trans>Create Project</Trans>
            </Button>
          </div>

          <NewProjectModal isOpened={isOpenNewProjectModal} onClose={onCloseNewProjectModal} />
        </div>
      </CardListItem>
    </>
  );
}

export default SetupNewProject;

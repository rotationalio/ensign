import { t, Trans } from '@lingui/macro';
import { Button, Loader, Table } from '@rotational/beacon-core';
import { Suspense, useEffect, useMemo, useState } from 'react';
import toast from 'react-hot-toast';
import { useNavigate, useParams } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application';
import { useFetchTopics } from '@/features/topics/hooks/useFetchTopics';

import { getHiddenColumns, getTopicInitialColumns, getTopics } from '../util';
import { NewTopicModal } from './NewTopicModal';
import TopicTableHeader from './TopicTableHeader';
export const TopicTable = () => {
  const navigate = useNavigate();
  const initialColumns = useMemo(() => getTopicInitialColumns(), []) as any;

  const [openNewTopicModal, setOpenNewTopicModal] = useState(false);
  const handleOpenNewTopicModal = () => setOpenNewTopicModal(true);
  const handleCloseNewTopicModal = () => setOpenNewTopicModal(false);

  const param = useParams<{ id: string }>();
  const { id: projectID } = param;
  const projID = projectID || (projectID as string);

  const { topics, isFetchingTopics, hasTopicsFailed, error } = useFetchTopics(projID);

  const redirectToTopicDetails = (topicID: string) => {
    if (!topicID) {
      toast.error(t`Topic ID is missing`);
    }
    navigate(`${PATH_DASHBOARD.TOPICS}/${topicID}`);
  };

  useEffect(() => {
    if (error && hasTopicsFailed) {
      toast.error(
        (error as any)?.response?.data?.error ||
          t`Sorry we are having trouble fetching your topics, please try again later.`
      );
    }
  }, [error, hasTopicsFailed]);

  return (
    <Suspense
      fallback={
        <div className="flex justify-center">
          <Loader />
        </div>
      }
    >
      <div className="mt-[26px]" data-cy="topicComp">
        {isFetchingTopics ? (
          <Loader />
        ) : (
          <>
            <TopicTableHeader />
            <div className="flex w-full justify-between bg-[#F7F9FB] p-2">
              <div className="flex items-center gap-3"></div>
              <Button
                variant="primary"
                size="small"
                className="!text-xs"
                onClick={handleOpenNewTopicModal}
                data-cy="addTopic"
              >
                <Trans>+ New Topic</Trans>
              </Button>
              <NewTopicModal open={openNewTopicModal} handleClose={handleCloseNewTopicModal} />
            </div>
            <div className="overflow-hidden text-sm">
              <Table
                trClassName="text-sm"
                columns={initialColumns}
                data={getTopics(topics)}
                onRowClick={(row: any) => {
                  redirectToTopicDetails(row?.values?.id);
                }}
                initialState={getHiddenColumns(['id'])}
                data-cy="topicTable"
              />
            </div>
          </>
        )}
      </div>
    </Suspense>
  );
};

export default TopicTable;

import { Trans } from '@lingui/macro';
import { Heading, mergeClassnames } from '@rotational/beacon-core';
import { useState } from 'react';
import { useParams } from 'react-router-dom';

import RefreshIcon from '@/components/icons/refresh';

import { useFetchTopicEvents } from '../hooks/useFetchTopicEvents';

const EventDetailTableHeader = () => {
  const param = useParams();
  const { id: topicID } = param as { id: string };
  const { getTopicEvents } = useFetchTopicEvents(topicID);

  const [isRefreshing, setIsRefreshing] = useState(false);

  const refreshHandler = () => {
    setIsRefreshing(true);
    setTimeout(() => {
      getTopicEvents();
      setIsRefreshing(false);
    }, 500);
  };

  return (
    <Heading as="h2" className="mt-8 flex gap-x-2 text-lg font-semibold">
      <Trans>Topic Usage</Trans>
      <button onClick={refreshHandler}>
        <div className={mergeClassnames(isRefreshing ? 'animate-spin-slow' : '')}>
          <RefreshIcon />
        </div>
      </button>
    </Heading>
  );
};

export default EventDetailTableHeader;

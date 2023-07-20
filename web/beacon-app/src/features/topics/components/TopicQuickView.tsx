import { Trans } from '@lingui/macro';
import { Heading, Loader } from '@rotational/beacon-core';
import { Suspense, useEffect, useState } from 'react';

import { QuickView } from '@/components/common/QuickView';
import { SentryErrorBoundary } from '@/components/Error';

import useFetchTopicStats from '../hooks/useFetchTopicStats';
import { getDefaultTopicStats, getTopicStatsHeaders } from '../utils';

interface TopicQuickViewProps {
  topicID: string;
}
const TopicQuickView: React.FC<TopicQuickViewProps> = ({ topicID }) => {
  const { topicStats, error } = useFetchTopicStats(topicID);
  const [topicData, setTopicData] = useState<any>(getDefaultTopicStats()); // by default we will show empty values

  // using useEffect will avoid infinite loop
  useEffect(() => {
    if (topicStats) {
      setTopicData(topicStats);
    }
  }, [topicStats]);

  useEffect(() => {
    if (error) {
      setTopicData(getDefaultTopicStats());
    }
  }, [error]);

  return (
    <Suspense fallback={<Loader />}>
      <SentryErrorBoundary fallback={<div>Something went wrong</div>}>
        <div>
          <Heading as="h1" className="mt-4 text-lg font-semibold">
            <Trans>Quick View</Trans>
          </Heading>
          <QuickView data={topicData} headers={getTopicStatsHeaders()} className="my-4" />
        </div>
      </SentryErrorBoundary>
    </Suspense>
  );
};

export default TopicQuickView;

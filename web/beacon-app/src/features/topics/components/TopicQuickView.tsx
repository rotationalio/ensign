import { Trans } from '@lingui/macro';
import { Heading, Loader } from '@rotational/beacon-core';
import { Suspense } from 'react';

import { QuickView } from '@/components/common/QuickView';
import { SentryErrorBoundary } from '@/components/Error';

import useFetchTopicStats from '../hooks/useFetchTopicStats';
import { getDefaultTopicStats } from '../utils';

interface TopicQuickViewProps {
  topicID: string;
}
const TopicQuickView: React.FC<TopicQuickViewProps> = ({ topicID }) => {
  const { topicStats, error } = useFetchTopicStats(topicID);

  const getTopicStatsData = () => {
    if (!topicStats || error) {
      console.log('getDefaultTopicStats', getDefaultTopicStats());
      return getDefaultTopicStats();
    }
    return topicStats;
  };

  return (
    <Suspense fallback={<Loader />}>
      <SentryErrorBoundary fallback={<div>Something went wrong</div>}>
        <div>
          <Heading as="h1" className="text-lg font-semibold">
            <Trans>Quick View</Trans>
          </Heading>
          <QuickView data={getTopicStatsData()} className="my-4" />
        </div>
      </SentryErrorBoundary>
    </Suspense>
  );
};

export default TopicQuickView;

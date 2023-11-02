import { Trans } from '@lingui/macro';
import { Heading } from '@rotational/beacon-core';
import React from 'react';

// import { useNavigate } from 'react-router-dom';
// import { useNavigate } from 'react-router-dom';
import { HelpTooltip } from '@/components/common/Tooltip/HelpTooltip';
const TopicTableHeader: React.FC = () => {
  // const navigate = useNavigate();

  // const handleRedirection = () => {
  //   navigate(PATH_DASHBOARD.TEMPLATES);
  // };

  return (
    <div>
      <div className="flex items-center">
        <Heading as={'h1'} className="text-lg font-semibold">
          <Trans>Topics</Trans>
        </Heading>
        <span className="ml-2" data-cy="topicHint">
          <HelpTooltip data-cy="topicInfo">
            <p>
              <Trans>
                A topic is a labeled, immutable stream (log) of information that you're interested
                in. Topics hold events in a logical order, giving you the ability to capture data as
                it changes over time and "time travel" back in time for reproducibility,
                explainability, and provenance. Events are sent to and read from your topics.
                Services that are <span className="font-semibold">publishers write data</span> to
                topics. <span className="font-semibold">Subscribers read data</span> from topics.
              </Trans>
            </p>
          </HelpTooltip>
        </span>
      </div>
      <p className="my-4">
        <Trans>Create topics to navigate, ingest, and transform data in real-time.</Trans>
      </p>
    </div>
  );
};

export default React.memo(TopicTableHeader);

import { Trans } from '@lingui/macro';
import { Heading } from '@rotational/beacon-core';
import React from 'react';

import { HelpTooltip } from '@/components/common/Tooltip/HelpTooltip';

const TopicTableHeader: React.FC = () => {
  return (
    <div>
      <Heading as={'h1'} className="flex items-center text-lg font-semibold capitalize">
        <Trans>Topics</Trans>
      </Heading>
      <p className="my-4">
        <Trans>
          You must have at least one topic in your project to publish and subscribe. Topics are
          categories or logs that hold messages and events in a logical order, allowing services and
          data sources to send and receive data between them with ease and accuracy.
        </Trans>
        <span className="ml-2" data-cy="topicHint">
          <HelpTooltip data-cy="topicInfo">
            <p>
              <Trans>
                {' '}
                Messages and events are sent to and read from specific topics. Services that are{' '}
                {''}
                <span className="font-bold">producers, write</span> data to topics. Services that
                are <span className="font-bold">consumers, read</span> data from topics. Topics are
                multi-subscriber, which means that a topic can have zero, one, or multiple consumers
                subscribing to that topic, with read access to the log.
              </Trans>
            </p>
          </HelpTooltip>
        </span>
      </p>
    </div>
  );
};

export default React.memo(TopicTableHeader);

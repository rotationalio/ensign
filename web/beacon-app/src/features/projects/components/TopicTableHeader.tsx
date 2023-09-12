import { Trans } from '@lingui/macro';
import { Heading } from '@rotational/beacon-core';
import React from 'react';

// import { useNavigate } from 'react-router-dom';
// import { useNavigate } from 'react-router-dom';
import { EXTERNAL_LINKS } from '@/application';
import { HelpTooltip } from '@/components/common/Tooltip/HelpTooltip';
import { Link } from '@/components/ui/Link';
const TopicTableHeader: React.FC = () => {
  // const navigate = useNavigate();

  // const handleRedirection = () => {
  //   navigate(PATH_DASHBOARD.TEMPLATES);
  // };

  return (
    <div>
      <Heading as={'h1'} className="flex items-center text-lg font-semibold capitalize">
        <Trans>Design Your Data Flows: Set Up Your Topics</Trans>
      </Heading>
      <p className="my-4">
        <Trans>
          Design your data flows for your use case. Think about where the data is produced and what
          new services, models, or applications benefit from the data. Then create topics or event
          streams, which are logs that hold messages and events in a logical order. As an event
          broker, Ensign navigates the data for you with speed, ease and accuracy. Need help? Watch
          our{' '}
          <Link href={EXTERNAL_LINKS.DATA_FLOW_OVERVIEW} openInNewTab>
            data flow overview,
          </Link>{' '}
          <Link href={EXTERNAL_LINKS.NAMING_TOPICS_GUIDE} openInNewTab>
            read our naming topics guide
          </Link>{' '}
          or{' '}
          <Link href={EXTERNAL_LINKS.OFFICE_HOURS_SCHEDULE} openInNewTab>
            schedule office hours!
          </Link>
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

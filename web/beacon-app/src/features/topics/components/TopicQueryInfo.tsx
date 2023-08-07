import { Trans } from '@lingui/macro';
import { memo } from 'react';
import { Link } from 'react-router-dom';

import { EXTRENAL_LINKS } from '@/application';
const TopicQueryInfo = () => {
  return (
    <div className="flex space-x-1" data-cy="topic-query-instructions">
      <p>
        <Trans>
          Query the topic for insights with{' '}
          <Link
            to={EXTRENAL_LINKS.ENSQL}
            className="font-semibold text-[#1D65A6] underline"
            target="_blank"
            data-cy="topic-query-ensql-link"
          >
            EnSQL
          </Link>{' '}
          e.g. the latest event or last 5 events. The maximum number of query results is 10. Use our{' '}
          <Link
            to={EXTRENAL_LINKS.SDKs}
            className="font-semibold text-[#1D65A6] underline"
            target="_blank"
            data-cy="topic-query-sdks-link"
          >
            SDKs
          </Link>{' '}
          for more results.
        </Trans>
      </p>
    </div>
  );
};

export default memo(TopicQueryInfo);

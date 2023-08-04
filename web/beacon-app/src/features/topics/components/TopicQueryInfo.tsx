import { Trans } from '@lingui/macro';
import { memo } from 'react';
import { Link } from 'react-router-dom';

import { EXTRENAL_LINKS } from '@/application';
const TopicQueryInfo = () => {
  return (
    <div className="flex space-x-1">
      <p>
        <Trans>
          Query the topic for insights with{' '}
          <Link
            to={EXTRENAL_LINKS.ENSQL}
            className="font-semibold text-[#1D65A6] underline"
            target="_blank"
          >
            EnSQL
          </Link>{' '}
          e.g. the latest event or last 5 events. The maximum number of query results is 10. Use our{' '}
          <Link
            to={EXTRENAL_LINKS.SDKs}
            className="font-semibold text-[#1D65A6] underline"
            target="_blank"
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

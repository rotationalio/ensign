import { t, Trans } from '@lingui/macro';
import { mergeClassnames } from '@rotational/beacon-core';
import { FC } from 'react';
interface ResultHeaderProps {
  mimeType: string;
  eventType: string;
  isBase64Encoded: boolean;
}

const ResultHeader: FC<ResultHeaderProps> = ({ mimeType, eventType, isBase64Encoded }) => {
  const isParsable =
    isBase64Encoded && mimeType !== 'application/json' && mimeType !== 'application/xml';

  return (
    <div className="mx-auto flex h-12 flex-row items-center justify-between bg-[#2F4858]/70 p-4">
      <div className="flex flex-row">
        <p className="font-bold text-white">
          <Trans>MIME Type:</Trans>
        </p>
        <p
          className={mergeClassnames('ml-2', isParsable ? 'text-[#FCF77C]' : ' text-white')}
          data-cy="result-mime-type"
        >
          {isParsable ? t`Could not parse. Rendered as base-64 encoded data.` : mimeType ?? 'N/A'}
        </p>
      </div>
      <div className="flex flex-row">
        <p className="font-bold text-white">
          <Trans>Event Type & Version: </Trans>
        </p>
        <p className="ml-2 text-white" data-cy="result-event-type">
          {eventType ?? 'N/A'}
        </p>
      </div>
    </div>
  );
};

export default ResultHeader;

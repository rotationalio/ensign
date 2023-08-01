import { mergeClassnames } from '@rotational/beacon-core';
import { FC } from 'react';

interface ResultHeaderProps {
  mimeType: string;
  eventType: string;
  isBase64Encoded: boolean;
}

const ResultHeader: FC<ResultHeaderProps> = ({ mimeType, eventType, isBase64Encoded }) => {
  return (
    <div className="flex flex-row justify-between bg-primary-900">
      <div className="flex flex-col">
        <div className="flex flex-row">
          <p className="text-white">Mime Type</p>
          <p className={mergeClassnames('ml-2 text-white', isBase64Encoded && 'text-warning-600')}>
            {mimeType}
          </p>
        </div>
        <div className="flex flex-row">
          <p className="text-gray-400">Event Type</p>
          <p className="ml-2 text-gray-400">{eventType}</p>
        </div>
      </div>
    </div>
  );
};

export default ResultHeader;

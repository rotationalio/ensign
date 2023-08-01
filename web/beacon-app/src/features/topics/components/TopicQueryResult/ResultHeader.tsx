import { FC } from 'react';

interface ResultHeaderProps {
  mimeType: string;
  eventType: string;
}

const ResultHeader: FC<ResultHeaderProps> = ({ mimeType, eventType }) => {
  return (
    <div className="flex flex-row justify-between bg-primary-900">
      <div className="flex flex-col">
        <div className="flex flex-row">
          <p className="text-gray-400">Mime Type</p>
          <p className="ml-2 text-gray-400">{mimeType}</p>
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

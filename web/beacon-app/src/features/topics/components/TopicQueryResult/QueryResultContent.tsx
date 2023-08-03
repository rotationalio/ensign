import { t } from '@lingui/macro';
import React from 'react';

import DisplayResultData from './DisplayResultData';
interface QueryResultContentProps {
  result: any;
  mimeType: string;
}

const QueryResultContent: React.FC<QueryResultContentProps> = ({ result, mimeType }) => {
  console.log('[] result', result); // added this avoid eslint error
  console.log('[] mimetype', mimeType);

  return (
    <div className="shadow-md h-20 max-h-80 overflow-y-auto bg-black p-4 text-white">
      <pre className="mx-auto text-sm">
        {result && <DisplayResultData result={result} mimeType={mimeType} />}
        {!result && t`No results found`}
        {!result &&
          !mimeType &&
          t`No query result. Try the default query or enter your own query. See EnSQL documentation for example queries.`}
      </pre>
    </div>
  );
};

export default QueryResultContent;

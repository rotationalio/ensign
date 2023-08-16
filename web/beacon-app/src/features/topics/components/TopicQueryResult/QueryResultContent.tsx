import { t } from '@lingui/macro';
import React from 'react';

//import { getXMLFixture , createBinaryFixture} from '../../__mocks__';
import DisplayResultData from './DisplayResultData';
export interface QueryResultContentProps {
  result: any;
  mimeType: string;
  isBase64Encoded?: boolean;
  error?: any;
  hasInvalidQuery: boolean;
}

const QueryResultContent: React.FC<QueryResultContentProps> = ({
  result,
  mimeType,
  error,
  isBase64Encoded,
  hasInvalidQuery,
}) => {
  const noQueryResult = !result && !mimeType && !hasInvalidQuery;

  return (
    <div className="shadow-md min-h-20 max-h-[480px] overflow-y-auto bg-black p-4 text-white">
      <pre className="font-base mx-auto">
        <code data-cy="topic-query-result">
          {result && (
            <DisplayResultData
              result={result}
              mimeType={mimeType}
              isBase64Encoded={isBase64Encoded}
            />
          )}
          {error && t`No results found.`}
          {noQueryResult &&
            t`No query result. Try the default query or enter your own query. See EnSQL documentation for example queries.`}
          {hasInvalidQuery &&
            t`Please enter a valid query. Please see EnSQL documentation for examples of valid queries.`}
        </code>
      </pre>
    </div>
  );
};

export default QueryResultContent;

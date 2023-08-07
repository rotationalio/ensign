import { t } from '@lingui/macro';
import React from 'react';

//import { getXMLFixture , createBinaryFixture} from '../../__mocks__';
import DisplayResultData from './DisplayResultData';
interface QueryResultContentProps {
  result: any;
  mimeType: string;
  isBase64Encoded?: boolean;
  error?: any;
}

const QueryResultContent: React.FC<QueryResultContentProps> = ({
  result,
  mimeType,
  error,
  isBase64Encoded,
}) => {
  // TODO: remove all those console.log after testing
  // console.log('[] result', result); // added this avoid eslint error
  // console.log('[] mimetype', mimeType);
  // commented out the below two lines to test the binary result
  // const mockMimeType = 'application/octet-stream';
  // const mockResult = createBinaryFixture();
  // commented out the above two lines and uncomment the below two lines to test the default result
  // result = result ?? mockResult;
  // mimeType = mimeType ?? mockMimeType;
  // commented out the above two lines and uncomment the below two lines to test the XML result
  // result = result ?? getXMLFixture();
  // mimeType = mimeType ?? 'application/xml';

  return (
    <div className="shadow-md min-h-20 max-h-[480px] overflow-y-auto bg-black p-4 text-white">
      <pre className="font-base mx-auto">
        <code>
          {result && (
            <DisplayResultData
              result={result}
              mimeType={mimeType}
              isBase64Encoded={isBase64Encoded}
            />
          )}
          {error && t`No results found.`}
          {!result &&
            !mimeType &&
            t`No query result. Try the default query or enter your own query. See EnSQL documentation for example queries.`}
        </code>
      </pre>
    </div>
  );
};

export default QueryResultContent;

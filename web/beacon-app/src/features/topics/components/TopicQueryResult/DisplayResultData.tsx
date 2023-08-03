// this component is used to display the result data according to the mimetype

import React from 'react';

import { MIME_TYPE } from '@/constants';

import { BinaryResult } from './MimeTypeResult';
interface DisplayResultDataProps {
  result: any;
  mimeType: string;
  isBase64Encoded?: boolean;
}

// Default component that should be rendered if the mimetype is binary based
const renderDefaultResultComponent = (result: any) => {
  if (result instanceof ArrayBuffer) {
    return <BinaryResult data={result} />;
  }

  return <>{result}</>;
};

const DisplayResultData: React.FC<DisplayResultDataProps> = ({ result, mimeType }) => {
  switch (mimeType) {
    case MIME_TYPE.JSON:
      return <>{JSON.stringify(result, null, 2)}</>; // TODO: add syntax highlighting with  sc-19457
    case MIME_TYPE.XML:
      return <>{result}</>; // TODO: beautify xml with 19456
    default:
      return renderDefaultResultComponent(result);
  }
};

export default DisplayResultData;

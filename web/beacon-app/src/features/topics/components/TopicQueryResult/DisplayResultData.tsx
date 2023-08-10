// this component is used to display the result data according to the mimetype

import React from 'react';

import { MIME_TYPE } from '@/constants';

import { BinaryResult, XMLResult } from './MimeTypeResult';
import HTMLResult from './MimeTypeResult/HTMLResult';
import JSONResult from './MimeTypeResult/JSONResult';

interface DisplayResultDataProps {
  result: any;
  mimeType: string;
  isBase64Encoded?: boolean;
}
const decodeBase64 = (data: string) => {
  return atob(data);
};

// Default component that should be rendered if the mimetype is binary based
const renderDefaultResultComponent = (result: any) => {
  if (result instanceof ArrayBuffer) {
    return <BinaryResult data={result} />;
  }

  return <>{result}</>;
};

const DisplayResultData: React.FC<DisplayResultDataProps> = ({
  result,
  mimeType,
  isBase64Encoded,
}) => {
  switch (mimeType) {
    case MIME_TYPE.JSON:
      return <JSONResult data={result} />;
    case MIME_TYPE.XML:
      return <XMLResult data={isBase64Encoded ? decodeBase64(result) : result} />;
    case MIME_TYPE.TEXT_HTML:
      return <HTMLResult data={result} />;
    default:
      return renderDefaultResultComponent(result);
  }
};

export default DisplayResultData;

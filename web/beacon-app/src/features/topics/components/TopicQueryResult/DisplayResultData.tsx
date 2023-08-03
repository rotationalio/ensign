// this component is used to display the result data according to the mimetype

import React from 'react';

import { MIME_TYPE } from '@/constants';

interface DisplayResultDataProps {
  result: any;
  mimeType: string;
  isBase64Encoded?: boolean;
}

type BinaryResultProps = {
  data: ArrayBuffer | null;
};

const BinaryResult: React.FC<BinaryResultProps> = ({ data }) => {
  const formatBinaryData = (binaryData: ArrayBuffer | null): string => {
    if (!binaryData) return '';

    // Convert the binary data to a Uint8Array
    const uint8Array = new Uint8Array(binaryData);
    // Convert the Uint8Array to a hexadecimal string
    const hexString = uint8Array.reduce(
      (acc, byte) => acc + byte.toString(16).padStart(2, '0'),
      ''
    );
    // Add spaces between every two hexadecimal characters for better readability
    const formattedHexString = hexString.replace(/(..)/g, '$1 ');

    return formattedHexString;
  };

  const formattedData = formatBinaryData(data);

  return <>{formattedData}</>;
};

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

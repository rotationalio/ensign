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

export default BinaryResult;

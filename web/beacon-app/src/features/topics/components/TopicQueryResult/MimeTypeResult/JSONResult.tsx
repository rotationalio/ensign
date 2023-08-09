import ReactJson from 'react-json-view';

type JSONResultProps = {
  data: any;
};

const JSONResult = ({ data }: JSONResultProps) => {
  const formattedJson = JSON.parse(data);
  return (
    <ReactJson
      src={formattedJson}
      theme="summerfruit"
      enableClipboard={false}
      displayObjectSize={false}
    />
  );
};

export default JSONResult;

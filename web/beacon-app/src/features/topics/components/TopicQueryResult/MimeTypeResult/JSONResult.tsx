import { Light as SyntaxHighlighter } from 'react-syntax-highlighter';
import { tomorrowNightBright } from 'react-syntax-highlighter/dist/esm/styles/hljs';

type JSONResultProps = {
  data: any;
};

const JSONResult = ({ data }: JSONResultProps) => {
  return (
    <SyntaxHighlighter language="json" style={tomorrowNightBright}>
      {data}
    </SyntaxHighlighter>
  );
};

export default JSONResult;

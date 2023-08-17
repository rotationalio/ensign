import pretty from 'pretty';
import { Light as SyntaxHighlighter } from 'react-syntax-highlighter';
import { irBlack } from 'react-syntax-highlighter/dist/esm/styles/hljs';

type HTMLResultProps = {
  data: any;
};

const HTMLResult = ({ data }: HTMLResultProps) => {
  const formattedHtml = pretty(data);
  return (
    <SyntaxHighlighter language="html" style={irBlack}>
      {formattedHtml}
    </SyntaxHighlighter>
  );
};

export default HTMLResult;

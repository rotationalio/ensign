// This component is responsible for rendering the XML result of a query.
import React from 'react';
import { Light as SyntaxHighlighter } from 'react-syntax-highlighter';
import { tomorrowNightBright } from 'react-syntax-highlighter/dist/esm/styles/hljs';
import xmlFormatter from 'xml-formatter';

type XMLResultProps = {
  data: any;
};

const XMLResult: React.FC<XMLResultProps> = ({ data }) => {
  const formattedData = xmlFormatter(data);

  return (
    <SyntaxHighlighter language="xml" style={tomorrowNightBright}>
      {formattedData}
    </SyntaxHighlighter>
  );
};

export default XMLResult;

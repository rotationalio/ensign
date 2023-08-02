import { Trans } from '@lingui/macro';
import { Button, TextField } from '@rotational/beacon-core';
import { useEffect, useState } from 'react';

type TopicNameProps = {
  name: string;
};

const QueryInput = ({ name }: TopicNameProps) => {
  const [topicQuery, setTopicQuery] = useState('');

  useEffect(() => {
    setTopicQuery(`SELECT * FROM ${name} LIMIT 10`);
  }, [name]);

  const handleInputQueryChange = (e: any) => {
    setTopicQuery(e.target.value);
  };

  const handleClearInputQuery = () => {
    setTopicQuery('');
  };

  return (
    <div className="mt-4 flex space-x-2">
      <TextField type="search" value={topicQuery} fullWidth onChange={handleInputQueryChange} />
      <div className="flex space-x-2">
        <Button variant="secondary">
          <Trans>Query</Trans>
        </Button>
        <Button onClick={handleClearInputQuery}>
          <Trans>Clear</Trans>
        </Button>
      </div>
    </div>
  );
};

export default QueryInput;

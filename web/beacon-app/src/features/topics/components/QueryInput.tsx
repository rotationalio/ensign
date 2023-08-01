import { Trans } from '@lingui/macro';
import { Button, TextField } from '@rotational/beacon-core';

type TopicNameProps = {
  name: string;
};

const QueryInput = ({ name }: TopicNameProps) => {
  return (
    <div className="mt-4 flex space-x-2">
      <TextField
        labelClassName="Topic Query"
        type="search"
        value={`SELECT * FROM ${name} LIMIT 1`}
        fullWidth
      ></TextField>
      <div className="flex space-x-2">
        <Button variant="secondary">
          <Trans>Query</Trans>
        </Button>
        <Button>
          <Trans>Clear</Trans>
        </Button>
      </div>
    </div>
  );
};

export default QueryInput;

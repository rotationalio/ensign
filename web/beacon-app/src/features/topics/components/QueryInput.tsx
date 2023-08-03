import { Trans } from '@lingui/macro';
import { Button, TextField } from '@rotational/beacon-core';

type TopicNameProps = {
  defaultEnSQL: string;
  queryHandler: (query: string) => void;
};

const QueryInput = ({ defaultEnSQL, queryHandler }: TopicNameProps) => {
  return (
    <div className="mt-4 flex space-x-2">
      <TextField
        labelClassName="Topic Query"
        type="search"
        defaultValue={defaultEnSQL}
        onChange={(e: any) => queryHandler(e?.target?.value ?? defaultEnSQL)}
        fullWidth
      ></TextField>
      <div className="flex space-x-2">
        <Button variant="secondary" type="submit">
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

import { Trans } from '@lingui/macro';
import { Button, TextField } from '@rotational/beacon-core';

type TopicNameProps = {
  defaultEnSQL: string;
  queryHandler: (query: string) => void;
};

const QueryInput = ({ defaultEnSQL, queryHandler }: TopicNameProps) => {
  const onQuery = (e: any) => {
    e.preventDefault();
    console.log('[] e.target.value', e.target.value);
    queryHandler(e.target.value);
  };

  return (
    <div className="mt-4 flex space-x-2">
      <TextField
        labelClassName="font-semibold"
        type="search"
        defaultValue={defaultEnSQL}
        onChange={onQuery}
        name="query"
        fullWidth
      ></TextField>
      <div className="flex space-x-2">
        <Button variant="secondary" onClick={onQuery}>
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

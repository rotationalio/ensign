import { Trans } from '@lingui/macro';
import { Heading } from '@rotational/beacon-core';

const EventDetailTableHeader = () => {
  return (
    <Heading as="h2" className="mt-8 ml-6 text-lg font-semibold">
      <Trans>Topic Usage</Trans>
    </Heading>
  );
};

export default EventDetailTableHeader;

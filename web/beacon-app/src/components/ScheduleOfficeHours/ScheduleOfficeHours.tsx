import { Trans } from '@lingui/macro';
import { Button } from '@rotational/beacon-core';
import { Link } from 'react-router-dom';

const ScheduleOfficeHours = () => {
  return (
    <Link to="https://calendly.com/rebecca-r8l" target="_blank">
      <Button variant="secondary">
        <Trans>Schedule Office Hours</Trans>
      </Button>
    </Link>
  );
};

export default ScheduleOfficeHours;

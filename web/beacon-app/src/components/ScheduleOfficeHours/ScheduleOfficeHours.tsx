import { Trans } from '@lingui/macro';
import { Button } from '@rotational/beacon-core';
import { Link } from 'react-router-dom';

import { PUBLIC_LINKS } from '@/constants';
const ScheduleOfficeHours = () => {
  return (
    <Link to={PUBLIC_LINKS.OFFICE_HOURS_SCHEDULE} target="_blank">
      <Button variant="secondary">
        <Trans>Schedule Office Hours</Trans>
      </Button>
    </Link>
  );
};

export default ScheduleOfficeHours;

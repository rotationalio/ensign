import { Trans } from '@lingui/macro';
import { Button } from '@rotational/beacon-core';
import { Link } from 'react-router-dom';

import { EXTRENAL_LINKS } from '@/application';
const ScheduleOfficeHours = () => {
  return (
    <Link to={EXTRENAL_LINKS.OFFICE_HOURS_SCHEDULE} target="_blank">
      <Button variant="secondary">
        <Trans>Schedule Office Hours</Trans>
      </Button>
    </Link>
  );
};

export default ScheduleOfficeHours;

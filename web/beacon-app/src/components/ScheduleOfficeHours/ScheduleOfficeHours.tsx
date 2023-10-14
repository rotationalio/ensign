import { t } from '@lingui/macro';
import { AiOutlineSchedule } from 'react-icons/ai';
import { Link } from 'react-router-dom';

import { EXTERNAL_LINKS } from '@/application';

import { IconTooltip } from '../common/Tooltip/IconTooltip';

const ScheduleOfficeHours = () => {
  return (
    <IconTooltip
      icon={
        <Link to={EXTERNAL_LINKS.OFFICE_HOURS_SCHEDULE} target="_blank">
          <AiOutlineSchedule className="office-hours-icon" fill="#1D65A6" fontSize={28} />
        </Link>
      }
      content={t`Schedule Office Hours`}
    />
  );
};

export default ScheduleOfficeHours;

import mergeClassnames from '../../../utils/mergeClassnames';
import { StatusColorIcon } from '../../Icon/Icons';
import { STATUS } from './util';
interface StatusPillProps {
  value: any;
  className?: string;
}
export function StatusPill({ value, className }: StatusPillProps) {
  const status = value ? value.toLowerCase() : 'unknown';

  const statusColorMap = {
    [STATUS.ACTIVE]: 'text-green-800',
    [STATUS.CONFIRMED]: 'text-green-800',
    [STATUS.COMPLETE]: 'text-green-800',
    [STATUS.PENDING]: 'text-warning-600',
    [STATUS.INCOMPLETE]: 'text-warning-600',
    [STATUS.INACTIVE]: 'text-warning-600',
    [STATUS.REVOKED]: 'text-primary-500',
    [STATUS.ERROR]: 'text-danger-600',
    [STATUS.UNUSED]: 'text-gray-600',
  } as any;

  const statusIconMap = {
    [STATUS.ACTIVE]: <StatusColorIcon fill="#34753E" />,
    [STATUS.COMPLETE]: <StatusColorIcon fill="#34753E" />,
    [STATUS.CONFIRMED]: <StatusColorIcon fill="#34753E" />,
    [STATUS.INACTIVE]: <StatusColorIcon fill="#C97900" />,
    [STATUS.PENDING]: <StatusColorIcon fill="#C97900" />,
    [STATUS.INCOMPLETE]: <StatusColorIcon fill="#C97900" />,
    [STATUS.REVOKED]: <StatusColorIcon fill="#F26800" />,
    [STATUS.ERROR]: <StatusColorIcon fill="#EB2A00" />,
    [STATUS.UNUSED]: <StatusColorIcon fill="#6C757D" />,
  } as any;

  return (
    <div className={mergeClassnames('flex items-center', className)}>
      {statusIconMap[status]}
      <span className={mergeClassnames('ml-1', statusColorMap[status] as string)}>{value}</span>
    </div>
  );
}

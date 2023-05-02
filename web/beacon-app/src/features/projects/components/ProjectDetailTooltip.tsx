import { Trans } from '@lingui/macro';
import * as Tooltip from '@radix-ui/react-tooltip';

import { SentryErrorBoundary } from '@/components/Error';
import HintIcon from '@/components/icons/hint';
import { formatDate } from '@/utils/formatDate';

import type { Project } from '../types/Project';
interface ProjectDetailTooltipProps {
  data: Project;
}

const ProjectDetailTooltip = ({ data }: ProjectDetailTooltipProps) => {
  const { name, description, status, created } = data || {};
  return (
    <SentryErrorBoundary
      fallback={
        <div>
          <Trans>Something went wrong while rendering the project detail tooltip.</Trans>
        </div>
      }
    >
      <Tooltip.Provider>
        <Tooltip.Root>
          <Tooltip.Trigger asChild>
            <button className="">
              <HintIcon />
            </button>
          </Tooltip.Trigger>
          <Tooltip.Portal>
            <Tooltip.Content
              className="w-full max-w-[550px] rounded-md bg-secondary-slate p-4 text-sm text-white"
              sideOffset={5}
              align="start"
            >
              <table className="table-auto border-separate border-spacing-y-2">
                <tbody>
                  <tr>
                    <td className="w-[150px] font-semibold">
                      <Trans>Project Status:</Trans>
                    </td>
                    <td>{status}</td>
                  </tr>
                  <tr>
                    <td className="font-semibold">
                      <Trans>Description:</Trans>
                    </td>
                    <td>{description}</td>
                  </tr>
                  <tr>
                    <td className="font-semibold">
                      <Trans>Owner:</Trans>
                    </td>
                    <td>{name}</td>
                  </tr>
                  <tr>
                    <td className="font-semibold">
                      <Trans>Created:</Trans>
                    </td>
                    <td>{formatDate(new Date(created))}</td>
                  </tr>
                </tbody>
              </table>
            </Tooltip.Content>
          </Tooltip.Portal>
        </Tooltip.Root>
      </Tooltip.Provider>
    </SentryErrorBoundary>
  );
};

export default ProjectDetailTooltip;

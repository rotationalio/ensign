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
  const { description, status, created, owner } = data || {};

  const getFormattedDecription = () => {
    if (!description) {
      return '---';
    }
    // cut off description at 100 characters

    return description.length > 100 ? `${description.slice(0, 100)}...` : description;
  };
  return (
    <SentryErrorBoundary
      fallback={
        <div className="item-center justify-center text-center text-sm text-danger">
          <Trans>
            Sorry, we were unable to load the project details. You can either refresh the page or
            get in touch with our support team for assistance.
          </Trans>
        </div>
      }
    >
      <Tooltip.Provider>
        <Tooltip.Root>
          <Tooltip.Trigger asChild>
            <button className="" data-cy="detailHint">
              <HintIcon />
            </button>
          </Tooltip.Trigger>
          <Tooltip.Portal>
            <Tooltip.Content
              className="w-full max-w-[550px] rounded-md bg-secondary-slate p-4 text-sm text-white"
              sideOffset={5}
              align="start"
              data-cy="prjDetail"
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
                    <td>{getFormattedDecription()}</td>
                  </tr>
                  <tr>
                    <td className="font-semibold">
                      <Trans>Owner:</Trans>
                    </td>
                    <td>{owner?.name}</td>
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

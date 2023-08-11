import { Trans } from '@lingui/macro';

import { SentryErrorBoundary } from '@/components/Error';
import { formatDate } from '@/utils/formatDate';

import type { Project } from '../../types/Project';
interface ProjectDetailInfoProps {
  data: Project;
}

const ProjectDetailInfo: React.FC<ProjectDetailInfoProps> = ({ data }) => {
  const { description, created, owner } = data || {};

  const getFormattedDecription = () => {
    if (!description) {
      return '---';
    }

    return description;
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
      <div className="my-[10px] border-b border-gray-500 ">
        <table className=" table-auto border-separate border-spacing-y-4">
          <tbody>
            <tr>
              <td className="w-[150px] font-semibold">
                <Trans>Description:</Trans>
              </td>
              <td>{getFormattedDecription()}</td>
            </tr>
            <tr>
              <td className="w-[150px] font-semibold">
                <Trans>Owner</Trans>
              </td>
              <td>{owner?.name}</td>
            </tr>
            <tr>
              <td className="w-[150px] font-semibold">
                <Trans>Created</Trans>
              </td>
              <td>{formatDate(new Date(created))}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </SentryErrorBoundary>
  );
};

export default ProjectDetailInfo;

import { Trans } from '@lingui/macro';
import { Heading } from '@rotational/beacon-core';

import { TagState } from '@/components/common/TagState';
import { SentryErrorBoundary } from '@/components/Error';
import { formatDate } from '@/utils/formatDate';

import type { Project } from '../../types/Project';
interface ProjectDetailInfoProps {
  data: Project;
}

const ProjectDetailInfo: React.FC<ProjectDetailInfoProps> = ({ data }) => {
  const { description, status, created, owner } = data || {};

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
      <div className="mt-[46px]">
        <Heading as={'h1'} className=" flex items-center text-lg font-semibold capitalize">
          <Trans>Project Details</Trans>
        </Heading>
        <table className="mt-4 table-auto border-separate border-spacing-y-4">
          <tbody>
            <tr>
              <td className="w-[150px] font-semibold">
                <Trans>Project Status:</Trans>
              </td>
              <td>
                <TagState status={status as string} />
              </td>
            </tr>
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

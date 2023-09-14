import { Trans } from '@lingui/macro';

import { SentryErrorBoundary } from '@/components/Error';
import { formatDate } from '@/utils/formatDate';

import type { MemberResponse } from '../types/memberServices';

interface MemberDetailInfoProps {
  data: MemberResponse;
}

const MemberDetailInfo: React.FC<MemberDetailInfoProps> = ({ data }) => {
  const { id, name, role, created } = data || {};

  return (
    <SentryErrorBoundary
      fallback={
        <div className="item-center justify-center text-center text-sm text-danger">
          <Trans>
            Sorry, we were unable to load the user details. You can either refresh the page or get
            in touch with our support team for assistance.
          </Trans>
        </div>
      }
    >
      <div className="my-[10px] ">
        <table className=" table-auto border-separate border-spacing-y-4">
          <tbody>
            <tr>
              <td className="w-[150px] font-semibold">
                <Trans>User ID:</Trans>
              </td>
              <td>{id}</td>
            </tr>
            <tr>
              <td className="w-[150px] font-semibold">
                <Trans>Name </Trans>
              </td>
              <td>{name}</td>
            </tr>
            <tr>
              <td className="w-[150px] font-semibold">
                <Trans>Role</Trans>
              </td>
              <td>{role}</td>
            </tr>
            <tr>
              <td className="w-[150px] font-semibold">
                <Trans>Date Created</Trans>
              </td>
              <td>{formatDate(new Date(created))}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </SentryErrorBoundary>
  );
};

export default MemberDetailInfo;

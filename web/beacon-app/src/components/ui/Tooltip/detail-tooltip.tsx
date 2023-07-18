import { Trans } from '@lingui/macro';
import * as Tooltip from '@radix-ui/react-tooltip';
import { mergeClassnames } from '@rotational/beacon-core';

import { SentryErrorBoundary } from '@/components/Error';
import HintIcon from '@/components/icons/hint';

interface DetailTooltipProps {
  headers: string[];
  data: string[];
}

const DetailTooltip = ({ data, headers }: DetailTooltipProps) => {
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
                  {headers.map((header, index) => (
                    <tr key={index}>
                      <td
                        className={mergeClassnames('font-semibold', index === 0 ? 'w-[150px]' : '')}
                      >
                        {header}
                      </td>
                      <td>{data[index]}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </Tooltip.Content>
          </Tooltip.Portal>
        </Tooltip.Root>
      </Tooltip.Provider>
    </SentryErrorBoundary>
  );
};

export default DetailTooltip;

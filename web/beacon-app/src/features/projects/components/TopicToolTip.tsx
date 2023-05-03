import { Trans } from '@lingui/macro';
import * as Tooltip from '@radix-ui/react-tooltip';

import HintIcon from '@/components/icons/hint';

function TopicTooltip() {
  return (
    <>
      <Tooltip.Provider>
        <Tooltip.Root>
          <Tooltip.Trigger asChild>
            <button>
              <HintIcon />
            </button>
          </Tooltip.Trigger>
          <Tooltip.Portal>
            <Tooltip.Content
              className="w-full max-w-[275px] rounded-md bg-secondary-slate p-4 text-sm text-white"
              sideOffset={5}
              align="start"
            >
              <p>
                <Trans>
                  Messages and events are sent to and read from specific topics. Services that are{' '}
                  {''}
                  <span className="font-bold">producers, write</span> data to topics. Services that
                  are <span className="font-bold">consumers, read</span> data from topics. Topics
                  are multi-subscriber, which means that a topic can have zero, one, or multiple
                  consumers subscribing to that topic, with read access to the log.
                </Trans>
              </p>
            </Tooltip.Content>
          </Tooltip.Portal>
        </Tooltip.Root>
      </Tooltip.Provider>
    </>
  );
}

export default TopicTooltip;

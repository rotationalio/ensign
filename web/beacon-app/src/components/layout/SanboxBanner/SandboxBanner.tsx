import { Trans } from '@lingui/macro';

import { EXTERNAL_LINKS } from '@/application';
import Alert from '@/components/common/Alert/Alert';

function SandboxBanner() {
  return (
    <Alert>
      <div className="flex h-auto w-full flex-col items-center justify-center gap-x-4 bg-[#EBF5FF] py-6 text-center font-bold text-[#1D65A6] lg:flex-row">
        <p>
          <Trans>
            You are using the Ensign Sandbox. Ready to deploy your models to production?
          </Trans>
        </p>
        <a
          href={EXTERNAL_LINKS.ENSIGN_PRICING}
          target="_blank"
          rel="noreferrer"
          className="mt-4 rounded-md border border-white bg-[#1D65A6] px-4 py-1.5 text-white hover:bg-[#1D65A6/50] lg:mt-0"
        >
          <Trans>Upgrade</Trans>
        </a>
      </div>
    </Alert>
  );
}

export default SandboxBanner;

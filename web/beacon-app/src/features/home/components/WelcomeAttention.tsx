import { Trans } from '@lingui/macro';
import { Button } from '@rotational/beacon-core';
import { useNavigate } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application';
const WelcomeAttention = () => {
  const navigate = useNavigate();
  const LINK = 'https://ensign.rotational.dev/';

  const redirectTo = () => {
    navigate(PATH_DASHBOARD.PROJECTS);
  };
  return (
    <>
      <div
        className="px-auto mb-8 mt-4 flex flex-row items-center justify-between space-x-4 rounded-md border border-neutral-500 bg-[#F7F9FB] p-2 px-5 text-justify"
        data-cy="projWelcome"
      >
        <p className="text-md">
          <Trans>
            Welcome to Ensign! Set up or manage your projects. A project is{' '}
            <a
              href={LINK}
              target="_blank"
              rel="noreferrer"
              className="font-bold text-[#1D65A6] underline hover:!underline"
            >
              a database for events.
            </a>{' '}
            We’ll guide you along the way!
          </Trans>
        </p>

        <Button variant="tertiary" size="small" onClick={redirectTo} data-cy="startSetupBttn">
          <Trans>Start</Trans>
        </Button>
      </div>
    </>
  );
};

export default WelcomeAttention;

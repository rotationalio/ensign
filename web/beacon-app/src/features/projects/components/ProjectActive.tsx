import { Trans } from '@lingui/macro';
import { Card } from '@rotational/beacon-core';

function ProjectActive() {
  const SDKLink = 'https://ensign.rotational.dev/sdk/';
  const DocsLink = 'https://ensign.rotational.dev/';
  const ExampleLink = 'https://github.com/rotationalio/ensign-examples';
  return (
    <>
      <Card
        style={{ borderRadius: '4px' }}
        contentClassName="m-[16px] w-full"
        className="mt-8 mb-8 w-full border-[1px] border-gray-600 p-[4px]"
      >
        <Card.Body>
          <Trans>
            Your project is active! Check out our{' '}
            <a href={SDKLink} target="_blank" rel="noreferrer" className="underline">
              SDKs,
            </a>{' '}
            <a href={DocsLink} target="_blank" rel="noreferrer" className="underline">
              documentation,
            </a>{' '}
            and{' '}
            <a href={ExampleLink} target="_blank" rel="noreferrer" className="underline">
              example code
            </a>{' '}
            to connect publishers and subscribers to your project (database).
          </Trans>
        </Card.Body>
      </Card>
    </>
  );
}

export default ProjectActive;

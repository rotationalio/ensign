import { Trans } from '@lingui/macro';
import { Button, Card } from '@rotational/beacon-core';
import { Dispatch, SetStateAction } from 'react';
import { AiOutlineClose } from 'react-icons/ai';
interface ProjectActiveProps {
  onActive: Dispatch<SetStateAction<boolean>>;
  projectID: string;
}
function ProjectActive({ onActive, projectID }: ProjectActiveProps) {
  const SDKLink = 'https://ensign.rotational.dev/sdk/';
  const DocsLink = 'https://ensign.rotational.dev/';
  const ExampleLink = 'https://github.com/rotationalio/ensign-examples';
  const closeBtnHandler = () => {
    const projectName = 'isActiveProject-' + projectID;
    localStorage.setItem(projectName, 'true');
    onActive(true);
  };
  return (
    <>
      <Card
        style={{ borderRadius: '4px', border: '1px solid #6DD19C' }}
        contentClassName="m-[16px] w-full"
        className="mt-8 mb-8 w-full border-[1px] border-green-800 bg-green/20 p-[4px]"
      >
        <Card.Body>
          <div className="flex items-center justify-between">
            <div>
              <Trans>
                Your project is active! Check out our{' '}
                <a
                  href={SDKLink}
                  target="_blank"
                  rel="noreferrer"
                  className="font-bold text-[#1D65A6] underline"
                >
                  SDKs,
                </a>{' '}
                <a
                  href={DocsLink}
                  target="_blank"
                  rel="noreferrer"
                  className="font-bold text-[#1D65A6] underline"
                >
                  documentation,
                </a>{' '}
                and{' '}
                <a
                  href={ExampleLink}
                  target="_blank"
                  rel="noreferrer"
                  className="font-bold text-[#1D65A6] underline"
                >
                  example code
                </a>{' '}
                to connect publishers and subscribers to your project (database).
              </Trans>
            </div>
            <Button
              onclick={closeBtnHandler}
              variant="ghost"
              size="custom"
              className="bg-transparent  h-4 w-4 border-none py-0"
            >
              <AiOutlineClose onClick={closeBtnHandler} />
            </Button>
          </div>
        </Card.Body>
      </Card>
    </>
  );
}

export default ProjectActive;

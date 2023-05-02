import { Trans } from '@lingui/macro';
import { Card, Heading } from '@rotational/beacon-core';

import HeavyCheckMark from '@/components/icons/heavy-check-mark';
interface ProjectSetupProps {
  warningMessage?: string;
  config: {
    isProjectCreated: boolean;
    isAPIKeyCreated: boolean;
    isTopicCreated: boolean;
  };
}
const ProjectSetup = ({ config, warningMessage }: ProjectSetupProps) => {
  const { isProjectCreated, isAPIKeyCreated, isTopicCreated } = config;
  return (
    <>
      <Card
        style={{ padding: '4px', borderRadius: '4px', marginTop: '2rem', marginBottom: '2rem' }}
        contentClassName="m-[16px] w-full rounded-[4px]"
        className="min-h-[150px] w-full border-[1px] border-gray-600 p-4"
        data-testid="project-setup"
      >
        <Card.Header>
          <Heading as="h3" className="pb-2 text-warning-500">
            <Trans>Project setup is incomplete. {warningMessage}</Trans>
          </Heading>
        </Card.Header>

        <Card.Body>
          <div className="space-y-4 pl-2">
            <table cellPadding={4} className="table-auto border-separate">
              <tr>
                <td>1.</td>
                <td>
                  <Trans>Create Project</Trans>
                </td>
                <td data-testid="project-created">
                  {isProjectCreated ? <HeavyCheckMark width={16} height={16} /> : null}
                </td>
              </tr>
              <tr>
                <td>2.</td>
                <td>
                  <Trans>Generate API Keys</Trans>
                </td>
                <td data-testid="api-key-created">
                  {isAPIKeyCreated ? <HeavyCheckMark width={16} height={16} /> : null}
                </td>
              </tr>

              <tr>
                <td>3.</td>
                <td>
                  <Trans>Create Topics</Trans>
                </td>
                <td data-testid="topic-created">
                  {isTopicCreated ? <HeavyCheckMark width={16} height={16} /> : null}
                </td>
              </tr>
            </table>
          </div>
        </Card.Body>
      </Card>
    </>
  );
};

export default ProjectSetup;

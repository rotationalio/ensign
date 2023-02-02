import { Button, Card } from '@rotational/beacon-core';

function ProjectDetails() {
  return (
    <>
      <Card contentClassName="w-full min-h-[200px] border border-primary-900 rounded-md p-4">
        <Card.Header>
          <h1 className="px-2 font-bold">Step 1: View Project Details</h1>
        </Card.Header>
        <Card.Body>
          <div className="space-y-3">
            <div className="my-3 mb-5 flex flex-col items-start justify-between gap-3 px-2 sm:mb-0 sm:flex-row sm:gap-0">
              <p className="text-sm sm:w-4/5">
                View project details below. Generate your API key next to connect producers and
                consumers to Ensign and start managing your project.
              </p>
              <div className="sm:w-1/5">
                <Button className="h-auto text-sm">Manage Project</Button>
              </div>
            </div>
            <table className="border-separate border-spacing-x-2 border-spacing-y-1 text-sm">
              <tr>
                <td className="font-bold">Project name</td>
                <td>Acme System Global Event Stream</td>
              </tr>
              <tr>
                <td className="font-bold">Topic name</td>
                <td>acsys</td>
              </tr>
              <tr>
                <td className="font-bold">Date created</td>
                <td>01/01/2023</td>
              </tr>
              <tr>
                <td className="font-bold">Project state</td>
                <td>Ready for use</td>
              </tr>
              <tr>
                <td className="font-bold">API key</td>
                <td>Not generated</td>
              </tr>
            </table>
          </div>
        </Card.Body>
      </Card>
    </>
  );
}

export default ProjectDetails;

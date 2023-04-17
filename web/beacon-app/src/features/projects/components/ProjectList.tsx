import { AriaButton as Button, Heading } from '@rotational/beacon-core';
import { Link } from 'react-router-dom';

import Union from '@/components/icons/union';

import ProjectsTable from './ProjectsTable';

function ProjectList() {
  return (
    <>
      <Heading as="h1" className="mb-4 mt-6 text-lg font-semibold">
        Projects
      </Heading>
      <div className="flex justify-between rounded-lg bg-[#F7F9FB] px-3 py-2">
        <div className="flex items-center gap-3"></div>
        <div>
          <Link to="project-setup">
            <Button className="flex items-center gap-1 text-xs" size="small">
              <Union className="fill-white" />
              Create Project
            </Button>
          </Link>
        </div>
      </div>
      <ProjectsTable />
    </>
  );
}

export default ProjectList;

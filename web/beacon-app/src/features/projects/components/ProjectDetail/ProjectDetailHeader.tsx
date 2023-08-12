import { Heading } from '@rotational/beacon-core';
import React from 'react';

import { TagState } from '@/components/common/TagState';
interface ProjectDetailHeaderProps {
  data: any;
}

const ProjectDetailHeader: React.FC<ProjectDetailHeaderProps> = ({ data }) => {
  const { status, name } = data || {};
  const getNormalizedProjectName = () => {
    return name?.split('-').join(' ');
  };

  return (
    <Heading as="h1" className="flex items-center gap-5 text-2xl font-semibold">
      <span className="mr-2" data-cy="project-name">
        {getNormalizedProjectName()}
      </span>
      <TagState status={status as string} />
    </Heading>
  );
};

export default ProjectDetailHeader;

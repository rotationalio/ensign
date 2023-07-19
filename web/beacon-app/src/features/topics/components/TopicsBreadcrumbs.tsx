import { Breadcrumbs } from '@rotational/beacon-core';
import React, { useCallback } from 'react';
import { useNavigate } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application/routes/paths';
import BreadcrumbsIcon from '@/components/ui/Breadcrumbs/breadcrumbs-icon';

import type { Topic } from '../types/topicService';

interface TopicBreadcrumbsProps {
  topic: Partial<Topic>;
}

const TopicsBreadcrumbs = ({ topic }: TopicBreadcrumbsProps) => {
  const { topic_name: name } = topic || {};
  const navigate = useNavigate();
  const CustomBreadcrumbs = useCallback(() => {
    return (
      <Breadcrumbs separator="/" className="ml-4 hidden md:block">
        <Breadcrumbs.Item
          onClick={() => navigate(PATH_DASHBOARD.HOME)}
          className="capitalize hover:underline"
        >
          <BreadcrumbsIcon className="inline" /> Home
        </Breadcrumbs.Item>
        <Breadcrumbs.Item className="capitalize" onClick={() => navigate(PATH_DASHBOARD.PROJECTS)}>
          Projects
        </Breadcrumbs.Item>
        <Breadcrumbs.Item className="capitalize">Topics</Breadcrumbs.Item>
        {name ? <Breadcrumbs.Item>{name}</Breadcrumbs.Item> : null}
      </Breadcrumbs>
    );
  }, [name, navigate]);

  return (
    <>
      <CustomBreadcrumbs />
    </>
  );
};

export default TopicsBreadcrumbs;

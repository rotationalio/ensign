// import { Project } from '../types';
import { t } from '@lingui/macro';

import { useOrgStore } from '@/store';
export const getRecentProject = (projects: any) => {
  const org = useOrgStore.getState() as any;
  let p = [] as any;
  if (projects && projects?.tenant_projects?.length) {
    const recent = projects?.tenant_projects[0]; // TODO: get most recent project instead of first

    const { name, id } = recent; // The project object response from the API
    org.setProjectID(id);
    org.setProject({
      name,
      id,
    });
    p = [
      {
        label: 'Project Name',
        value: name,
      },
      {
        label: 'Project ID',
        value: id,
      },
    ];
  }
  return p;
};

export const getDefaultHomeStats = () => {
  return [
    {
      name: t`Active Projects`,
      value: 0,
    },
    {
      name: t`Topics`,
      value: 0,
    },
    {
      name: t`API Keys`,
      value: 0,
    },
    {
      name: t`Storage`,
      value: 0,
      units: 'GB',
    },
  ];
};
export const getHomeStatsHeaders = () => {
  return [t`Active Projects`, t`Topics`, t`API Keys`, t`Storage`];
};

export const STARTER_VIDEOS = [
  {
    title: 'Batch vs Streaming',
    preview_image: 'https://i.ytimg.com/vi/HDRQ9Fe9g7c/maxres1.jpg',
    ytVideoId: 'HDRQ9Fe9g7c',
  },
  {
    title: 'Data Architecture',
    preview_image: 'https://i.ytimg.com/vi/3AxNSJ9oB24/maxres1.jpg',
    ytVideoId: '3AxNSJ9oB24',
  },
  {
    title: 'Creating Projects',
    // preview_image: 'https://img.youtube.com/vi/VskNgAVMORQ/1.jpg',
    preview_image: 'https://i.ytimg.com/vi/VskNgAVMORQ/maxres1.jpg',
    ytVideoId: 'VskNgAVMORQ',
  },

  {
    title: 'Naming Topics',
    preview_image: 'https://i.ytimg.com/vi/1XuVPl_Ki4U/maxres1.jpg',
    ytVideoId: '1XuVPl_Ki4U',
  },
  {
    title: 'Creating API Keys',
    preview_image: 'https://i.ytimg.com/vi/KMejrUIouMw/maxres1.jpg',
    ytVideoId: 'KMejrUIouMw',
  },
  {
    title: 'Protecting Your API Keys',
    preview_image: 'https://i.ytimg.com/vi/EEpIDkKJopY/maxres1.jpg',
    ytVideoId: 'EEpIDkKJopY',
  },
];

export const TEMPLATE_DATA = [
  {
    title: 'Data Flow Templates',
    links: [
      {
        name: 'Sentiment Analysis',
        url: 'https://ensign.rotational.dev/getting-started/data-flow/',
      },
      {
        name: 'NLP',
        url: 'https://www.youtube.com/watch?v=Zz0Cw8x2X0M',
      },
      {
        name: 'Detection',
        url: 'https://ensign.rotational.dev/getting-started/data-flow/',
      },
      {
        name: 'IoT',
        url: 'https://www.youtube.com/watch?v=Zz0Cw8x2X0M',
      },
    ],
  },
  {
    title: 'Code Examples',
    links: [
      {
        name: 'Data Archtecture',
        url: 'https://ensign.rotational.dev/getting-started/data-architecture/',
      },
      {
        name: 'Creating Projects',
        url: 'https://www.youtube.com/watch?v=VskNgAVMORQ',
      },
      {
        name: 'Naming Topics',
        url: 'https://www.youtube.com/watch?v=1XuVPl_Ki4U',
      },
      {
        name: 'Creating API Keys',
        url: 'https://www.youtube.com/watch?v=KMejrUIouMw',
      },
    ],
  },
  {
    title: 'Code Challenges',
    links: [
      {
        name: 'Publish to Topic',
        url: 'https://ensign.rotational.dev/getting-started/data-flow/',
      },
      {
        name: 'Subscribe to Topic',
        url: 'https://www.youtube.com/watch?v=Zz0Cw8x2X0M',
      },
      {
        name: 'Model Deployment',
        url: 'https://ensign.rotational.dev/getting-started/data-flow/',
      },
      {
        name: 'Monitoring Service',
        url: 'https://www.youtube.com/watch?v=Zz0Cw8x2X0M',
      },
    ],
  },
];

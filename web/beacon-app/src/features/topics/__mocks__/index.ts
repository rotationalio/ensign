export const getTopicsStatsMockData = () => {
  const data = [
    {
      name: 'publishers',
      value: 2,
    },
    {
      name: 'subscribers',
      value: 3,
    },
    {
      name: 'total_events',
      value: 1000000,
    },
    {
      name: 'storage',
      value: 203,
      units: 'MB',
    },
  ];
  return data;
};

export const getTopicEventsMockData = () => {
  return [
    {
      type: 'Document',
      version: '1.0.0',
      mimetype: 'application/json',
      events: {
        value: 12345678,
        percent: 96.0,
      },
      storage: {
        value: 512,
        units: 'MB',
        percent: 98.5,
      },
    },
    {
      type: 'Feed Item',
      version: '0.8.1',
      mimetype: 'application/rss',
      events: {
        value: 98765,
        percent: 4.0,
      },
      storage: {
        value: 4.3,
        units: 'KB',
        percent: 1.5,
      },
    },
  ];
};

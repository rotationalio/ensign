import { render, screen } from '@testing-library/react';
import commaNumber from 'comma-number';
import React from 'react';
import { vi } from 'vitest';

import { getTopicEventsMockData } from '../../__mocks__';
import { TopicEvents } from '../../types/topicEventsService';
import EventDetailTable from '../EventDetailTable';
// Mock the custom hook `useFetchTopicEvents`
vi.mock('../../hooks/useFetchTopicEvents', () => ({
  useFetchTopicEvents: () => ({
    topicEvents: getTopicEventsMockData(),
    isFetchingTopicEvents: false,
  }),
}));

const renderComponent = () => {
  return render(<EventDetailTable />);
};

// mock getEventDetailColumns
vi.mock('../../utils', async (importOrginial: any) => ({
  ...importOrginial,
  getEventDetailColumns: () => [
    {
      Header: `Event Type`,
      accessor: 'type',
    },
    {
      Header: `Version`,
      accessor: 'version',
    },
    {
      Header: `MIME Type`,
      accessor: 'mimetype',
    },
    // all columns with accessor function are not working yet
    // TODO: fix this later
    {
      Header: `# of Events`,
      accessor: (event: TopicEvents) => {
        return `${event?.events?.value}`;
      },
    },
    {
      Header: `% of Events`,
      accessor: (event: TopicEvents) => {
        return `${event?.events?.percent}%`;
      },
    },
    {
      Header: `Storage Volume`,
      accessor: (event: TopicEvents) => {
        return `${event?.storage?.value} ${event?.storage?.units}`;
      },
    },
    {
      Header: `% of Volume`,
      accessor: (event: TopicEvents) => {
        return `${event?.storage?.percent}%`;
      },
    },
  ],
  getFormattedEventDetailData: (data) => {
    return data.map((event: any) => ({
      ...event,
      events: {
        ...event?.events,
        value: commaNumber(event?.events?.value, '.', ','),
      },
    }));
  },
}));

const MockTable = ({ columns, data, isLoading }: any) => {
  return (
    <table data-testid="event-detail-table">
      <thead>
        <tr>
          {columns.map((column) => (
            <th key={column.Header}>{column.Header}</th>
          ))}
        </tr>
      </thead>
      <tbody>
        {isLoading && <div>Loading...</div>}
        {data.map((row, rowIndex) => (
          <tr key={rowIndex}>
            {columns.map((column: any) => (
              <td key={column.accessor}>{row[column.accessor]}</td>
            ))}
          </tr>
        ))}
      </tbody>
    </table>
  );
};
vi.mock('react-router-dom', async (importOrginial: any) => ({
  ...importOrginial,
  useNavigate: () => vi.fn(),
  useParams: () => ({ id: '1' }),
}));
// mock
//  mock table from beacon-core
vi.mock('@rotational/beacon-core', async (importOrginial: any) => ({
  ...importOrginial,
  Table: (props) => <MockTable {...props} />,
}));

vi.mock('@lingui/macro', async (importOrginial: any) => ({
  ...importOrginial,
  t: () => vi.fn(),
  Trans: () => vi.fn(),
}));

describe('EventDetailTable', () => {
  it('should render the component with default value', () => {
    const { container } = renderComponent();
    expect(container).toMatchSnapshot();
  });

  it('should render the correct columns', () => {
    renderComponent();
    expect(screen.getByText('Event Type')).toBeInTheDocument();
    expect(screen.getByText('Version')).toBeInTheDocument();
    expect(screen.getByText('MIME Type')).toBeInTheDocument();
    expect(screen.getByText('# of Events')).toBeInTheDocument();
    expect(screen.getByText('% of Events')).toBeInTheDocument();
    expect(screen.getByText('Storage Volume')).toBeInTheDocument();
    expect(screen.getByText('% of Volume')).toBeInTheDocument();
  });

  it('should contain the correct value', () => {
    renderComponent();
    expect(screen.getByText('Document')).toBeInTheDocument();
    expect(screen.getByText('application/json')).toBeInTheDocument();
  });
});

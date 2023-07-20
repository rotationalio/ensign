import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

/// import { dynamicActivate } from '../../../../I18n';
import TopicQuickView from '../TopicQuickView';

// mock SentryErrorBoundary
vi.mock('SentryErrorBoundary', () => ({
  __esModule: true,
  default: ({ children }: any) => <div>{children}</div>,
}));

// mock Suspense
vi.mock('Suspense', () => ({
  __esModule: true,
  default: ({ children }: any) => <div>{children}</div>,
}));
// mock useFetchTopicStats
const renderComponent = (props: { topicID: string }) => {
  const queryClient = new QueryClient();
  const wrapper = ({ children }: any) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
  return render(<TopicQuickView {...props} />, { wrapper });
};

// mock Trans tag from lingui
vi.mock('@lingui/macro', () => ({
  t: (str) => str,
  Trans: ({ children }) => children,
}));

describe('TopicQuickView', () => {
  it('should render the component', () => {
    vi.mock('../../hooks/useFetchTopicStats', () => ({
      __esModule: true,
      default: () => ({
        topicStats: [
          {
            name: 'Online Publishers',
            value: 1,
          },
          {
            name: 'Online Subscribers',
            value: 2,
          },
          {
            name: 'Events',
            value: 3,
          },
          {
            name: 'Data Storage',
            value: 4,
            units: 'GB',
          },
        ] as IStats[],

        error: false,
      }),
    }));
    const { container } = renderComponent({ topicID: '1' });
    expect(container).toMatchSnapshot();
  });

  it('should return the right data', () => {
    vi.mock('../../hooks/useFetchTopicStats', () => ({
      __esModule: true,
      default: () => ({
        topicStats: [
          {
            name: 'Online Publishers',
            value: 1,
          },
          {
            name: 'Online Subscribers',
            value: 2,
          },
          {
            name: 'Avg Events/ Second',
            value: 3,
            units: 'eps',
          },
          {
            name: 'Data Storage',
            value: 4,
            units: 'GB',
          },
        ] as IStats[],

        error: false,
      }),
    }));
    renderComponent({ topicID: '1' });
    expect(screen.getByTestId('quick-view-card-0')).toHaveTextContent('1');
    expect(screen.getByTestId('quick-view-card-1')).toHaveTextContent('2');
    expect(screen.getByTestId('quick-view-card-2')).toHaveTextContent('3');
    expect(screen.getByTestId('quick-view-card-3')).toHaveTextContent('4');
    expect(screen.getByTestId('quick-view-card-3')).toHaveTextContent('GB');
  });

  it('should return default values if error', () => {
    vi.mock('../hooks/useFetchTopicStats', () => ({
      __esModule: true,
      default: () => ({
        topicStats: null,
        error: true,
      }),
    }));
    renderComponent({ topicID: '' });

    expect(screen.getByTestId('quick-view-card-0')).toHaveTextContent('0');
    expect(screen.getByTestId('quick-view-card-1')).toHaveTextContent('0');
    expect(screen.getByTestId('quick-view-card-2')).toHaveTextContent('0');
    expect(screen.getByTestId('quick-view-card-3')).toHaveTextContent('0');
  });
});

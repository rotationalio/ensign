import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

import TopicQuickView from '../TopicQuickView';

// mock SentryErrorBoundary
vi.mock('SentryErrorBoundary', () => ({
  __esModule: true,
  default: ({ children }) => children,
}));

// mock Suspense
vi.mock('Suspense', () => ({
  __esModule: true,
  default: ({ children }) => children,
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

// vi.mock("src/features/topics/hooks/useFetchTopicStats", async () => {
//   const actual = await vi.importActual("src/features/topics/hooks/useFetchTopicStats");
//   return {
//     ...actual,
//     // your mocked methods
//   },
// })

// mock useFetchTopicStats return value

describe('TopicQuickView', () => {
  it('should render the component with default value', () => {
    const { container } = renderComponent({ topicID: '1' });
    expect(container).toMatchSnapshot();
    expect(screen.getByTestId('quick-view-card-0')).toHaveTextContent('0');
    expect(screen.getByTestId('quick-view-card-1')).toHaveTextContent('0');
    expect(screen.getByTestId('quick-view-card-2')).toHaveTextContent('0');
    expect(screen.getByTestId('quick-view-card-3')).toHaveTextContent('0');
  });
});

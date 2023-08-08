/* eslint-disable testing-library/no-wait-for-multiple-assertions */
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import React from 'react';
import { vi } from 'vitest';

import { getTopicQueryResponseMockData } from '../../__mocks__';
import TopicQueryResult from '../TopicQueryResult';

// mock Trans tag from lingui
vi.mock('@lingui/macro', () => ({
  t: (str: string) => str,
  Trans: ({ children }) => children,
}));

const renderComponent = () => {
  const propsMock = getTopicQueryResponseMockData();
  return render(<TopicQueryResult data={propsMock} />);
};

describe('TopicQueryResult', () => {
  it('should render the component with default value', () => {
    const { container } = renderComponent();
    expect(container).toMatchSnapshot();
  });

  it('should render with the first result data', async () => {
    renderComponent();

    expect(screen.getByText('1 results of 10 total')).toBeInTheDocument();
    expect(screen.getByText('text/plain')).toBeInTheDocument();
    expect(screen.getByText('Message v1.0.0')).toBeInTheDocument();
    expect(screen.getByText('hello world')).toBeInTheDocument();

    expect(screen.getByTestId('prev-query-btn')).toBeDisabled();

    expect(screen.getByTestId('next-query-btn')).toBeEnabled();
  });

  it('should display the correct results when the next button is clicked', async () => {
    renderComponent();
    fireEvent.click(screen.getByTestId('next-query-btn')); // Simulate a click

    await waitFor(() => {
      expect(screen.getByText('text/csv')).toBeInTheDocument();
      expect(screen.getByText('Spreadsheet v1.1.0')).toBeInTheDocument();
      expect(screen.getByText('hello,world')).toBeInTheDocument();
      expect(screen.getByText('2 results of 10 total')).toBeInTheDocument();
    });
    // Additional assertions after the state has loaded
    expect(screen.getByTestId('prev-query-btn')).toBeEnabled();
    expect(screen.getByTestId('next-query-btn')).toBeEnabled();
  });
});

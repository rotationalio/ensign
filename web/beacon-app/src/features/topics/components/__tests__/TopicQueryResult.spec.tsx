import { fireEvent, render, screen } from '@testing-library/react';
import React from 'react';
import { vi } from 'vitest';

import { getTopicQueryResponseMockData } from '../../__mocks__';
import TopicQueryResult from '../TopicQueryResult';

// mock Trans tag from lingui
vi.mock('@lingui/macro', () => ({
  t: (str) => str,
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

  it('should render with the first result data', () => {
    renderComponent();
    expect(screen.getByText('text/plain')).toBeInTheDocument();
    expect(screen.getByText('Message v1.0.0')).toBeInTheDocument();
    expect(screen.getByText('hello world')).toBeInTheDocument();

    expect(screen.getByText('1 results of 10 total')).toBeInTheDocument();

    expect(screen.getByTestId('prev-query-btn')).toBeDisabled();

    expect(screen.getByTestId('next-query-btn')).toBeEnabled();
  });

  it(' should display the correct result when click on next button', () => {
    renderComponent();
    // click on next button
    fireEvent.click(screen.getByTestId('next-query-btn'));
    expect(screen.getByText('text/csv')).toBeInTheDocument();
    expect(screen.getByText('Spreadsheet v1.1.0')).toBeInTheDocument();
    expect(screen.getByText('hello,world')).toBeInTheDocument();

    expect(screen.getByText('2 results of 10 total')).toBeInTheDocument();

    expect(screen.getByTestId('prev-query-btn')).toBeEnabled();

    expect(screen.getByTestId('next-query-btn')).toBeEnabled();
  });
});

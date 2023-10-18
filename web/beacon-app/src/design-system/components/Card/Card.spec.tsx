import { render } from '@testing-library/react';
import React from 'react';

import Card from './Card';

describe('<Card>', () => {
  it('should display the children', () => {
    const { getByText } = render(<Card>Card</Card>);
    expect(getByText('Card')).toBeInTheDocument();
  });

  it('should apply containerClasses props', () => {
    const { getByText } = render(
      <Card containerClasses="bg-white shadow-lg rounded-lg">Card</Card>
    );
    expect(getByText('Card')).toHaveClass('bg-white shadow-lg rounded-lg');
  });

  it('should apply contentClasses props  ', () => {
    const { getByText } = render(<Card contentClasses="p-4">Card</Card>);
    expect(getByText('Card')).toHaveClass('p-4');
  });
});

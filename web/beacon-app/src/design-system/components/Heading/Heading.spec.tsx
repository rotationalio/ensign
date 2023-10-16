import { render } from '@testing-library/react';

import Heading from './Heading';

describe('Heading component', () => {
  it('renders correctly with default level of 1 and as prop of h3', () => {
    const { getByText } = render(<Heading>Heading Content</Heading>);
    expect(getByText('Heading Content').tagName).toBe('H3');
  });

  it('renders correctly with custom level and as prop', () => {
    const { getByText } = render(
      <Heading as="h2" level={2}>
        Heading Content
      </Heading>
    );
    expect(getByText('Heading Content').tagName).toBe('H2');
  });
});

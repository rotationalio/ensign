import { render } from '@testing-library/react';

import Select from './Select';

describe('Select', () => {
  it('should render successfully with options', () => {
    const { baseElement } = render(
      <Select
        options={[
          { value: '1', label: 'Option 1' },
          { value: '2', label: 'Option 2' },
          { value: '3', label: 'Option 3' },
        ]}
      />
    );
    expect(baseElement).toBeTruthy();
  });

  it('should render successfully with selected option', () => {
    const { baseElement } = render(
      <Select
        options={[
          { value: '1', label: 'Option 1' },
          { value: '2', label: 'Option 2' },
          { value: '3', label: 'Option 3' },
        ]}
        selectedOption="2"
      />
    );
    expect(baseElement).toBeTruthy();
  });
});

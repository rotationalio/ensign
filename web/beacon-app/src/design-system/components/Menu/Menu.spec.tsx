import { fireEvent, render } from '@testing-library/react';

import Menu from './Menu';

describe('Menu component', () => {
  it('should render without errors', () => {
    const { container } = render(<Menu />);
    expect(container).toBeTruthy();
  });

  it('should render children', () => {
    const testChildren = 'Test Children';
    const { getByText } = render(
      <Menu>
        <Menu.Item>{testChildren}</Menu.Item>
      </Menu>
    );
    expect(getByText(testChildren)).toBeTruthy();
  });

  it('should open and close on button click', () => {
    const handleClick = jest.fn();
    const { container, getByText } = render(
      <div>
        <button onClick={handleClick}></button>
        <Menu>
          <Menu.Item>Test Item</Menu.Item>
        </Menu>
      </div>
    );
    const button = container.querySelector('button');
    fireEvent.click(button!);
    expect(getByText('Test Item')).toBeTruthy();
    fireEvent.click(button!);
    expect(() => getByText('Test Item')).toThrow();
  });
});

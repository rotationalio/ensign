import { render, screen } from '@testing-library/react';
import React from 'react';

import Modal from './Modal';

describe('Modal', () => {
  it('renders the title', () => {
    const title = 'Modal Title';
    render(
      <Modal title={title} open={true}>
        <p>Modal Content</p>
      </Modal>
    );
    expect(screen.getByText(title)).toBeInTheDocument();
  });

  it('renders the children', () => {
    const children = <p>Modal Content</p>;
    render(<Modal open>{children}</Modal>);
    expect(screen.getByText('Modal Content')).toBeInTheDocument();
  });

  it('does not render the close button when isDismissible is false', () => {
    render(
      <Modal open>
        <p>Modal Content</p>
      </Modal>
    );
    expect(screen.queryByRole('button', { name: 'Close' })).toBeNull();
  });

  /* it('renders the container with the correct className', () => {
    const containerClassName = 'bg-red-500';
    render(
      <Modal open containerClassName={containerClassName}>
        <p>Modal Content</p>
      </Modal>
    );
    expect(screen.getByTestId('container')).toHaveClass(containerClassName);
  });

  it('renders the backdrop when fullScreen is true', () => {
    render(
      <Modal open fullScreen={true}>
        <p>Modal Content</p>
      </Modal>
    );
    expect(screen.getByTestId('backdrop')).toBeInTheDocument();
  });

  it('does not render the backdrop when fullScreen is false', () => {
    render(
      <Modal open fullScreen={false}>
        <p>Modal Content</p>
      </Modal>
    );
    expect(screen.queryByTestId('backdrop')).toBeNull();
  });

  it('renders the modal with the correct size', () => {
    const size = 'large';
    render(
      <Modal open size={size}>
        <p>Modal Content</p>
      </Modal>
    );
    expect(screen.getByTestId('modal')).toHaveClass(`modal-${size}`);
  });

  it('renders the modal with the default size if size prop is not provided', () => {
    render(
      <Modal open>
        <p>Modal Content</p>
      </Modal>
    );
    expect(screen.getByTestId('modal')).toHaveClass(`modal-small`);
  });

  it('renders the modalTitle with correct props', () => {
    const titleProps: ModalTitleProps = {
      children: 'Modal Title',
    };
    render(
      <Modal open titleProps={titleProps}>
        <p>Modal Content</p>
      </Modal>
    );
    expect(screen.getByText(titleProps.children as string)).toBeInTheDocument();
  }); */
});

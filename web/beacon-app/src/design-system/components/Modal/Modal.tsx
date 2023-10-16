import React from 'react';

import { mergeClassnames } from '../../utils';
import { Button } from '../Button';
import CloseIcon from './CloseIcon';
import { Backdrop, Container, StyledModal, Title } from './Modal.styles';
import { ModalProps } from './Modal.types';
function Modal(props: ModalProps, ref: React.ForwardedRef<HTMLDivElement>) {
  const {
    slots,
    children,
    title,
    containerClassName,
    modalCloseBtnClassName,
    fullScreen,
    size = 'small',
    titleProps,
    onClose,
    ...rest
  } = props;
  return (
    <StyledModal slots={{ backdrop: Backdrop, ...slots }} {...rest} ref={ref}>
      <Container
        size={size}
        fullScreen={fullScreen}
        className={mergeClassnames(
          'no-scrollbar w-[25vw] max-w-[80vw] overflow-scroll lg:max-w-[50vw]',
          containerClassName
        )}
      >
        {onClose && (
          <Button
            variant="ghost"
            size="custom"
            className={
              mergeClassnames(
                'bg-transparent absolute top-0 right-4 m-4 h-4 w-4 border-none py-0',
                modalCloseBtnClassName
              ) as string
            }
          >
            <CloseIcon onClick={onClose} />
          </Button>
        )}

        {title && (
          <Title {...titleProps} ref={null}>
            {title}
          </Title>
        )}
        <div>{children}</div>
      </Container>
    </StyledModal>
  );
}

export type ModalTitleProps = {
  children: React.ReactNode;
};

export default React.forwardRef(Modal);

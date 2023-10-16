import React from 'react';

import mergeClassnames from '../../utils/mergeClassnames';
import { CardBody, CardHeader, CardWrapper } from './Card.styles';
import { CardComposition, CardProps } from './Card.types';
import { getImgWrapperStyle } from './utils';

const hasHeaderElement = (child: React.ReactNode): child is React.ReactElement =>
  React.isValidElement(child) && child.type === CardHeader;

// const isBodyElement = (child: React.ReactNode): child is React.ReactElement =>
//   React.isValidElement(child) && child.type === CardBody;
type Props = {
  contentClassName?: string;
} & CardProps &
  CardComposition;

function Card<T>(props: Props & T) {
  const {
    as: Component,
    children,
    title,
    imgSrc,
    className,
    imgAlt,
    imgPosition,
    imgClassNames,
    contentClassName,
    headerClassNames,
    ...rest
  } = props;

  const hasHeaderChildren = React.Children.toArray(children).some((child) =>
    hasHeaderElement(child)
  );

  let header = null;

  if (!hasHeaderChildren && title) {
    header = <CardHeader className={headerClassNames}>{title}</CardHeader>;
  }

  const styleClasses = imgSrc && getImgWrapperStyle(imgPosition || 'top');

  return (
    <CardWrapper
      as={Component}
      className={mergeClassnames(
        'border-sm shadow-sm padding-4 overflow-hidden rounded-sm border bg-white',
        styleClasses,
        className
      )}
      {...rest}
    >
      {imgSrc && (
        <div>
          <img src={imgSrc} alt={imgAlt || ''} className={imgClassNames} />
        </div>
      )}

      <div className={contentClassName}>
        {header}
        {children}
      </div>
    </CardWrapper>
  );
}

Card.Header = CardHeader;
Card.Body = CardBody;

Card.displayName = 'Card';

export default Card;

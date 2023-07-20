import { mergeClassnames } from '@rotational/beacon-core';
import { useCallback } from 'react';

import QuickViewCard from './QuickViewCard';

export interface QuickViewProps {
  data: any;
  headers: string[]; // this implementation is not great cz if the order of the data changes, the headers will be wrong , we need to ask the backend team to add key to the data
}
const BRAND_COLORS = ['#ECF6FF', '#E5ECF6', '#ECF6FF', '#E5ECF6'];

const QUICKVIEW_CARD_LENGTH = 4;

const QuickView: React.FC<any> = ({ data, headers, ...props }) => {
  // TODO: create an abstraction for this logic
  const getValidData = useCallback(() => {
    const isDataValid = data?.length >= QUICKVIEW_CARD_LENGTH;
    if (data && !isDataValid) {
      return [];
    }
    return data.slice(0, QUICKVIEW_CARD_LENGTH) as IStats[];
  }, [data]);

  return (
    <div
      className={mergeClassnames(
        'grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-4',
        props?.className
      )}
    >
      {getValidData().map((item, index) => (
        <QuickViewCard
          key={index}
          title={headers ? headers[index] : item.name}
          color={BRAND_COLORS[index]}
        >
          {item.value} {item.units ? item.units : ''}
        </QuickViewCard>
      ))}
    </div>
  );
};

export default QuickView;

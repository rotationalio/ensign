import { useCallback } from 'react';

import { capitalize } from '@/utils/strings';

import QuickViewCard from './QuickViewCard';
interface QuickViewData {
  name: string;
  value: number;
}
export interface QuickViewProps {
  data: any;
}
const BRAND_COLORS = ['#ECF6FF', '#E5ECF6', '#ECF6FF', '#E5ECF6'];

/* we should have 4 statistic cards in the quick view
 * if we have less than 4 cards, we should not render the quick view
 * if we have more than 4 cards, we should only render the first 4 cards
 * The reason is that we have 4 colors for the quick view cards,
 * if we have more than 4 cards, we will have to use the same color for the cards
 * which is not a good user experience
 */
const QUICKVIEW_CARD_LENGTH = 4;
const defaultData: QuickViewData[] = [
  {
    name: 'Active Projects',
    value: 0,
  },
  {
    name: 'Topics',
    value: 0,
  },
  {
    name: 'API Keys',
    value: 0,
  },
  {
    name: 'Data Storage',
    value: 0,
  },
];

const QuickView: React.FC<any> = ({ data }) => {
  // TODO: create an abstraction for this logic
  const getValidData = useCallback(() => {
    const isDataValid = data?.length >= QUICKVIEW_CARD_LENGTH;
    if (data && !isDataValid) {
      return [];
    }
    if (!data) {
      return defaultData;
    }
    return data.slice(0, QUICKVIEW_CARD_LENGTH) as QuickViewData[];
  }, [data]);

  return (
    <div className="grid grid-cols-2 gap-y-10 gap-x-5 lg:grid-cols-4" data-cy="quick_view">
      {getValidData().map((item, index) => (
        <QuickViewCard key={item.name} title={capitalize(item.name)} color={BRAND_COLORS[index]}>
          {item.value}
        </QuickViewCard>
      ))}
    </div>
  );
};

export default QuickView;

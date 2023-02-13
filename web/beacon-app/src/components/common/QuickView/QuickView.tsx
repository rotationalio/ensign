import { useCallback } from 'react';

import QuickViewCard from './QuickViewCard';

interface QuickViewData {
  name: string;
  value: number;
}
export interface QuickViewProps {
  data: QuickViewData[];
}
const BRAND_COLORS = ['#ECF6FF', '#FFE9DD', '#ECFADC', '#FBF8EC'];

/* we should have 4 statistic cards in the quick view
 * if we have less than 4 cards, we should not render the quick view
 * if we have more than 4 cards, we should only render the first 4 cards
 * The reason is that we have 4 colors for the quick view cards,
 * if we have more than 4 cards, we will have to use the same color for the cards
 * which is not a good user experience
 */
const QUICKVIEW_CARD_LENGTH = 4;

const QuickView: React.FC<QuickViewProps> = ({ data }) => {
  const getValidData = useCallback(() => {
    const isDataValid = data?.length >= QUICKVIEW_CARD_LENGTH;
    if (!isDataValid) {
      return [];
    }
    return data.slice(0, QUICKVIEW_CARD_LENGTH);
  }, [data]);

  return (
    <div className="grid grid-cols-2 gap-y-10 gap-x-20 lg:grid-cols-4">
      {getValidData().map((item, index) => (
        <QuickViewCard key={item.name} title={item.name} color={BRAND_COLORS[index]}>
          {item.value}
        </QuickViewCard>
      ))}
    </div>
  );
};

export default QuickView;

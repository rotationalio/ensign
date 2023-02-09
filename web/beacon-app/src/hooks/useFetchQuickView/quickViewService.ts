import { QuickViewProps } from '@/components/common/QuickView/QuickView';
export type QuickViewKey = 'project' | 'tenant';

export interface QuickViewDTO {
  id: string;
  key: QuickViewKey;
}

export interface QuickViewResponse {
  data: QuickViewProps;
}

export interface QuickViewQuery {
  getQuickView: () => void;
  hasQuickViewFailed: boolean;
  isFetchingQuickView: boolean;
  quickView: any;
  wasQuickViewFetched: boolean;
  error: any;
}

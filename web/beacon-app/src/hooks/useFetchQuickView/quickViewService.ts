export type QuickViewKey = 'project' | 'tenant';

export interface QuickViewData {
  name: string;
  value: number;
  units?: string;
  percent?: number;
}
export interface QuickViewDTO {
  id: string;
  key: QuickViewKey;
}

export interface QuickViewResponse {
  data: QuickViewData[];
}

export interface QuickViewQuery {
  getQuickView: () => void;
  hasQuickViewFailed: boolean;
  isFetchingQuickView: boolean;
  quickView: any;
  wasQuickViewFetched: boolean;
  error: any;
}

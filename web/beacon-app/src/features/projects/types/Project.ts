import { MemberResponse } from '@/features/members/types/memberServices';
import { QuickViewData } from '@/hooks/useFetchQuickView/quickViewService';
export type Project = {
  created: string;
  id: string;
  modified: string;
  name: string;
  tenant_id: string;
  description?: string;
  status?: string;
  owner: Partial<MemberResponse>;
  active_topics?: number;
  data_storage?: QuickViewData;
};

export enum ProjectStatus {
  ACTIVE = 'ACTIVE',
  INACTIVE = 'INACTIVE',
  DELETED = 'DELETED',
  INCOMPLETE = 'Incomplete',
  COMPLETE = 'Complete',
}

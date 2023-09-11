import { t } from '@lingui/macro';
import invariant from 'invariant';
import { useEffect } from 'react';

import { useFetchTenants } from '@/features/tenants/hooks/useFetchTenants';

// import { useOrgStore } from '@/store';
import { MainStyle } from './DashLayout.styles';
import MobileFooter from './MobileFooter';
import { Sidebar } from './Sidebar';

type DashLayoutProps = {
  children?: React.ReactNode;
};

const DashLayout: React.FC<DashLayoutProps> = ({ children }) => {
  // const Store = useOrgStore.getState() as any;
  const { tenants, wasTenantsFetched } = useFetchTenants();
  const hasTenants =
    tenants?.tenants && Array.isArray(tenants?.tenants) && tenants?.tenants?.length > 0;
  // const tenantID = tenants?.tenants[0]?.id;

  // ensure the tenant is loaded correctly otherwise we need to throw an error
  useEffect(() => {
    if (wasTenantsFetched) {
      invariant(
        hasTenants,
        t`Tenant is not loaded correctly. Please contact us at support@rotational.io for assistance.`
      );
    }
  }, [hasTenants, wasTenantsFetched]);

  // useEffect(() => {
  //   if (wasTenantsFetched && tenantID) {
  //     Store.setTenantID(tenantID);
  //   }
  // }, [tenantID, wasTenantsFetched, Store]);

  return (
    <div className="flex flex-col md:pl-[250px]">
      <Sidebar className="hidden md:block" />
      <MainStyle>{children}</MainStyle>
      <MobileFooter />
    </div>
  );
};

export default DashLayout;

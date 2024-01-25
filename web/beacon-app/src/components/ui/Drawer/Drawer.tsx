import * as RadixDialog from '@radix-ui/react-dialog';
import { useEffect } from 'react';
import { useLocation } from 'react-router-dom';

import { Sidebar } from '@/components/layout/Sidebar';
import useDrawer from '@/hooks/useDrawer';

const Dialog = () => {
  const { isOpen, closeDrawer } = useDrawer();
  const location = useLocation();

  useEffect(() => {
    if (isOpen) {
      closeDrawer();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [location?.pathname]);

  return (
    <RadixDialog.Root open={isOpen} onOpenChange={closeDrawer}>
      <RadixDialog.Portal>
        <RadixDialog.Overlay className="data-[state=open]:animate-overlayShow fixed inset-0" />
        <RadixDialog.Content className="data-[state=open]:animate-contentShow fixed left-0 top-0 h-screen w-[250px] bg-white p-[25px] shadow-[hsl(206_22%_7%_/_35%)_0px_10px_38px_-10px,_hsl(206_22%_7%_/_20%)_0px_10px_20px_-15px] focus:outline-none">
          <Sidebar className="w-[250px]" />
        </RadixDialog.Content>
      </RadixDialog.Portal>
    </RadixDialog.Root>
  );
};

export default Dialog;

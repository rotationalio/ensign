import { AriaButton } from '@rotational/beacon-core';
import { useState } from 'react';

import { CancelAcctModal } from '@/components/ui/CancelModal';

export default function CancelAccount(props: any) {
  const [showModal, setShowModal] = useState(false);

  const handleOpen = () => setShowModal(true);
  return (
    <AriaButton variant="tertiary" className="rounded-sm" onClick={handleOpen}>
      Cancel Account
      {showModal && <CancelAcctModal close={props.close} />}
    </AriaButton>
  );
}

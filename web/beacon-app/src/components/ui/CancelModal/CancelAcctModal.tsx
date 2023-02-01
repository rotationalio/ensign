import { memo, useState } from "react"
import { Fragment } from "react"

import { Modal } from "@rotational/beacon-core"

import { Close } from "@/components/icons/close"

function CancelAcctModal() {
  const [isOpen, setIsOpen] = useState(true)

  const closeModal = () => setIsOpen(false)

    return(
    <div>
      <Modal
        title="Cancel Account"
        open={isOpen}
        className="max-w-md"
      >
        <Fragment key=".0">
          <Close onClick={closeModal} className="absolute top-4 right-8"></Close>
          <p className="pb-4">Please contact us at <span className="font-bold">support@rotational.io</span> to cancel your account. Please include your name, email, and Org ID in your request to cancel your account. We are working on an automated process to cancel accounts and appreciate your patience.</p>
          <p className="pb-4">You are the <span className="font-bold">Owner</span> of this account. If you cancel your account, your Organization, Tenant, Project, and Topic and all associated data will be <span className="font-bold">permanently</span> deleted.</p>
        </Fragment>
      </Modal>
     </div>
    )
}

export default memo(CancelAcctModal)
import { memo } from "react"
import { Fragment } from "react"

import { AriaButton, Modal } from "@rotational/beacon-core"

function CancelAcctModal() {
    return(
    <div>
      <Modal
        title="Cancel Account"
        open={true}
        className="max-w-md"
        >
          <Fragment key=".0">
            <p className="pb-4">Please contact us at <span className="font-bold">support@rotational.io</span> to cancel your account. Please include your name, email, and Org ID in your request to cancel your account. We are working on an automated process to cancel accounts and appreciate your patience.</p>
            <p className="pb-4">You are the <span className="font-bold">Owner</span> of this account. If you cancel your account, your Organization, Tenant, Project, and Topic and all associated data will be <span className="font-bold">permanently</span> deleted.</p>
          </Fragment>
      </Modal>
     </div>
    )
}

export default memo(CancelAcctModal)
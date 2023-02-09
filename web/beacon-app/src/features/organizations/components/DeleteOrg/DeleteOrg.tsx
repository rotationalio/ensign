import { AriaButton } from "@rotational/beacon-core";
import { useState } from "react";
import DeleteOrgModal from "../DeleteOrgModal/DeleteOrgModal";

export default function DeleteOrg(props: any) {
  const [showModal, setShowModal] = useState(false)

  const handleOpen = () => setShowModal(true)

    return(
      <div>
      <AriaButton
      variant='tertiary'
      className="rounded-sm"
      onClick={handleOpen}
        >
          Delete Org
            {showModal && <DeleteOrgModal close={props.close} />}
        </AriaButton>
      </div>
    )
}
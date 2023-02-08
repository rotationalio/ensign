import { BlueBars } from "@/components/icons/blueBars";
import { useState } from "react";
import { DeleteOrg } from "../DeleteOrg";

export default function OrganizationDetails() {
    const [showButton, setShowButton] = useState(false)

    const handleOpen = () => setShowButton(true)
    const handleClose = () => setShowButton(false)

    return(
    <>
        <h3 className="mt-10 font-bold text-2xl">Organization Dashboard</h3>
        <h4 className="mt-10 pt-4 border-t border-primary-900 max-w-4xl font-bold text-xl">Organization Details</h4>
        <section className="mt-8 pl-6 max-w-4xl border-2 border-secondary-500 rounded-md">
            <div className="mt-4 absolute right-28">
                <BlueBars onClick={handleOpen} />
                <div className="relative left-12">
                    {showButton && <DeleteOrg close={handleClose} />}
                </div>    
            </div>
            <div className="flex py-8 gap-4">
                <h6 className="font-bold">Organization Name:</h6>
                <span>Name</span>
            </div>
            <div className="flex pb-8 gap-36">
                <h6 className="font-bold">URL:</h6>
                <span>Domain</span>
            </div>
            <div className="flex pb-8 gap-32">
                <h6 className="font-bold">Org ID:</h6>
                <span>ID</span>
            </div>
            <div className="flex pb-8 gap-28">
                <h6 className="font-bold">Owner:</h6>
                <span className="ml-3">Owner</span>
            </div>
            <div className="flex pb-8 gap-28">
                <h6 className="font-bold">Created:</h6>
                <span className="ml-1">Created</span>
            </div>
        </section>
    </>
    )
}
export default function OrganizationDetails() {
    return(
    <>
        <h3 className="mt-10">Organization Dashboard</h3>
        <h4 className="mt-10 pt-4 border-t max-w-4xl">Organization Details</h4>
        <section className="mt-8 pl-6 max-w-4xl border border-primary-900 rounded-md">
            <div className="flex py-6">
                <h6>Name:</h6>
                <span className="ml-12">Name</span>
            </div>
            <div className="flex pb-6">
                <h6>URL:</h6>
                <span className="ml-16">Domain</span>
            </div>
            <div className="flex pb-6">
                <h6>ID:</h6>
                <span className="ml-20">ID</span>
            </div>
            <div className="flex pb-6">
                <h6>Owner:</h6>
                <span className="ml-11">Owner</span>
            </div>
            <div className="flex pb-6">
                <h6>Created:</h6>
                <span className="ml-8">Created</span>
            </div>
        </section>
    </>
    )
}
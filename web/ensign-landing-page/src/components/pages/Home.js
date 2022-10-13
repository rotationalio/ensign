import React from "react";
import Header from '../layout/Header'
import Main from '../layout/Main';
import BuildApps from '../layout/BuildApps';
import Diagram from '../layout/Diagram';
import DevExperience from '../layout/DevExperience';
import Footer from '../layout/Footer'
import withTracker from "../../lib/analytics";

const Home = () => {
    return (
        <>
          <Header />
          <Main />
          <Diagram />
          <BuildApps />
          <DevExperience />
         <div className="bg-[#ECF6FF]">
            <Footer />
         </div>
        </>
      );
}

export default withTracker(Home); 

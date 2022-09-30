import './App.css';
import Header from './Components/layout/Header'
import Main from './Components/layout/Main';
import BuildApps from './Components/layout/BuildApps';
import Diagram from './Components/layout/Diagram';
import DevExperience from './Components/layout/DevExperience';
import Footer from './Components/layout/Footer'

function App() {
  return (
    <div className="App">
      <Header />
      <Main />
      <Diagram />
     <BuildApps />
     <DevExperience />
     <Footer />
    </div>
  );
}

export default App;

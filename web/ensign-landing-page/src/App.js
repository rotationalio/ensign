import {BrowserRouter, Routes, Route} from 'react-router-dom'
import Home from './components/pages/Home';
import Confirmation from './components/pages/Confirmation'

function App() {
  return (
    <BrowserRouter>
    <Routes>
      <Route path="/" element={<Home />} />
     <Route path="/ensign-access" element={<Confirmation />}/>
    </Routes>
    </BrowserRouter>
  );
}

export default App;

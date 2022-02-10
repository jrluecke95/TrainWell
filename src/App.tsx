import logo from './logo.svg';
import './App.css';
import { BrowserRouter } from 'react-router-dom';
import Appbar from './Components/Appbar';

function App() {
  return (
    <div className="App">
      <BrowserRouter>
        <Appbar />
        
      </BrowserRouter>
      
    </div>
  );
}

export default App;

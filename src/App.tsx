import { BrowserRouter, Route, Routes } from 'react-router-dom';
import Appbar from './Components/Appbar';
import ClientLoginPage from './Pages/ClientLoginPage';
import ClientRegisterPage from './Pages/ClientRegisterPage';
import CoachLoginPage from './Pages/CoachLoginPage';
import CoachRegisterPage from './Pages/CoachRegisterPage';
import HomePage from './Pages/HomePage';

function App() {
  return (
    <div className="App">
      <BrowserRouter>
        <Appbar />
        <Routes>
          <Route path='/' element={<HomePage />} />
          <Route path='/coach/register' element={<CoachRegisterPage />} />
          <Route path='/coach/login' element={<CoachLoginPage />} />
          <Route path='/client/register' element={<ClientRegisterPage />} />
          <Route path='/client/login' element={<ClientLoginPage />} />
        </Routes> 
      </BrowserRouter>
    </div>
  );
}

export default App;

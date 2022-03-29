import { BrowserRouter, Route, Routes } from 'react-router-dom';
import Appbar from './Components/Appbar';
import ClientLoginPage from './Pages/ClientLoginPage';
import ClientRegisterPage from './Pages/ClientRegisterPage';
import CoachLoginPage from './Pages/CoachLoginPage';
import CoachRegisterPage from './Pages/CoachRegisterPage';
import CoachWorkoutPage from './Pages/CoachWorkoutPage';
import CoachWorkoutsPage from './Pages/CoachWorkoutsPage';
import CreateWorkoutPage from './Pages/CreateWorkoutPage';
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
          <Route path='/coach/workoutPlan/create' element={ <CreateWorkoutPage /> }/>
          <Route path='/coach/workoutPlans' element={ <CoachWorkoutsPage />} />
          <Route path='coach/workoutPlan/:id' element={ <CoachWorkoutPage/> } />
        </Routes> 
      </BrowserRouter>
    </div>
  );
}

export default App;

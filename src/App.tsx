import { BrowserRouter, Route, Routes } from 'react-router-dom';
import Appbar from './Components/Appbar';
import ClientLoginPage from './Pages/ClientLoginPage';
import ClientRegisterPage from './Pages/ClientRegisterPage';
import CoachLoginPage from './Pages/CoachLoginPage';
import CoachRegisterPage from './Pages/CoachRegisterPage';
import { CoachWorkoutPage } from './Pages/CoachWorkoutPage';
import CoachWorkoutPlanPage from './Pages/CoachWorkoutPlanPage';
import CoachWorkoutPlansPage from './Pages/CoachWorkoutPlansPage';
import CreateWorkoutPlanPage from './Pages/CreateWorkoutPlanPage';
import HomePage from './Pages/HomePage';

function App() {
  return (
    <div className="App">
      <BrowserRouter>
        <Appbar />
        <Routes>
          <Route path='/' element={<HomePage />} />
          <Route path='/coach/register' element={<CoachRegisterPage/>} />
          <Route path='/coach/login' element={<CoachLoginPage/>} />
          <Route path='/client/register' element={<ClientRegisterPage/>} />
          <Route path='/client/login' element={<ClientLoginPage/>} />
          <Route path='/coach/workoutPlan/create' element={ <CreateWorkoutPlanPage/> }/>
          <Route path='/coach/workoutPlans' element={ <CoachWorkoutPlansPage/>} />
          <Route path='coach/workoutPlan/:id' element={ <CoachWorkoutPlanPage/> } />
          <Route path='/coach/workout/details/:id' element={ <CoachWorkoutPage/> } />
        </Routes> 
      </BrowserRouter>
    </div>
  );
}

export default App;

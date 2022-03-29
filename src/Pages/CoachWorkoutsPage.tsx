import { Button } from "@mui/material";
import { useEffect, useState } from "react";
import { NavLink } from "react-router-dom";


export default function CoachorkoutsPage() {
  const [workoutPlans, setWorkoutPlans] = useState([]);

  useEffect(() => {
    fetch('/api/coach/workoutPlans')
    .then(res => res.json())
    .then(data => {
      setWorkoutPlans(data)
    })
  }, [])

  return(
    <>
      <h2>Coach Workouts</h2>
      <div>
        {workoutPlans.map((workout: any) => (
          <>
            <Button
            component={NavLink}
            to={`/coach/workoutPlan/${workout._id}`}><span>{workout.name}</span>
            </Button>
          </>
          
        ))}
      </div>
    </>
  )
}
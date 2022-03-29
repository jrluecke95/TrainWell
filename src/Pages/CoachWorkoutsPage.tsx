import { Button } from "@mui/material";
import { useEffect, useState } from "react";
import { NavLink } from "react-router-dom";


export default function CoachorkoutsPage() {
  const [workouts, setWorkouts] = useState([]);

  useEffect(() => {
    fetch('/api/coach/workouts')
    .then(res => res.json())
    .then(data => {
      console.log(data)
      setWorkouts(data)
    })
  }, [])

  return(
    <>
      <h2>Coach Workouts</h2>
      <div>
        {workouts.map((workout: any) => (
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
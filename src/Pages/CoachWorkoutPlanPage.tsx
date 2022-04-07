import { Box, Button, FormControl, InputLabel, MenuItem, Select, TextField } from "@mui/material";
import { useEffect, useState } from "react";
import { NavLink, useParams } from "react-router-dom";

export default function CoachWorkoutPlanPage() {
  // TODO this state doesn't make sense to have iniated to empty array since state is only one thing - however, it works for now so roll with it
  const [workouts, setWorkouts] = useState<any>([])
  const { id } = useParams()

  useEffect(() => {
    fetch(`/api/workoutPlan/${id}`)
    .then(res => res.json())
    .then(data => {
      setWorkouts(data)
    })
    
  }, [id])

  return(
    <>
      <h2>Workout</h2>
      <div>
        {workouts && workouts.map((workout: any, i: any) => {
          return (
            <>
              <Button
                component={NavLink}
                to={`/coach/workout/details/${workout._id}`}><span>{`Day ${i+1}`}</span>
            </Button>
            </>
          )
        })}
      </div>
    </>
  )
}
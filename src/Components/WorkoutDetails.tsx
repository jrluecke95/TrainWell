import { Box } from "@mui/material"
import { useEffect, useState } from "react"

export function WorkoutDetails(props: any) {
  const [workout, setWorkout] = useState<any>([])
  useEffect(() => {
    fetch(`/api/workout/details/${props.id}`, {
      method: 'GET',
    })
    .then(res => res.json())
    .then(data => {
      setWorkout(data)
    })
  }, [])

  return(

    <Box>
      <div>
        {workout.exercisesDetails?.map((details: any, i: any) => {
          return (
            <div key={workout.exercisesDetails._id}>
              <p>Exercise Number {i+1}</p>
              <p>Desc: {details.description}</p>
              <p>Sets: {details.sets}</p>
              <p>Reps: {details.reps}</p>
              <p>Weight: {details.weight}</p>
              <p>Name: {details.exercise.name}</p>
            </div>
          )
        })}
      </div>
      </Box>    
  ) 
}
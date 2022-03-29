import { Box, Button, MenuItem, TextField } from "@mui/material";
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

export default function CoachWorkoutPage() {
  // TODO this state doesn't make sense to have iniated to empty array since state is only one thing - however, it works for now so roll with it
  const [workouts, setWorkouts] = useState<any>([]);
  const [exercises, setExercises] = useState<any>([])
  const { id } = useParams();
  const [form, setForm] = useState({
    exercise: {
      id: '',
      name: ''
    },
    sets: '',
    reps: '',
    weight: '',
    description: ''
  })

  useEffect(() => {
    fetch(`/api/workoutPlan/${id}`)
    .then(res => res.json())
    .then(data => {
      setWorkouts(data)
    })

    fetch('/api/exercises')
    .then(res => res.json())
    .then(data => {
      setExercises(data)
    })
    
  }, [])

  const handleChange = (e: any) => {
    setForm({
        ...form,
        [e.target.name]: e.target.value
    })
  }

  const handleSubmit = (e: any) => {
    e.PreventDefault()
    const data = {
      exercise: {
        id: form.exercise.id,
        name: form.exercise.name
      },
      sets: form.sets,
      reps: form.reps,
      weight: form.weight,
      description: form.description
    }

    fetch('/api/workoutPlan/addNewWorkout', {
      headers: { 'Content-Type': 'application/json' },
      method: 'POST',
      body: JSON.stringify(data),
    })
    .then(res => {
    return res.json()
    })
    .then(data => {
        if (data.error) {
            alert(data.error)
        } else {
            alert('exercise added successfully')
            // probably need to refresh the page/component after
        }
    })
  }


  return(
    <>
      <h2>Workout</h2>
      <div>
        {workouts[0] && workouts[0].exercisesDetails.map((details: any, i: any) => {
          return (
            <>
              <p>Exercise Number {i+1}</p>
              <p>Desc: {details.description}</p>
              <p>Sets: {details.sets}</p>
              <p>Reps: {details.reps}</p>
              <p>Weight: {details.weight}</p>
              <p>Name: {details.exercise.name}</p>
            </>
          )
        })}
      </div>
      <Box sx={{ flexGros: 1 }}>
      <form noValidate autoComplete="off" onSubmit={handleSubmit}>
        <TextField id="filled-basic" label="Sets"name="sets" value={form.sets}/>
        <TextField id="filled-basic" label="Reps" variant="filled" onChange={handleChange} name="reps" value={form.reps}/>
        <TextField id="filled-basic" label="Weight" variant="filled" onChange={handleChange} name="weight" value={form.weight}/>
        <TextField id="filled-basic" label="Description" variant="filled" onChange={handleChange} name="description" value={form.description}/>
        <TextField id="standard-select" label="exercise">
          {console.log(exercises)}
          {exercises.map((exercise: any) => {
            return (
              <>
                <MenuItem key={exercise._id} value={exercise.name}>{exercise.name}</MenuItem>
              </>
            )
          })}
          
        </TextField>
        <Button variant="contained" type="submit">Submit</Button>
      </form>      
      </Box>
    </>
  )
}
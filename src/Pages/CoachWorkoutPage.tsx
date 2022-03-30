import { Box, Button, FormControl, InputLabel, MenuItem, Select, TextField } from "@mui/material";
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

export default function CoachWorkoutPage() {
  // TODO this state doesn't make sense to have iniated to empty array since state is only one thing - however, it works for now so roll with it
  const [workouts, setWorkouts] = useState<any>([]);
  const [exercises, setExercises] = useState<any>([])
  const { id } = useParams();
  const [form, setForm] = useState({
    exerciseId: '',
    exerciseName: '',
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
    
  }, [id])
// TODO cant get this to stop fetching on form change
  useEffect(() => {
    fetch('/api/exercises')
    .then(res => res.json())
    .then(data => {
      setExercises(data)
    })
  }, [id])

  const handleChange = (e: any) => {
    setForm({
        ...form,
        [e.target.name]: e.target.value
    })
  }

  // TODO need to restructure wokrout plan vs workout page
  // need a workout plan page that displays workouts
  // then this page should go on each workout 
  // right now this page serves as both and it is fucked

  const handleSubmit = (e: any) => {
    console.log('this is working')
    e.preventDefault()
    const data = {
      workoutID: id,
      exercise: {
        id: form.exerciseId,
        name: form.exerciseName
      },
      sets: form.sets,
      reps: form.reps,
      weight: form.weight,
      description: form.description
    }
    console.log(data)

    fetch('/api/workout/exercise', {
      headers: { 'Content-Type': 'application/json' },
      method: 'POST',
      body: JSON.stringify(data),
    })
    .then(res => {
      console.log(res)
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
        <TextField onChange={handleChange} id="filled-basic" label="Sets" name="sets" value={form.sets}/>
        <TextField id="filled-basic" label="Reps" variant="filled" onChange={handleChange} name="reps" value={form.reps}/>
        <TextField id="filled-basic" label="Weight" variant="filled" onChange={handleChange} name="weight" value={form.weight}/>
        <TextField id="filled-basic" label="Description" variant="filled" onChange={handleChange} name="description" value={form.description}/>
        <FormControl fullWidth>
          <InputLabel>Exercise</InputLabel>
          <Select
            value={form.exerciseName}
            label="Exercise"
            name="exerciseName"
            onChange={handleChange}
          >
            
            {exercises.map((exercise: any) => {
            return (
                <MenuItem key={exercise._id} value={exercise.name}>{exercise.name}</MenuItem>
              )
            })}
          </Select>
        </FormControl>
        <Button variant="contained" type="submit">Submit</Button>
      </form>      
      </Box>
    </>
  )
}
import { Box, Button, FormControl, InputLabel, MenuItem, Select, TextField } from "@mui/material";
import { Suspense, useEffect, useState } from "react"
import { useParams } from "react-router-dom";
import { WorkoutDetails } from "../Components/WorkoutDetails";

export function CoachWorkoutPage() {
  const { id } = useParams();
  const [exercises, setExercises] = useState<any>([])
  
  const [form, setForm] = useState({
    exerciseId: '',
    exerciseName: '',
    sets: '',
    reps: '',
    weight: '',
    description: ''
  })
  // TODO cant get this to stop fetching on form change
  useEffect(() => {
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
    console.log(form)
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
    fetch('/api/workout/exercise', {
      headers: { 'Content-Type': 'application/json' },
      method: 'POST',
      body: JSON.stringify(data),
    })
    .then(res => {
      console.log(JSON.stringify(data))
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
      <WorkoutDetails id={id}/>
      
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
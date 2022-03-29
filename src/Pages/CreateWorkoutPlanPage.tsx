import { Box, Button, TextField } from "@mui/material"
import { useState } from "react"

export default function CreateWorkoutPage() {
  const [form, setForm] = useState({
    workoutPlanName: ''
  })

  const handleChange = (e: any) => {
    setForm({
      ...form,
      [e.target.name]: e.target.value
    })
  }

  const handleSubmit = (e: any) => {
    e.preventDefault()
    const data = {
      name: form.workoutPlanName
    }
    
    fetch('/api/workoutPlan', {
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
            alert('successfully created workout plan')
            // woudl ideally redifrect here - figure out that part later
        }
    });
  }

  return (
    <>
      <Box sx={{ flexGros: 1 }}>
      <form noValidate autoComplete="off" onSubmit={handleSubmit}>
        <TextField id="filled-basic" label="WorkoutPlanName" variant="filled" onChange={handleChange} name="workoutPlanName" value={form.workoutPlanName}/>
        <Button variant="contained" type="submit">Submit</Button>
      </form>      
    </Box>
    </>
  ) 
}
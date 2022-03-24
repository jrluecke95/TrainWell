import { Box, Button, TextField } from "@mui/material"
import { useState } from "react"


export default function CoachLoginPage() {
  const [form, setForm] = useState({
    email: '',
    password: ''
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
      email: form.email,
      password: form.password
    }
    
    fetch('/api/coach/login', {
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
            alert('coach logged in Successfully')
            // woudl ideally redifrect here - figure out that part later
        }
    });
  }
  return (
    <Box sx={{ flexGros: 1 }}>
      <form noValidate autoComplete="off" onSubmit={handleSubmit}>
        <TextField id="filled-basic" label="Email" variant="filled" onChange={handleChange} name="email" value={form.email}/>
        <TextField id="filled-basic" label="Password" variant="filled" onChange={handleChange} name="password" value={form.password}/>
        <Button variant="contained" type="submit">Submit</Button>
      </form>      
    </Box>
  )
}
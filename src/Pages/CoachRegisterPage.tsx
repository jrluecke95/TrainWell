import { Box, Button, TextField } from "@mui/material";
import { useState } from 'react'

export default function CoachRegisterPage() {
  const [form, setForm] = useState({
    firstName: '',
    lastName: '',
    email: '',
    phoneNumber: '',
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
      personalInfo: {
        firstName: form.firstName,
        lastName: form.lastName,
        email: form.email,
        phoneNumber: form.phoneNumber,
        password: form.password
      }
    }
    
    fetch('/api/coach/create', {
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
            alert('User registered Successfully')
            // woudl ideally redifrect here - figure out that part later
        }
    });
    }
  return (
    <Box sx={{ flexGros: 1 }}>
      <form noValidate autoComplete="off" onSubmit={handleSubmit}>
        <TextField id="filled-basic" label="First Name" variant="filled" onChange={handleChange} name="firstName" value={form.firstName}/>
        <TextField id="filled-basic" label="Last Name" variant="filled" onChange={handleChange} name="lastName" value={form.lastName}/>
        <TextField id="filled-basic" label="Email" variant="filled" onChange={handleChange} name="email" value={form.email}/>
        <TextField id="filled-basic" label="Phone No." variant="filled" onChange={handleChange} name="phoneNumber" value={form.phoneNumber}/>
        <TextField id="filled-basic" label="Password" variant="filled" onChange={handleChange} name="password" value={form.password}/>
        <Button variant="contained" type="submit">Submit</Button>
      </form>      
    </Box>
  )
}
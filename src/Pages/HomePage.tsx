import { Button } from "@mui/material";


export default function HomePage() {


  return (
    <>
      <div>Welcome to Trainwell</div>

      <Button color="inherit" href="/coach/register">Coach Register</Button>
      <Button color="inherit" href="/coach/login">Coach Login</Button>

      <Button color="inherit" href="/client/register">Client Register</Button>
      <Button color="inherit" href="/client/login">Client Login</Button>

    
    </>
  )
}
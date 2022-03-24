import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import type { RootState } from '../Redux/store'

// Define a type for the slice state
export interface UserState {
  firstName: string
  lastName: string
  email: string
  phoneNumber: string
}

// Define the initial state using that type
const initialState: UserState = {
  firstName: '',
  lastName: '',
  email: '',
  phoneNumber: ''
}

export const userSlice = createSlice({
  name: 'user',
  // `createSlice` will infer the state type from the `initialState` argument
  initialState,
  reducers: {
    // signIn: state => {
    //   state.firstName = 
    // },
    // decrement: state => {
    //   state.value -= 1
    // },
    // Use the PayloadAction type to declare the contents of `action.payload`
    signIn: (state, action: PayloadAction<UserState>) => {
      state.firstName = action.payload.firstName
      state.lastName = action.payload.lastName
      state.email = action.payload.email
      state.phoneNumber = action.payload.phoneNumber
    },
    signOut: state => {
      state.firstName = ''
      state.lastName = ''
      state.email = ''
      state.phoneNumber = ''
    }
  }
})

export const { signIn, signOut } = userSlice.actions

// Other code such as selectors can use the imported `RootState` type
export const selectUser = (state: RootState) => state.users

export default userSlice.reducer
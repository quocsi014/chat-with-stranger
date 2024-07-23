import { createSlice } from "@reduxjs/toolkit";

const initialState = {
  name: "",
}

export const userSlice = createSlice({
  name: 'user',
  initialState,
  reducers: {
    updateName: (state, action)=>{
      state.name = action.payload.name
    }
  }
})

export const {updateName} = userSlice.actions

export default userSlice.reducer

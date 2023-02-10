import { combineReducers } from '@reduxjs/toolkit';

const rootReducer = combineReducers({
  // Add  reducers here\
});

export type RootState = ReturnType<typeof rootReducer>;

export default rootReducer;

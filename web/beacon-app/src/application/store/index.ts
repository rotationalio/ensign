import { Action, configureStore } from '@reduxjs/toolkit';
import { useDispatch } from 'react-redux';
import { PersistConfig, persistReducer, persistStore } from 'redux-persist';
import storage from 'redux-persist/lib/storage';
import thunk, { ThunkAction } from 'redux-thunk';

import { appConfig } from '../config';
import rootReducer, { RootState } from './rootReducer';
const persistConfig: PersistConfig<any> = {
  key: 'root',
  storage,
};

const persistedReducer = persistReducer(persistConfig, rootReducer);

const store = configureStore({
  reducer: persistedReducer,
  middleware: [thunk],
  devTools: appConfig.nodeENV !== 'production',
});
export const persistor = persistStore(store);

export const resetStore = async () => {
  await persistor.purge();
  store.dispatch({ type: 'RESET_STORE' });
  await persistor.flush();
};

export type AppDispatch = typeof store.dispatch;
export const useAppDispatch = () => useDispatch<AppDispatch>();
export type AppThunk = ThunkAction<void, RootState, unknown, Action<string>>;

export default store;

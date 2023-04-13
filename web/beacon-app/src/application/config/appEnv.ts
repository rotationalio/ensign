import appConfig from './appConfig';
export const isProdEnv = appConfig.nodeENV === 'production';
export const isDevEnv = appConfig.nodeENV === 'development';
export const isTestEnv = appConfig.nodeENV === 'test';
export const isStagingEnv = appConfig.nodeENV === 'staging';

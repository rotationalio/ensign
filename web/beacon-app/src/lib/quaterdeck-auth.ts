import axios from 'axios';

import { getValidApiError, getValidApiResponse } from '@/application/api/ApiService';
import { appConfig } from '@/application/config';
import { QDK_API_ROUTE } from '@/constants';
import { AuthUser, NewUserAccount, UserAuthResponse } from '@/features/auth';

class QuarterDeckAuth {
  readonly baseUrl: string = appConfig.quaterDeckApiUrl;
  readonly refreshTokenUrl: string = `${appConfig.quaterDeckApiUrl}/${QDK_API_ROUTE.REFRESH_TOKEN}`;
  readonly loginUrl: string = `${appConfig.quaterDeckApiUrl}/${QDK_API_ROUTE.LOGIN}`;
  readonly registerUrl: string = `${appConfig.quaterDeckApiUrl}/${QDK_API_ROUTE.REGISTER}`;
  axiosInstance: any; // Axios;
  constructor() {
    this.axiosInstance = axios.create({
      baseURL: this.baseUrl,
      headers: {
        'Content-Type': 'application/json',
      },
    });
  }
  public async login(data: AuthUser) {
    try {
      const response = await axios.post(`${this.loginUrl}`, {
        data: JSON.stringify(data),
      });
      return getValidApiResponse<UserAuthResponse>(response);
    } catch (error: any) {
      getValidApiError(error);
    }
  }
  public async register(data: NewUserAccount) {
    const response = await axios.post(`${this.registerUrl}`, {
      data: JSON.stringify(data),
    });
    return response;
  }

  public async refreshToken(refreshToken: string) {
    const response = await axios.post(`${this.refreshTokenUrl}`, {
      data: JSON.stringify({ refreshToken }),
    });
    return response;
  }
}

export default QuarterDeckAuth;

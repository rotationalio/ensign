import { t } from '@lingui/macro';

//TODO: Add more error messages here and improve this class to be more robust

export default class ErrorMessage {
  static readonly somethingWentWrong = t`Something went wrong. Please contact us at support@rotational.io for assistance.`;
  static readonly NETWORK_ERROR = t`No internet connection. Please check your internet connection and try again.`;
}

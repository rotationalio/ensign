import { Container } from '@rotational/beacon-core';

import ForgotPasswordForm from '../components/ForgotPassword/ForgotPasswordForm';

const ForgotPasswordPage = () => {
  const submitFormHandler = (values: any) => {
    console.log(values);
  };
  return (
    <>
      <Container className="my-20">
        <div className="mx-auto min-h-min max-w-xl rounded-lg border border-solid border-primary-800 p-12">
          <div className="">
            <p className="mb-4">
              Lost your password? No problem. Enter your email address to recover your login
              credentials.
            </p>
            <ForgotPasswordForm onSubmit={submitFormHandler} />
          </div>
        </div>
      </Container>
    </>
  );
};

export default ForgotPasswordPage;

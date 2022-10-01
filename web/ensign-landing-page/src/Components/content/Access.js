import React from 'react';
import { Formik, Field, Form, ErrorMessage } from 'formik';
import { useNavigate } from 'react-router-dom'
import * as Yup from 'yup';

export default function Access () {
    const navigate = useNavigate() 
    return (
      <Formik
        initialValues={{
          firstName: '',
          lastName: '',
          email: '',
          title: '',
          organization: '',
        }}
        validationSchema={Yup.object().shape({
          firstName: Yup.string().required('First Name is required'),
          lastName: Yup.string().required('Last Name is required'),
          email: Yup.string()
            .email('Email is invalid')
            .required('Email is required'),
          title: Yup.string(),
          organization: Yup.string(),
        })}
        onSubmit={fields => {
          navigate('/ensign-access')
          console.log(fields)
        }}
        render={({ errors, status, touched }) => (
          <Form class="w-96 p-7">
            <div>
              <h3 class="pb-2 text-2xl text-center font-bold">Request Alpha Access Today</h3>
              <p class="text-center pb-3">We're opening up Ensign on a limited basis. No credit card required.</p>
            </div>
            <div className="form-group" class="pb-3 text-left">
              <label htmlFor="firstName">First Name </label>
              <Field
                name="firstName"
                type="text"
                placeholder="*First Name"
                class="w-full p-1"
                className={
                  'form-control' +
                  (errors.firstName && touched.firstName ? ' is-invalid' : '')
                }
              />
              <ErrorMessage
                name="firstName"
                component="div"
                className="invalid-feedback"
              />
            </div>
            <div className="form-group" class="pb-3">
              <label htmlFor="lastName">Last Name </label>
              <Field
                name="lastName"
                type="text"
                placeholder="*Last Name"
                class="w-full p-1"
                className={
                  'form-control' +
                  (errors.lastName && touched.lastName ? ' is-invalid' : '')
                }
              />
              <ErrorMessage
                name="lastName"
                component="div"
                className="invalid-feedback"
              />
            </div>
            <div className="form-group" class="pb-3">
              <label htmlFor="email">Email address </label>
              <Field
                name="email"
                type="text"
                placeholder="*Email address"
                class="w-full p-1"
                className={
                  'form-control' +
                  (errors.email && touched.email ? ' is-invalid' : '')
                }
              />
              <ErrorMessage
                name="email"
                component="div"
                className="invalid-feedback"
              />
            </div>
            <div className="form-group" class="pb-3">
              <label htmlFor="title">Title </label>
              <Field 
                name="title" 
                type="title" 
                placeholder="Title" 
                class="w-full p-1"
                className={'form-control'} />
            </div>
            <div className="form-group" class="pb-3">
              <label htmlFor="organization">Organization </label>
              <Field
                name="organization"
                type="organization"
                placeholder="Organization"
                class="w-full p-1"
                className={'form-control'}
              />
            </div>
            <div class="pb-5">
            <label>
                <Field 
                type="checkbox" 
                name="notifications" 
                value="notifications"
                class="m-1" />
                I agree to receive notifications about Ensign from Rotational Labs. Your contact information will not be shared with external parties. Unsubscribe any time. 
              </label>
            </div>
            <div className="form-group w-52 mx-auto p-2 text-2xl text-center text-white bg-[#37A36E]">
              <button type="submit">
                Request Access
              </button>
            </div>
          </Form>
        )}
      />
    );
  }
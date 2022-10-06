import React from 'react';
import { Formik, Field, Form, ErrorMessage } from 'formik';
import { useNavigate } from 'react-router-dom'
import * as Yup from 'yup';
import './Access.css';

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
          notifications: false,
        }}
        validationSchema={Yup.object().shape({
          firstName: Yup.string().required('First Name is required'),
          lastName: Yup.string().required('Last Name is required'),
          email: Yup.string()
            .email('Email is invalid')
            .required('Email is required'),
          title: Yup.string(),
          organization: Yup.string(),
          notifications: Yup.bool().oneOf([true], 'Must allow notifications to continue'),
        })}
        onSubmit={fields => {
          fetch("https://api.rotational.app/v1/notifications/signup", {
            method: "POST",
            headers: {
              'Accept': 'application/json',
              'Content-Type': 'application/json'
            },
            body: JSON.stringify(fields),
            cache: 'default'
          }).then((rep)=> {
              console.log(rep);
              navigate('/ensign-access');
            });
        }}
        render={({ errors, status, touched }) => (
          <Form className="w-96 p-7" method="post">
            <div>
              <h3 className="pb-2 text-2xl text-center font-bold">Request Alpha Access Today</h3>
              <p className="text-center pb-3">We're opening up Ensign on a limited basis. No credit card required.</p>
            </div>
            <div className="form-group pb-3">
              <label htmlFor="firstName" className="hidden">First Name </label>
              <Field
                name="firstName"
                type="text"
                placeholder="First Name"
                className={
                  'w-full form-input' +
                  (errors.firstName && touched.firstName ? ' is-invalid' : '')
                }
              />
              <ErrorMessage
                name="firstName"
                component="div"
                className="text-red-500"
              />
            </div>
            <div className="form-group pb-3">
              <label htmlFor="lastName" className="hidden">Last Name </label>
              <Field
                name="lastName"
                type="text"
                placeholder="Last Name"
                className={
                  'w-full form-input' +
                  (errors.lastName && touched.lastName ? ' is-invalid' : '')
                }
              />
              <ErrorMessage
                name="lastName"
                component="div"
                className="text-red-500"
              />
            </div>
            <div className="form-group pb-3">
              <label htmlFor="email" className="hidden">Email address </label>
              <Field
                name="email"
                type="email"
                placeholder="Email address"
                className={
                  'w-full form-input' +
                  (errors.email && touched.email ? ' is-invalid' : '')
                }
              />
              <ErrorMessage
                name="email"
                component="div"
                className="text-red-500"
              />
            </div>
            <div className="form-group pb-3">
              <label htmlFor="title" className="hidden">Title </label>
              <Field
                name="title"
                type="title"
                placeholder="Title"
                className='w-full form-input' />
            </div>
            <div className="form-group pb-3">
              <label htmlFor="organization" className="hidden">Organization </label>
              <Field
                name="organization"
                type="organization"
                placeholder="Organization"
                className='w-full form-input'
              />
            </div>
            <div className="pb-5">
            <label>
                <Field
                type="checkbox"
                name="notifications"
                className={
                  'w-full form-checkbox' +
                  (errors.notifications && touched.notifications ? ' is-invalid' : '')
                } />
                <span className="ml-2">I agree to receive notifications about Ensign from Rotational Labs. Your contact information will not be shared with external parties. Unsubscribe any time.</span>
              </label>
              <ErrorMessage
                name="notifications"
                component="div"
                className="text-red-500"
              />
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
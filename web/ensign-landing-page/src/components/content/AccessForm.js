import React from 'react';
import { Formik, Field, Form, ErrorMessage } from 'formik';
import { useNavigate } from 'react-router-dom'
import * as Yup from 'yup';

export default function AccessForm () {
    const navigate = useNavigate()

    return (
      <Formik
        initialValues={{
          firstName: '',
          lastName: '',
          email: '',
          title: '',
          organization: '',
          cloudServiceProvider: '',
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
          cloudServiceProvider: Yup.string(),
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
        children={({ errors, status, touched }) => (
          <Form className="w-[26rem] p-7 bg-[#DED6C5] mx-auto">
            <div>
              <h3 className="pb-2 text-2xl font-bold">Request Alpha Access Today</h3>
              <p className="pb-3">We're opening up Ensign on a limited basis. No credit card required.</p>
            </div>
            <div className="form-group pb-3">
              <label htmlFor="firstName" className="hidden">First Name </label>
              <Field
                name="firstName"
                type="text"
                placeholder="First Name *"
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
                placeholder="Last Name *"
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
                placeholder="Email address *"
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
            <div className="form-group pb-3 form-multiselect">
            <label className="hidden">Cloud Service Provider</label>
              <Field 
                component="select"
                name="cloudServiceProvider"
                multiple={false}
                >
                  <option value="Cloud service provider not selected">Cloud service provider</option>
                  <option value="Amazon Web Services (AWS)">Amazon Web Services (AWS)</option>
                  <option value="Microsoft Azure">Microsoft Azure</option>
                  <option value="Google Cloud">Google Cloud</option>
                  <option value="Alibaba Cloud">Alibaba Cloud</option>
                  <option value="Digital Ocean">Digital Ocean</option>
                  <option value="IBM Cloud">IBM Cloud</option>
                  <option value="Oracle">Oracle</option>
                  <option value="Salesforce">Salesforce</option>
                  <option value="SAP">SAP</option>
                  <option value="Rackspace Cloud">Rackspace Cloud</option>
                  <option value="VMWare">VMWare</option>
                  <option value="On premises">On premises</option>
                  <option value="Other">Other</option>
                </Field>
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
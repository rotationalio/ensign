import { yupResolver } from '@hookform/resolvers/yup';
import { object, string } from 'yup';

const FORM_VALIDATION_SCHEMA = object({
    email: string().required('Email is required').email('Email is invalid'),
    password: string().required('Password is required'),

});

export const LOGIN_FORM_OPTIONS = {
    resolver: yupResolver(FORM_VALIDATION_SCHEMA),
    defaultValues: {
        email: '',
        password: '',
    },
};

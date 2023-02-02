/* eslint-disable prettier/prettier */
export const checkPassword = (password: string) => {
    const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#$%^&*])(?=.{12,})/;
    return passwordRegex.test(password);
};

export const checkPasswordContains12Characters = (password: string) => {
    const passwordRegex = /^(?=.{12,})/;
    return passwordRegex.test(password);
};

export const checkPasswordContainsOneLowerCase = (password: string) => {
    const passwordRegex = /^(?=.*[a-z])/;
    return passwordRegex.test(password);
};

export const checkPasswordContainsOneUpperCase = (password: string) => {
    const passwordRegex = /^(?=.*[A-Z])/;
    return passwordRegex.test(password);
};

export const checkPasswordContainsOneNumber = (password: string) => {
    const passwordRegex = /^(?=.*[0-9])/;
    return passwordRegex.test(password);
};

export const checkPasswordContainsOneSpecialChar = (password: string) => {
    const passwordRegex = /^(?=.*[#$%&'()*+,-./:;<=>?@[\\]^_`{|}~])/;
    return passwordRegex.test(password);
};

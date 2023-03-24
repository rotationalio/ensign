import React, { useEffect, useState } from 'react';

import {
  checkPasswordContains12Characters,
  checkPasswordContainsOneLowerCase,
  checkPasswordContainsOneNumber,
  checkPasswordContainsOneSpecialChar,
  checkPasswordContainsOneUpperCase,
} from '@/utils/passwordChecker';

interface PasswordStrengthProps {
  string: string;
  onMatch?: (result: boolean) => void;
}

const PasswordStrength = ({ string, onMatch }: PasswordStrengthProps) => {
  const [isContains12Characters, setIsContains12Characters] = useState<boolean>(false);
  const [isContainsOneLowerCase, setIsContainsOneLowerCase] = useState<boolean>(false);
  const [isContainsOneUpperCase, setIsContainsOneUpperCase] = useState<boolean>(false);
  const [isContainsOneNumber, setIsContainsOneNumber] = useState<boolean>(false);
  const [isContainsOneSpecialChar, setIsContainsOneSpecialChar] = useState<boolean>(false);

  const checkPasswordValidity = (password: string) => {
    setIsContains12Characters(checkPasswordContains12Characters(password));
    setIsContainsOneLowerCase(checkPasswordContainsOneLowerCase(password));
    setIsContainsOneUpperCase(checkPasswordContainsOneUpperCase(password));
    setIsContainsOneNumber(checkPasswordContainsOneNumber(password));
    setIsContainsOneSpecialChar(checkPasswordContainsOneSpecialChar(password));
  };

  useEffect(() => {
    checkPasswordValidity(string);
  }, [string]);

  // send result to parent component
  useEffect(() => {
    if (onMatch) {
      onMatch(
        isContains12Characters &&
          isContainsOneLowerCase &&
          isContainsOneUpperCase &&
          isContainsOneNumber &&
          isContainsOneSpecialChar
      );
    }
  }, [
    isContains12Characters,
    isContainsOneLowerCase,
    isContainsOneUpperCase,
    isContainsOneNumber,
    isContainsOneSpecialChar,
    onMatch,
  ]);

  //
  const bgStyle = (isMatch: boolean) => {
    if (isMatch) {
      return 'green-success';
    }
    return 'gray-500';
  };

  return (
    <div className="flex flex-col space-y-2">
      <div className="flex items-center space-x-2">
        <div className={`h-4 w-4 rounded-full bg-${bgStyle(isContains12Characters)}`} />
        <div className={`text-xs text-${bgStyle(isContains12Characters)}`}>12 characters</div>
      </div>
      <div className="flex items-center space-x-2">
        <div className={`h-4 w-4 rounded-full bg-${bgStyle(isContainsOneLowerCase)}`} />
        <div className={`text-xs text-${bgStyle(isContainsOneLowerCase)}`}>1 lowercase</div>
      </div>
      <div className="flex items-center space-x-2">
        <div className={`h-4 w-4 rounded-full bg-${bgStyle(isContainsOneUpperCase)}`} />
        <div className={`text-xs text-${bgStyle(isContainsOneUpperCase)}`}>1 uppercase</div>
      </div>
      <div className="flex items-center space-x-2">
        <div className={`h-4 w-4 rounded-full bg-${bgStyle(isContainsOneNumber)}`} />
        <div className={`text-xs text-${bgStyle(isContainsOneNumber)}`}>1 number</div>
      </div>
      <div className="flex items-center space-x-2">
        <div className={`h-4 w-4 rounded-full bg-${bgStyle(isContainsOneSpecialChar)}`} />
        <div className={`text-xs text-${bgStyle(isContainsOneSpecialChar)}`}>
          1 special character
        </div>
      </div>
    </div>
  );
};

export default PasswordStrength;

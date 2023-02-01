import {
  checkPasswordContains12Characters,
  checkPasswordContainsOneLowerCase,
  checkPasswordContainsOneNumber,
  checkPasswordContainsOneSpecialChar,
  checkPasswordContainsOneUpperCase,
} from '../password';

describe('checkPasswordContains12Characters', () => {
  it('should return false if password is less than 12 characters', () => {
    const result = checkPasswordContains12Characters('1234567890');
    expect(result).toEqual(false);
  });

  it('should return true if password is 12 characters', () => {
    const result = checkPasswordContains12Characters('123456789012');
    expect(result).toEqual(true);
  });

  it('should return true if password is more than 12 characters', () => {
    const result = checkPasswordContains12Characters('1234567890123');
    expect(result).toEqual(true);
  });
});

describe('checkPasswordContainsOneLowerCase', () => {
  it('should return false if password does not contain lowercase', () => {
    const result = checkPasswordContainsOneLowerCase('123456789012');
    expect(result).toEqual(false);
  });

  it('should return true if password contains lowercase', () => {
    const result = checkPasswordContainsOneLowerCase('123456789012a');
    expect(result).toEqual(true);
  });
});

describe('checkPasswordContainsOneUpperCase', () => {
  it('should return false if password does not contain uppercase', () => {
    const result = checkPasswordContainsOneUpperCase('123456789012');
    expect(result).toEqual(false);
  });

  it('should return true if password contains uppercase', () => {
    const result = checkPasswordContainsOneUpperCase('123456789012A');
    expect(result).toEqual(true);
  });
});

describe('checkPasswordContainsOneNumber', () => {
  it('should return false if password does not contain number', () => {
    const result = checkPasswordContainsOneNumber('abcdefghijk');
    expect(result).toEqual(false);
  });

  it('should return true if password contains number', () => {
    const result = checkPasswordContainsOneNumber('123456789012');
    expect(result).toEqual(true);
  });
});

describe('checkPasswordContainsSpecialCharacter', () => {
  it('should return false if password does not contain special character', () => {
    const result = checkPasswordContainsOneSpecialChar('123456789012');
    expect(result).toEqual(false);
  });
});

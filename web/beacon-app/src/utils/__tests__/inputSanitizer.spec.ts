import { inputSanitizer } from '../inputSanitzer';

describe('Input Sanitizer', () => {
  it('should remove leading and trailing whitespaces from the input', () => {
    const inputWithWhitespaces = '  Hello, World!   ';
    const sanitizedInput = inputSanitizer(inputWithWhitespaces);

    const expectedOutput = 'Hello, World!';

    // Assertion
    expect(sanitizedInput).toEqual(expectedOutput);
  });

  it('should prevent SQL injection', () => {
    // Input data with potential SQL injection
    const inputWithSqlInjection = "John'; DROP TABLE users;";
    const sanitizedInput = inputSanitizer(inputWithSqlInjection);

    const expectedOutput = "John\\'; DROP TABLE users;";

    // Assertion
    expect(sanitizedInput).toEqual(expectedOutput);
  });

  it('should prevent JS injection', () => {
    // Input data with potential JS injection
    const inputWithJsInjection = '<script>alert("XSS Attack");</script>';
    const sanitizedInput = inputSanitizer(inputWithJsInjection);
    const expectedOutput = '';

    // Assertion
    expect(sanitizedInput).toEqual(expectedOutput);
  });

  it('should prevent XSS attack', () => {
    // Input data with potential XSS attack
    const inputWithXssAttack = '<img src="invalid-url" onerror="alert(\'XSS Attack\');" />';
    const sanitizedInput = inputSanitizer(inputWithXssAttack);
    const expectedOutput = '&lt;img src="invalid-url"&gt;';

    // Assertion
    expect(sanitizedInput).toEqual(expectedOutput);
  });
});

import { inputSanitizer, sqlInputSanitizer } from '../utils';

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

describe('SQL Input Sanitizer', () => {
  it('should remove leading and trailing whitespaces from the input', () => {
    const inputWithWhitespaces = '  Hello, World!   ';
    const sanitizedInput = sqlInputSanitizer(inputWithWhitespaces);

    const expectedOutput = 'Hello, World!';

    // Assertion
    expect(sanitizedInput).toEqual(expectedOutput);
  });

  it('should prevent JS injection', () => {
    // Input data with potential JS injection
    const inputWithJsInjection = '<script>alert("XSS Attack");</script>';
    const sanitizedInput = sqlInputSanitizer(inputWithJsInjection);
    const expectedOutput = '';

    // Assertion
    expect(sanitizedInput).toEqual(expectedOutput);
  });

  it('should prevent XSS attack', () => {
    // Input data with potential XSS attack
    const inputWithXssAttack = '<img src="invalid-url" onerror="alert(\'XSS Attack\');" />';
    const sanitizedInput = sqlInputSanitizer(inputWithXssAttack);
    const expectedOutput = '<img src="invalid-url">';

    // Assertion
    expect(sanitizedInput).toEqual(expectedOutput);
  });

  it('should return sql query', () => {
    // Input data with potential SQL injection
    const inputWithSqlInjection = "John'; DROP TABLE users;";
    const sanitizedInput = sqlInputSanitizer(inputWithSqlInjection);

    const expectedOutput = "John'; DROP TABLE users;";

    // Assertion
    expect(sanitizedInput).toEqual(expectedOutput);
  });

  // add sql query with operator to check if it is working

  it('should return sql query', () => {
    const sqlWithOperator = 'SELECT * FROM users WHERE name = "John" AND age > 30';
    const sanitizedInput = sqlInputSanitizer(sqlWithOperator);

    const expectedOutput = 'SELECT * FROM users WHERE name = "John" AND age > 30';

    // Assertion
    expect(sanitizedInput).toEqual(expectedOutput);
  });
});

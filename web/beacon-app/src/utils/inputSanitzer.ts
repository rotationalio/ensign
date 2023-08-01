import DOMPurify from 'dompurify';

export const inputSanitizer = (input: string) => {
  //  prevent XSS attacks
  const sanitizedInput = DOMPurify.sanitize(input);
  // prevent SQL injection
  const sanitizedSqlInjection = sanitizedInput.replace(/'/g, "\\'");
  // prevent JS injection
  const jsInjectionSafeInput = sanitizedSqlInjection.replace(/</g, '&lt;').replace(/>/g, '&gt;');
  // prevent leading and trailing spaces
  const finalSanitizedInput = jsInjectionSafeInput.trim();

  return finalSanitizedInput;
};

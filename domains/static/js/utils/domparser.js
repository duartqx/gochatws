/*
 * Parses a string into an HTML document and returns the body of the document 
 * as an HTML element.
 *
 * @param {string} body - The string to parse into an HTML document.
 * @returns {HTMLBodyElement} The body of the parsed HTML document.
 *
 * @example
 * // returns <body>...</body>
 * domParser('<html><head></head><body>...</body></html>');
 */
function domParser(body) {
  const parser = new DOMParser();
  const doc = parser.parseFromString(body, "text/html");
  return doc.body;
}

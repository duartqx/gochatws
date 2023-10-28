/**
 * Fetches the HTML content from a specified URL.
 *
 * @async
 * @function fetchHtml
 * @param {string} url - The URL to fetch the HTML content from.
 * @returns {Promise<string>} A promise that resolves to the fetched HTML content.
 * @throws Will throw an error if the network response is not OK or if an
 * error occurs while fetching the HTML content.
 */
async function fetchHtml(url) {
  try {
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error("Network response was not OK");
    }
    const htmlData = await response.text();
    return htmlData;
  } catch (error) {
    console.error("Error:", error);
  }
}

/**
 * Fetches data from the provided URL and dispatches an alert event.
 *
 * @async
 * @function fetchTransitionAndDispatches
 * @param {string} url - The URL to fetch data from.
 * @param {Object} detail - The detail object passed to the alert event.
 * @returns {Promise<void>} A promise that resolves when the operation is complete.
 * @see {@link https://htmx.org/docs/#ajax} htmx.ajax
 * @see {@link https://developer.mozilla.org/en-US/docs/Web/API/Window/pushState} window.history.pushState
 */
async function fetchTransitionAndDispatches(url, detail) {
  await htmx.ajax("GET", url, { swap: "transition:true" });
  window.history.pushState({}, "", url);
  dispatchAlert(detail);
}

/**
 * Dispatches a custom alert event after a delay.
 *
 * @function dispatchAlert
 * @param {Object} detail - The detail object passed to the alert event.
 * @see {@link https://developer.mozilla.org/en-US/docs/Web/API/Window/setTimeout} setTimeout
 * @see {@link https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/dispatchEvent} dispatchEvent
 */
function dispatchAlert(detail) {
  setTimeout(() => {
    window.dispatchEvent(new CustomEvent("alert-message", detail));
  }, 500);
}

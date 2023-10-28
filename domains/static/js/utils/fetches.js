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
 * Fetches the HTML content from a specified URL, starts a view transition,
 * processes the HTML content with HTMX, and dispatches an alert.
 *
 * @async
 * @function fetchTransitionAndDispatches
 * @param {string} url - The URL to fetch the HTML content from.
 * @param {Object} detail - The detail of the alert to dispatch.
 * @returns {Promise<void>} A promise that resolves when the view transition
 * is finished and the alert is dispatched.
 */
async function fetchTransitionAndDispatches(url, detail) {
  const htmlData = await fetchHtml(url);

  const transition = document.startViewTransition(() => {
    document.documentElement.innerHTML = htmlData;
  });
  await transition.finished;
  htmx.process(document.documentElement);
  window.history.pushState({}, "", "/");
  dispatchAlert(detail);
}

/**
 * Dispatches an alert on '#generic-alert'
 *
 * @function dispatchAlert
 * @param {Object} detail - The detail of the alert to dispatch.
 */
function dispatchAlert(detail) {
  htmx
    .find("#generic-alert")
    .dispatchEvent(new CustomEvent("alert-message", detail));
}

/**
 * Handles user registration by sending a POST request to the server with user
 * credentials. If the registration is successful, it redirects the user to the
 * login page. If the registration fails, it displays an error message.
 *
 * @param {Event} e - The event object from the form submission.
 * @returns {Promise<void>} A promise that resolves when the registration 
 * process is complete.
 *
 * @example
 * register(event);
 */
function register(e) {
  let values = htmx.values(e.target);
  if (!values.username || !values.password) {
    return false;
  }

  if (values.password !== values.password2) {
    return false
  }

  fetch("/api/register", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(values),
  })
    .then((res) => {
      if (!res.ok) return res.json();
      htmx.ajax("GET", "/login", { target: "body", swap: "transition:true" });
      window.history.pushState({}, "", "/login");
    })
    .then((res) => {
      if (res && res.error) registerAlert(res.error);
    });
}

/**
 * Displays an alert message with the provided message.
 * @param {string} msg - The message to display in the alert.
 */
function registerAlert(msg) {
  let alertEl = htmx.find("#registerAlert");
  alertEl.innerText = msg;
  htmx.removeClass(alertEl, "collapse");
}

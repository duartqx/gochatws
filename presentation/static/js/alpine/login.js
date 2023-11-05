/*
 * Logs in the user by sending a POST request to the "/api/login"
 * endpoint with the provided username and password.
 * If the login is successful, it redirects the user to the home page ("/").
 * If there is an error, it displays an error message using the `loginAlert` function.
 *
 * @param {Event} e - The event object.
 * @returns {boolean} Returns false if the username or password is missing.
 *
 */
function login(e) {
  let values = htmx.values(e.target);
  if (!values.username || !values.password) {
    return false;
  }

  fetch("/api/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(values),
  })
    .then((res) => {
      if (!res.ok) return res.json();
      htmx.ajax("GET", "/", { target: "body", swap: "transition:true" });
      window.history.pushState({}, "", "/");
    })
    .then((res) => {
      if (res && res.error) loginAlert(res.error);
    });
}

/**
 * Displays an alert message with the provided message.
 * @param {string} msg - The message to display in the alert.
 */
function loginAlert(msg) {
  let alertEl = htmx.find("#loginAlert");
  alertEl.innerText = msg;
  htmx.removeClass(alertEl, "collapse");
}

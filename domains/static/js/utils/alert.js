function alertWithMessage(msg, alertClass) {
  let alertMessageEl = htmx.find("#generic-alert");
  alertMessageEl.innerText = msg;
  htmx.toggleClass(alertMessageEl, alertClass);

  setTimeout(() => {
    htmx.toggleClass(alertMessageEl, "hidden");
    setTimeout(() => {
      htmx.toggleClass(alertMessageEl, "hidden");
      setTimeout(() => {
        htmx.toggleClass(alertMessageEl, alertClass);
        alertMessageEl.innerText = "_";
      }, 500)
    }, 5000);
  }, 500);
}

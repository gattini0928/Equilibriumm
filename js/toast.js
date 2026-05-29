document.body.addEventListener("htmx:afterSwap", (event) => {
  if (event.detail.target.id !== "toast") return;

  const toast = event.detail.target.querySelector(".toast");
  if (!toast) return;

  clearTimeout(window._toastTimer);

  window._toastTimer = setTimeout(() => {
    toast.classList.add("hide");

    setTimeout(() => {
      event.detail.target.innerHTML = "";
    }, 500);
  }, 3000);
});
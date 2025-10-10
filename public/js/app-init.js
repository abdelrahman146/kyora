// Alpine global store for app UI state & flash messages
window.app = () => ({
  theme:
    localStorage.getItem("kyora.theme") ||
    (window.matchMedia("(prefers-color-scheme: dark)").matches
      ? "dark"
      : "light"),
  sidebarOpen: false,
  active: "",
  toggleTheme() {
    this.theme = this.theme === "dark" ? "light" : "dark";
    localStorage.setItem("kyora.theme", this.theme);
  },
});

// Flash message helper
document.addEventListener("alpine:init", () => {
  Alpine.store("flash", {
    messages: [],
    push(text, type = "info", timeout = 2800) {
      const idx = this.messages.push({ text, type }) - 1;
      if (timeout) setTimeout(() => this.dismiss(idx), timeout);
    },
    dismiss(idx) {
      this.messages.splice(idx, 1);
    },
  });
});

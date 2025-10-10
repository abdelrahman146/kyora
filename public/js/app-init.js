function appState() {
  const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
  const savedTheme = localStorage.getItem("theme");
  const savedDir = localStorage.getItem("dir");
  return {
    theme: savedTheme || (prefersDark ? "dark" : "light"),
    dir: savedDir || document.documentElement.getAttribute("dir") || "ltr",
    toggleTheme() {
      this.theme = this.theme === "dark" ? "light" : "dark";
      localStorage.setItem("theme", this.theme);
    },
    toggleDir() {
      this.dir = this.dir === "rtl" ? "ltr" : "rtl";
      localStorage.setItem("dir", this.dir);
    },
  };
}

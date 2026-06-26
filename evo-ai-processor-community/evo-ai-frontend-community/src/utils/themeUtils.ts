// Função auxiliar para inicializar o tema antes do React montar
export function initTheme() {
  // Verificar localStorage
  const savedTheme = localStorage.getItem('theme');
  if (savedTheme === 'dark') {
    document.documentElement.classList.add('dark');
  } else if (savedTheme === 'light') {
    document.documentElement.classList.remove('dark');
  } else {
    // Padrão para dark mode (escuro) quando o cache é limpo
    document.documentElement.classList.add('dark');
    try {
      localStorage.setItem('theme', 'dark');
    } catch (e) {
      // Ignore if localStorage is not available
    }
  }
}

import { useState, useMemo, createContext, useContext } from 'react';
import { createTheme, Theme } from '@mui/material/styles';

interface ThemeContextType {
  theme: Theme;
  toggleThemeMode: () => void;
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

export const ThemeProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [mode, setMode] = useState<'light' | 'dark'>('light');

  const theme = useMemo(
    () =>
      createTheme({
        direction: 'rtl',
        typography: {
          fontFamily: 'Vazirmatn, Roboto, "Helvetica Neue", Arial, sans-serif',
        },
        palette: {
          mode,
          primary: {
            main: mode === 'light' ? '#3A86FF' : '#6A9BFF',
          },
          background: {
            default: mode === 'light' ? '#F8F8F8' : '#1A1A2E',
            paper: mode === 'light' ? '#FFFFFF' : '#242A42',
          },
          text: {
            primary: mode === 'light' ? '#212121' : '#E0E0E0',
            secondary: mode === 'light' ? '#757575' : '#A0A0A0',
          },
          success: {
            main: '#4CAF50',
          },
          warning: {
            main: '#FFC107',
          },
          error: {
            main: '#F44336',
          },
        },
        components: {
          MuiButton: {
            styleOverrides: {
              root: {
                borderRadius: 4,
                padding: '8px 16px',
              },
            },
          },
          MuiCssBaseline: {
            styleOverrides: {
              body: {
                direction: 'rtl',
              },
            },
          },
        },
      }),
    [mode]
  );

  const toggleThemeMode = () => {
    setMode((prevMode) => (prevMode === 'light' ? 'dark' : 'light'));
  };

  return (
    <ThemeContext.Provider value={{ theme, toggleThemeMode }}>
      {children}
    </ThemeContext.Provider>
  );
};

export const useTheme = () => {
  const context = useContext(ThemeContext);
  if (context === undefined) {
    throw new Error('useTheme must be used within a ThemeProvider');
  }
  return context;
};

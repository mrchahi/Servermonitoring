import React from 'react';
import {
  AppBar,
  Toolbar,
  Typography,
  IconButton,
  Box,
  useTheme,
} from '@mui/material';
import {
  Brightness4 as DarkModeIcon,
  Brightness7 as LightModeIcon,
} from '@mui/icons-material';
import { useThemeMode } from '../hooks/useTheme';

export const Header: React.FC = () => {
  const theme = useTheme();
  const { toggleThemeMode } = useThemeMode();

  return (
    <AppBar
      position="fixed"
      sx={{
        backgroundColor: 'background.paper',
        color: 'text.primary',
        boxShadow: 1,
      }}
    >
      <Toolbar>
        <Typography
          variant="h6"
          component="h1"
          sx={{
            flexGrow: 1,
            fontFamily: 'Vazirmatn',
            fontWeight: 700,
          }}
        >
          پنل مدیریت سرور
        </Typography>
        <Box>
          <IconButton onClick={toggleThemeMode} color="inherit">
            {theme.palette.mode === 'dark' ? (
              <LightModeIcon />
            ) : (
              <DarkModeIcon />
            )}
          </IconButton>
        </Box>
      </Toolbar>
    </AppBar>
  );
};

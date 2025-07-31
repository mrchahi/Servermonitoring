import React from 'react';
import {
  Drawer,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  IconButton,
  useTheme,
  useMediaQuery,
} from '@mui/material';
import {
  Dashboard as DashboardIcon,
  Storage as ServicesIcon,
  Router as PortsIcon,
  Security as FirewallIcon,
  Description as LogsIcon,
  Settings as SettingsIcon,
  Menu as MenuIcon,
} from '@mui/icons-material';
import { useNavigate, useLocation } from 'react-router-dom';

const menuItems = [
  { text: 'داشبورد', icon: <DashboardIcon />, path: '/' },
  { text: 'سرویس‌ها', icon: <ServicesIcon />, path: '/services' },
  { text: 'پورت‌ها', icon: <PortsIcon />, path: '/ports' },
  { text: 'فایروال', icon: <FirewallIcon />, path: '/firewall' },
  { text: 'لاگ‌ها', icon: <LogsIcon />, path: '/logs' },
  { text: 'تنظیمات', icon: <SettingsIcon />, path: '/settings' },
];

const drawerWidth = 240;

export const Sidebar: React.FC = () => {
  const [mobileOpen, setMobileOpen] = React.useState(false);
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  const navigate = useNavigate();
  const location = useLocation();

  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen);
  };

  const drawer = (
    <List>
      {menuItems.map((item) => (
        <ListItem
          button
          key={item.text}
          onClick={() => navigate(item.path)}
          selected={location.pathname === item.path}
          sx={{
            '&.Mui-selected': {
              backgroundColor: 'primary.main',
              color: 'white',
              '&:hover': {
                backgroundColor: 'primary.dark',
              },
              '& .MuiListItemIcon-root': {
                color: 'white',
              },
            },
          }}
        >
          <ListItemIcon sx={{ minWidth: 40 }}>{item.icon}</ListItemIcon>
          <ListItemText 
            primary={item.text} 
            sx={{ 
              '& .MuiTypography-root': { 
                fontFamily: 'Vazirmatn' 
              } 
            }} 
          />
        </ListItem>
      ))}
    </List>
  );

  return (
    <>
      {isMobile && (
        <IconButton
          color="inherit"
          aria-label="منو"
          edge="start"
          onClick={handleDrawerToggle}
          sx={{ position: 'fixed', right: 16, top: 16, zIndex: 1200 }}
        >
          <MenuIcon />
        </IconButton>
      )}
      <Drawer
        variant={isMobile ? 'temporary' : 'permanent'}
        anchor="right"
        open={isMobile ? mobileOpen : true}
        onClose={handleDrawerToggle}
        ModalProps={{ keepMounted: true }}
        sx={{
          '& .MuiDrawer-paper': {
            width: drawerWidth,
            boxSizing: 'border-box',
            borderLeft: 0,
            borderRight: '1px solid rgba(0, 0, 0, 0.12)',
          },
        }}
      >
        {drawer}
      </Drawer>
    </>
  );
};

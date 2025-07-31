import React, { useEffect, useState } from 'react';
import {
  Grid,
  Paper,
  Typography,
  Box,
  CircularProgress,
} from '@mui/material';
import {
  Memory as CpuIcon,
  Storage as RamIcon,
  Save as DiskIcon,
  NetworkCheck as NetworkIcon,
} from '@mui/icons-material';

interface SystemStats {
  cpu: {
    usagePercent: number;
    temperature?: number;
  };
  memory: {
    total: number;
    used: number;
    free: number;
    usagePercent: number;
  };
  disk: {
    total: number;
    used: number;
    free: number;
    usagePercent: number;
  };
  network: {
    bytesSent: number;
    bytesReceived: number;
  };
  system: {
    hostname: string;
    uptime: number;
    loadAverage: number[];
  };
}

const formatBytes = (bytes: number): string => {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`;
};

const ResourceCard: React.FC<{
  title: string;
  icon: React.ReactNode;
  value: number;
  total?: number;
  unit?: string;
  color: string;
}> = ({ title, icon, value, total, unit, color }) => (
  <Paper
    elevation={2}
    sx={{
      p: 2,
      display: 'flex',
      flexDirection: 'column',
      height: 180,
      position: 'relative',
    }}
  >
    <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
      <Box sx={{ color: 'text.secondary', mr: 1 }}>{icon}</Box>
      <Typography variant="h6" component="h2" sx={{ fontFamily: 'Vazirmatn' }}>
        {title}
      </Typography>
    </Box>
    <Box
      sx={{
        position: 'relative',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        flexGrow: 1,
      }}
    >
      <CircularProgress
        variant="determinate"
        value={value}
        size={80}
        sx={{
          color,
          '& .MuiCircularProgress-circle': {
            strokeLinecap: 'round',
          },
        }}
      />
      <Box
        sx={{
          position: 'absolute',
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
        }}
      >
        <Typography
          variant="h6"
          component="div"
          sx={{ fontFamily: 'Roboto', lineHeight: 1 }}
        >
          {value.toFixed(1)}%
        </Typography>
        {total && (
          <Typography
            variant="caption"
            sx={{ color: 'text.secondary', fontFamily: 'Roboto' }}
          >
            {formatBytes(total)}
          </Typography>
        )}
      </Box>
    </Box>
  </Paper>
);

export const Dashboard: React.FC = () => {
  const [stats, setStats] = useState<SystemStats | null>(null);
  const [ws, setWs] = useState<WebSocket | null>(null);

  useEffect(() => {
    const socket = new WebSocket('ws://localhost:8443/ws/stats');

    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      setStats(data);
    };

    socket.onclose = () => {
      // Try to reconnect in 5 seconds
      setTimeout(() => {
        setWs(new WebSocket('ws://localhost:8443/ws/stats'));
      }, 5000);
    };

    setWs(socket);

    return () => {
      socket.close();
    };
  }, []);

  if (!stats) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  const getStatusColor = (value: number): string => {
    if (value >= 90) return '#F44336';
    if (value >= 70) return '#FFC107';
    return '#4CAF50';
  };

  return (
    <Grid container spacing={3}>
      <Grid item xs={12} sm={6} md={3}>
        <ResourceCard
          title="CPU"
          icon={<CpuIcon />}
          value={stats.cpu.usagePercent}
          color={getStatusColor(stats.cpu.usagePercent)}
        />
      </Grid>
      <Grid item xs={12} sm={6} md={3}>
        <ResourceCard
          title="حافظه"
          icon={<RamIcon />}
          value={stats.memory.usagePercent}
          total={stats.memory.total}
          color={getStatusColor(stats.memory.usagePercent)}
        />
      </Grid>
      <Grid item xs={12} sm={6} md={3}>
        <ResourceCard
          title="دیسک"
          icon={<DiskIcon />}
          value={stats.disk.usagePercent}
          total={stats.disk.total}
          color={getStatusColor(stats.disk.usagePercent)}
        />
      </Grid>
      <Grid item xs={12} sm={6} md={3}>
        <ResourceCard
          title="شبکه"
          icon={<NetworkIcon />}
          value={(stats.network.bytesSent + stats.network.bytesReceived) / 1024 / 1024}
          unit="MB/s"
          color="#3A86FF"
        />
      </Grid>
    </Grid>
  );
};

import React, { useEffect, useState } from 'react';
import {
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
  Box,
  Chip,
  CircularProgress,
} from '@mui/material';

interface Port {
  number: number;
  protocol: string;
  status: string;
  service: string;
  description: string;
  allowedIPs?: string[];
}

export const Ports: React.FC = () => {
  const [ports, setPorts] = useState<Port[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchPorts = async () => {
      try {
        const response = await fetch('http://localhost:8443/api/ports');
        if (!response.ok) throw new Error('خطا در دریافت لیست پورت‌ها');
        const data = await response.json();
        setPorts(data);
        setError(null);
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchPorts();
  }, []);

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', mt: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <>
      <Typography variant="h5" component="h1" gutterBottom sx={{ fontFamily: 'Vazirmatn' }}>
        مدیریت پورت‌ها
      </Typography>

      {error && (
        <Typography color="error" sx={{ mb: 2 }}>
          {error}
        </Typography>
      )}

      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell align="right">پورت</TableCell>
              <TableCell align="right">پروتکل</TableCell>
              <TableCell align="right">وضعیت</TableCell>
              <TableCell align="right">سرویس</TableCell>
              <TableCell align="right">IP‌های مجاز</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {ports.map((port) => (
              <TableRow key={`${port.number}-${port.protocol}`}>
                <TableCell align="right">{port.number}</TableCell>
                <TableCell align="right">
                  <Chip
                    label={port.protocol.toUpperCase()}
                    size="small"
                    color="primary"
                    variant="outlined"
                  />
                </TableCell>
                <TableCell align="right">
                  <Chip
                    label={port.status === 'open' ? 'باز' : 'بسته'}
                    color={port.status === 'open' ? 'success' : 'error'}
                    size="small"
                  />
                </TableCell>
                <TableCell align="right">
                  <Typography sx={{ fontFamily: 'Vazirmatn' }}>
                    {port.service}
                  </Typography>
                </TableCell>
                <TableCell align="right">
                  {port.allowedIPs?.length ? (
                    port.allowedIPs.map((ip) => (
                      <Chip
                        key={ip}
                        label={ip}
                        size="small"
                        sx={{ m: 0.5 }}
                        variant="outlined"
                      />
                    ))
                  ) : (
                    <Typography variant="caption" color="text.secondary">
                      همه IP‌ها
                    </Typography>
                  )}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </>
  );
};

import React, { useEffect, useState } from 'react';
import {
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  IconButton,
  Chip,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Typography,
  Box,
  CircularProgress,
} from '@mui/material';
import {
  PlayArrow as StartIcon,
  Stop as StopIcon,
  Refresh as RestartIcon,
  Check as EnableIcon,
  Close as DisableIcon,
} from '@mui/icons-material';

interface Service {
  name: string;
  displayName: string;
  status: string;
  port?: number;
  description: string;
  autoStart: boolean;
}

interface ServiceAction {
  action: 'start' | 'stop' | 'restart' | 'enable' | 'disable';
  name: string;
}

export const Services: React.FC = () => {
  const [services, setServices] = useState<Service[]>([]);
  const [loading, setLoading] = useState(true);
  const [confirmDialog, setConfirmDialog] = useState<{
    open: boolean;
    service?: Service;
    action?: string;
  }>({ open: false });
  const [error, setError] = useState<string | null>(null);

  const fetchServices = async () => {
    try {
      const response = await fetch('http://localhost:8443/api/services');
      if (!response.ok) throw new Error('خطا در دریافت لیست سرویس‌ها');
      const data = await response.json();
      setServices(data);
      setError(null);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchServices();
  }, []);

  const handleServiceAction = async (service: Service, action: ServiceAction['action']) => {
    try {
      const response = await fetch(`http://localhost:8443/api/services/${service.name}/action`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ action }),
      });

      if (!response.ok) throw new Error('خطا در اجرای عملیات');

      await fetchServices(); // Refresh services list
      setError(null);
    } catch (err) {
      setError(err.message);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active':
        return 'success';
      case 'inactive':
        return 'warning';
      case 'failed':
        return 'error';
      default:
        return 'default';
    }
  };

  const getStatusText = (status: string) => {
    switch (status) {
      case 'active':
        return 'فعال';
      case 'inactive':
        return 'غیرفعال';
      case 'failed':
        return 'خطا';
      default:
        return 'نامشخص';
    }
  };

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
        مدیریت سرویس‌ها
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
              <TableCell align="right">نام سرویس</TableCell>
              <TableCell align="right">وضعیت</TableCell>
              <TableCell align="right">پورت</TableCell>
              <TableCell align="right">توضیحات</TableCell>
              <TableCell align="right">راه‌اندازی خودکار</TableCell>
              <TableCell align="center">عملیات</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {services.map((service) => (
              <TableRow key={service.name}>
                <TableCell align="right">
                  <Typography sx={{ fontFamily: 'Vazirmatn' }}>
                    {service.displayName}
                  </Typography>
                  <Typography variant="caption" sx={{ color: 'text.secondary' }}>
                    {service.name}
                  </Typography>
                </TableCell>
                <TableCell align="right">
                  <Chip
                    label={getStatusText(service.status)}
                    color={getStatusColor(service.status)}
                    size="small"
                  />
                </TableCell>
                <TableCell align="right">{service.port || '-'}</TableCell>
                <TableCell align="right">{service.description}</TableCell>
                <TableCell align="right">
                  <Chip
                    label={service.autoStart ? 'فعال' : 'غیرفعال'}
                    color={service.autoStart ? 'success' : 'default'}
                    size="small"
                  />
                </TableCell>
                <TableCell align="center">
                  <IconButton
                    size="small"
                    color="success"
                    onClick={() => setConfirmDialog({
                      open: true,
                      service,
                      action: 'start'
                    })}
                    disabled={service.status === 'active'}
                  >
                    <StartIcon />
                  </IconButton>
                  <IconButton
                    size="small"
                    color="error"
                    onClick={() => setConfirmDialog({
                      open: true,
                      service,
                      action: 'stop'
                    })}
                    disabled={service.status === 'inactive'}
                  >
                    <StopIcon />
                  </IconButton>
                  <IconButton
                    size="small"
                    color="primary"
                    onClick={() => setConfirmDialog({
                      open: true,
                      service,
                      action: 'restart'
                    })}
                  >
                    <RestartIcon />
                  </IconButton>
                  {service.autoStart ? (
                    <IconButton
                      size="small"
                      onClick={() => setConfirmDialog({
                        open: true,
                        service,
                        action: 'disable'
                      })}
                    >
                      <DisableIcon />
                    </IconButton>
                  ) : (
                    <IconButton
                      size="small"
                      onClick={() => setConfirmDialog({
                        open: true,
                        service,
                        action: 'enable'
                      })}
                    >
                      <EnableIcon />
                    </IconButton>
                  )}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      <Dialog
        open={confirmDialog.open}
        onClose={() => setConfirmDialog({ open: false })}
      >
        <DialogTitle sx={{ fontFamily: 'Vazirmatn' }}>تأیید عملیات</DialogTitle>
        <DialogContent>
          <Typography sx={{ fontFamily: 'Vazirmatn' }}>
            {confirmDialog.action === 'start' && 'آیا می‌خواهید این سرویس را شروع کنید؟'}
            {confirmDialog.action === 'stop' && 'آیا می‌خواهید این سرویس را متوقف کنید؟'}
            {confirmDialog.action === 'restart' && 'آیا می‌خواهید این سرویس را مجدداً راه‌اندازی کنید؟'}
            {confirmDialog.action === 'enable' && 'آیا می‌خواهید راه‌اندازی خودکار را فعال کنید؟'}
            {confirmDialog.action === 'disable' && 'آیا می‌خواهید راه‌اندازی خودکار را غیرفعال کنید؟'}
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button
            onClick={() => setConfirmDialog({ open: false })}
            color="inherit"
            sx={{ fontFamily: 'Vazirmatn' }}
          >
            لغو
          </Button>
          <Button
            onClick={() => {
              if (confirmDialog.service && confirmDialog.action) {
                handleServiceAction(confirmDialog.service, confirmDialog.action);
                setConfirmDialog({ open: false });
              }
            }}
            color={confirmDialog.action === 'stop' ? 'error' : 'primary'}
            sx={{ fontFamily: 'Vazirmatn' }}
          >
            تأیید
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
};

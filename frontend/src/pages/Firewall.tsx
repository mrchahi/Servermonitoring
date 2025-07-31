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
  Button,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  CircularProgress,
  Chip,
} from '@mui/material';
import {
  Add as AddIcon,
  Delete as DeleteIcon,
} from '@mui/icons-material';

interface FirewallRule {
  id: number;
  action: string;
  protocol: string;
  port: number;
  source: string;
  description: string;
  enabled: boolean;
}

interface FirewallRuleRequest {
  action: string;
  protocol: string;
  port: number;
  source: string;
  description: string;
}

export const Firewall: React.FC = () => {
  const [rules, setRules] = useState<FirewallRule[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [openDialog, setOpenDialog] = useState(false);
  const [deleteConfirm, setDeleteConfirm] = useState<number | null>(null);
  const [newRule, setNewRule] = useState<FirewallRuleRequest>({
    action: 'allow',
    protocol: 'tcp',
    port: 80,
    source: '',
    description: '',
  });

  const fetchRules = async () => {
    try {
      const response = await fetch('http://localhost:8443/api/firewall/rules');
      if (!response.ok) throw new Error('خطا در دریافت قوانین فایروال');
      const data = await response.json();
      setRules(data);
      setError(null);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchRules();
  }, []);

  const handleAddRule = async () => {
    try {
      const response = await fetch('http://localhost:8443/api/firewall/rules', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(newRule),
      });

      if (!response.ok) throw new Error('خطا در افزودن قانون');
      await fetchRules();
      setOpenDialog(false);
      setNewRule({
        action: 'allow',
        protocol: 'tcp',
        port: 80,
        source: '',
        description: '',
      });
    } catch (err: any) {
      setError(err.message);
    }
  };

  const handleDeleteRule = async (id: number) => {
    try {
      const response = await fetch(`http://localhost:8443/api/firewall/rules/${id}`, {
        method: 'DELETE',
      });

      if (!response.ok) throw new Error('خطا در حذف قانون');
      await fetchRules();
      setDeleteConfirm(null);
    } catch (err: any) {
      setError(err.message);
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
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h5" component="h1" sx={{ fontFamily: 'Vazirmatn' }}>
          مدیریت فایروال
        </Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => setOpenDialog(true)}
          sx={{ fontFamily: 'Vazirmatn' }}
        >
          افزودن قانون جدید
        </Button>
      </Box>

      {error && (
        <Typography color="error" sx={{ mb: 2 }}>
          {error}
        </Typography>
      )}

      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell align="right">شماره</TableCell>
              <TableCell align="right">عملیات</TableCell>
              <TableCell align="right">پروتکل</TableCell>
              <TableCell align="right">پورت</TableCell>
              <TableCell align="right">منبع</TableCell>
              <TableCell align="right">توضیحات</TableCell>
              <TableCell align="center">حذف</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {rules.map((rule) => (
              <TableRow key={rule.id}>
                <TableCell align="right">{rule.id}</TableCell>
                <TableCell align="right">
                  <Chip
                    label={rule.action === 'allow' ? 'مجاز' : 'مسدود'}
                    color={rule.action === 'allow' ? 'success' : 'error'}
                    size="small"
                  />
                </TableCell>
                <TableCell align="right">
                  <Chip
                    label={rule.protocol.toUpperCase()}
                    variant="outlined"
                    size="small"
                  />
                </TableCell>
                <TableCell align="right">{rule.port}</TableCell>
                <TableCell align="right">{rule.source || 'همه'}</TableCell>
                <TableCell align="right">{rule.description}</TableCell>
                <TableCell align="center">
                  <IconButton
                    size="small"
                    color="error"
                    onClick={() => setDeleteConfirm(rule.id)}
                  >
                    <DeleteIcon />
                  </IconButton>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      {/* Add Rule Dialog */}
      <Dialog open={openDialog} onClose={() => setOpenDialog(false)} maxWidth="sm" fullWidth>
        <DialogTitle sx={{ fontFamily: 'Vazirmatn' }}>افزودن قانون جدید</DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 2 }}>
            <FormControl fullWidth>
              <InputLabel id="action-label">عملیات</InputLabel>
              <Select
                labelId="action-label"
                value={newRule.action}
                label="عملیات"
                onChange={(e) => setNewRule({ ...newRule, action: e.target.value })}
              >
                <MenuItem value="allow">مجاز</MenuItem>
                <MenuItem value="deny">مسدود</MenuItem>
              </Select>
            </FormControl>

            <FormControl fullWidth>
              <InputLabel id="protocol-label">پروتکل</InputLabel>
              <Select
                labelId="protocol-label"
                value={newRule.protocol}
                label="پروتکل"
                onChange={(e) => setNewRule({ ...newRule, protocol: e.target.value })}
              >
                <MenuItem value="tcp">TCP</MenuItem>
                <MenuItem value="udp">UDP</MenuItem>
                <MenuItem value="any">همه</MenuItem>
              </Select>
            </FormControl>

            <TextField
              label="پورت"
              type="number"
              value={newRule.port}
              onChange={(e) => setNewRule({ ...newRule, port: parseInt(e.target.value) })}
              fullWidth
            />

            <TextField
              label="آدرس منبع (اختیاری)"
              value={newRule.source}
              onChange={(e) => setNewRule({ ...newRule, source: e.target.value })}
              fullWidth
              placeholder="مثال: 192.168.1.0/24"
            />

            <TextField
              label="توضیحات"
              value={newRule.description}
              onChange={(e) => setNewRule({ ...newRule, description: e.target.value })}
              fullWidth
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenDialog(false)} sx={{ fontFamily: 'Vazirmatn' }}>
            لغو
          </Button>
          <Button
            onClick={handleAddRule}
            variant="contained"
            sx={{ fontFamily: 'Vazirmatn' }}
          >
            افزودن
          </Button>
        </DialogActions>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      <Dialog
        open={deleteConfirm !== null}
        onClose={() => setDeleteConfirm(null)}
      >
        <DialogTitle sx={{ fontFamily: 'Vazirmatn' }}>تأیید حذف</DialogTitle>
        <DialogContent>
          <Typography sx={{ fontFamily: 'Vazirmatn' }}>
            آیا از حذف این قانون اطمینان دارید؟
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button
            onClick={() => setDeleteConfirm(null)}
            sx={{ fontFamily: 'Vazirmatn' }}
          >
            لغو
          </Button>
          <Button
            onClick={() => deleteConfirm && handleDeleteRule(deleteConfirm)}
            color="error"
            sx={{ fontFamily: 'Vazirmatn' }}
          >
            حذف
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
};

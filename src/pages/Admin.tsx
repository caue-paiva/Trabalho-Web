/**
 * Admin Dashboard Page (Protected)
 *
 * This page is only accessible to authenticated users.
 * Demonstrates protected route functionality and auth state access.
 */

import { useAuth } from '@/auth';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { useNavigate } from 'react-router-dom';
import { Shield, User, Mail, Key, LogOut } from 'lucide-react';
import { useLanguage } from '@/hooks/useLanguage';

const Admin = () => {
  const { currentUser, logout, isMockMode } = useAuth();
  const navigate = useNavigate();
  const { t } = useLanguage();

  const handleLogout = async () => {
    try {
      await logout();
      navigate('/login');
    } catch (error) {
      console.error(t('admin.alerts.logoutFailed'), error);
    }
  };

  // Get user email safely
  const userEmail = currentUser && 'email' in currentUser ? currentUser.email : 'Unknown';
  const userId = currentUser && 'uid' in currentUser ? currentUser.uid : 'Unknown';
  const displayName = currentUser && 'displayName' in currentUser ? currentUser.displayName : null;
  const providerId = currentUser && 'providerId' in currentUser ? currentUser.providerId : 'Unknown';

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800 p-8">
      <div className="max-w-4xl mx-auto space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold flex items-center gap-2">
              <Shield className="h-8 w-8 text-primary" />
              {t('admin.title')}
            </h1>
            <p className="text-muted-foreground mt-1">{t('admin.subtitle')}</p>
          </div>
          <Button variant="destructive" onClick={handleLogout}>
            <LogOut className="mr-2 h-4 w-4" />
            {t('admin.logout')}
          </Button>
        </div>

        {/* Mock Mode Warning */}
        {isMockMode && (
          <Alert>
            <Shield className="h-4 w-4" />
            <AlertTitle>{t('admin.mockMode.title')}</AlertTitle>
            <AlertDescription>
              {t('admin.mockMode.description')}
            </AlertDescription>
          </Alert>
        )}

        {/* User Information Card */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <User className="h-5 w-5" />
              {t('admin.userInfo.title')}
            </CardTitle>
            <CardDescription>{t('admin.userInfo.description')}</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <div className="flex items-center gap-2 text-sm font-medium text-muted-foreground">
                  <Mail className="h-4 w-4" />
                  {t('admin.userInfo.email')}
                </div>
                <p className="text-lg font-mono">{userEmail || t('admin.userInfo.notAvailable')}</p>
              </div>

              <div className="space-y-2">
                <div className="flex items-center gap-2 text-sm font-medium text-muted-foreground">
                  <User className="h-4 w-4" />
                  {t('admin.userInfo.displayName')}
                </div>
                <p className="text-lg">{displayName || t('admin.userInfo.notSet')}</p>
              </div>

              <div className="space-y-2">
                <div className="flex items-center gap-2 text-sm font-medium text-muted-foreground">
                  <Key className="h-4 w-4" />
                  {t('admin.userInfo.userId')}
                </div>
                <p className="text-sm font-mono break-all">{userId}</p>
              </div>

              <div className="space-y-2">
                <div className="flex items-center gap-2 text-sm font-medium text-muted-foreground">
                  <Shield className="h-4 w-4" />
                  {t('admin.userInfo.provider')}
                </div>
                <Badge variant="secondary">{providerId}</Badge>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default Admin;

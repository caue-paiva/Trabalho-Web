/**
 * Admin Dashboard Page (Protected)
 *
 * This page is only accessible to authenticated users.
 * Demonstrates protected route functionality and auth state access.
 */

import { useAuth } from '@/auth';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Separator } from '@/components/ui/separator';
import { useNavigate } from 'react-router-dom';
import { useState } from 'react';
import { fetchWithAuth, post } from '@/auth/httpClient';
import { Shield, User, Mail, Key, LogOut, TestTube } from 'lucide-react';

const Admin = () => {
  const { currentUser, logout, getToken, isMockMode } = useAuth();
  const navigate = useNavigate();
  const [testResult, setTestResult] = useState<string | null>(null);
  const [testing, setTesting] = useState(false);

  const handleLogout = async () => {
    try {
      await logout();
      navigate('/login');
    } catch (error) {
      console.error('Logout failed:', error);
    }
  };

  const testAuthenticatedRequest = async () => {
    setTesting(true);
    setTestResult(null);

    try {
      const response = await fetchWithAuth('https://jsonplaceholder.typicode.com/posts/1');
      const data = await response.json();
      setTestResult(`✅ Success! Fetched post: "${data.title}"`);
    } catch (error) {
      setTestResult(`❌ Failed: ${error instanceof Error ? error.message : 'Unknown error'}`);
    } finally {
      setTesting(false);
    }
  };

  const showToken = async () => {
    try {
      const token = await getToken();
      if (!token) {
        alert('No token available. Please log in first.');
        return;
      }

      // Print the entire token to console
      console.log('Current ID Token (Full):', token);
      console.log('Token Length:', token.length);

      // Copy the entire token to clipboard
      try {
        await navigator.clipboard.writeText(token);
        alert(`Full token copied to clipboard and logged to console!\n\nToken length: ${token.length} characters\n\nCheck the browser console (F12) to see the full token.`);
      } catch (clipboardError) {
        // Fallback if clipboard API is not available
        console.error('Failed to copy to clipboard:', clipboardError);
        alert(`Token logged to console!\n\nFull Token:\n${token}\n\nToken length: ${token.length} characters\n\nNote: Could not copy to clipboard automatically. Please copy from console or this alert.`);
      }
    } catch (error) {
      console.error('Failed to get token:', error);
      alert('Failed to get token. Check console for details.');
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
              Admin Dashboard
            </h1>
            <p className="text-muted-foreground mt-1">Protected route - authentication required</p>
          </div>
          <Button variant="destructive" onClick={handleLogout}>
            <LogOut className="mr-2 h-4 w-4" />
            Logout
          </Button>
        </div>

        {/* Mock Mode Warning */}
        {isMockMode && (
          <Alert>
            <Shield className="h-4 w-4" />
            <AlertTitle>🎭 Mock Authentication Mode</AlertTitle>
            <AlertDescription>
              You're using simulated authentication. All operations are logged to the console.
              Check the browser console to see authentication events and token generation.
            </AlertDescription>
          </Alert>
        )}

        {/* User Information Card */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <User className="h-5 w-5" />
              User Information
            </CardTitle>
            <CardDescription>Your current authentication details</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <div className="flex items-center gap-2 text-sm font-medium text-muted-foreground">
                  <Mail className="h-4 w-4" />
                  Email
                </div>
                <p className="text-lg font-mono">{userEmail || 'Not available'}</p>
              </div>

              <div className="space-y-2">
                <div className="flex items-center gap-2 text-sm font-medium text-muted-foreground">
                  <User className="h-4 w-4" />
                  Display Name
                </div>
                <p className="text-lg">{displayName || 'Not set'}</p>
              </div>

              <div className="space-y-2">
                <div className="flex items-center gap-2 text-sm font-medium text-muted-foreground">
                  <Key className="h-4 w-4" />
                  User ID
                </div>
                <p className="text-sm font-mono break-all">{userId}</p>
              </div>

              <div className="space-y-2">
                <div className="flex items-center gap-2 text-sm font-medium text-muted-foreground">
                  <Shield className="h-4 w-4" />
                  Provider
                </div>
                <Badge variant="secondary">{providerId}</Badge>
              </div>
            </div>

            <Separator />

            <div className="flex gap-2">
              <Button variant="outline" onClick={showToken} className="flex-1">
                <Key className="mr-2 h-4 w-4" />
                Show Token in Console
              </Button>
            </div>
          </CardContent>
        </Card>

        {/* API Testing Card */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <TestTube className="h-5 w-5" />
              Test Authenticated Requests
            </CardTitle>
            <CardDescription>
              Test the HTTP client with automatic token injection
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-sm text-muted-foreground">
              This will make a request to a public API (JSONPlaceholder) with your authentication
              token automatically injected in the Authorization header. Check the console to see
              the full request details and token preview.
            </p>

            {testResult && (
              <Alert variant={testResult.startsWith('✅') ? 'default' : 'destructive'}>
                <AlertDescription>{testResult}</AlertDescription>
              </Alert>
            )}
          </CardContent>
          <CardFooter>
            <Button
              onClick={testAuthenticatedRequest}
              disabled={testing}
              className="w-full"
            >
              {testing ? 'Testing...' : 'Test API Request with Auth Token'}
            </Button>
          </CardFooter>
        </Card>

        {/* Instructions Card */}
        <Card>
          <CardHeader>
            <CardTitle>🎯 Testing Instructions</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3 text-sm">
            <div className="space-y-2">
              <p className="font-semibold">1. Check Browser Console</p>
              <p className="text-muted-foreground pl-4">
                Open DevTools (F12) and check the Console tab to see all authentication logs
                and token operations.
              </p>
            </div>

            <Separator />

            <div className="space-y-2">
              <p className="font-semibold">2. Test Token Display</p>
              <p className="text-muted-foreground pl-4">
                Click "Show Token in Console" to see your current Firebase ID token.
                This token is automatically included in all authenticated requests.
              </p>
            </div>

            <Separator />

            <div className="space-y-2">
              <p className="font-semibold">3. Test API Request</p>
              <p className="text-muted-foreground pl-4">
                Click "Test API Request" to see how the httpClient automatically injects
                your token into outgoing requests.
              </p>
            </div>

            <Separator />

            <div className="space-y-2">
              <p className="font-semibold">4. Test Session Persistence</p>
              <p className="text-muted-foreground pl-4">
                Refresh this page (F5) - you should remain logged in. The auth state
                is persisted in localStorage.
              </p>
            </div>

            <Separator />

            <div className="space-y-2">
              <p className="font-semibold">5. Test Protected Route</p>
              <p className="text-muted-foreground pl-4">
                Click Logout, then try to access /admin directly. You should be
                redirected to the login page.
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default Admin;

/**
 * Login Page Component
 *
 * Provides a user interface for authentication with:
 * - Email/password login
 * - Google social login
 * - Error handling and loading states
 */

import { useState, FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '@/auth/useAuth';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { FcGoogle } from 'react-icons/fc';
import { Loader2 } from 'lucide-react';

const Login = () => {
  const navigate = useNavigate();
  const { signInWithEmail, signInWithGoogle, isMockMode } = useAuth();

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [googleLoading, setGoogleLoading] = useState(false);

  const handleEmailLogin = async (e: FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      await signInWithEmail(email, password);
      navigate('/');
    } catch (err) {

      if (err instanceof Error) {
        if (err.message.includes('auth/invalid-credential')) {
          setError('Invalid email or password. Please try again.');
        } else if (err.message.includes('auth/user-not-found')) {
          setError('No account found with this email.');
        } else if (err.message.includes('auth/wrong-password')) {
          setError('Incorrect password. Please try again.');
        } else if (err.message.includes('auth/too-many-requests')) {
          setError('Too many failed attempts. Please try again later.');
        } else {
          setError('Failed to sign in. Please try again.');
        }
      } else {
        setError('An unexpected error occurred.');
      }
    } finally {
      setLoading(false);
    }
  };

  const handleGoogleLogin = async () => {
    setError(null);
    setGoogleLoading(true);

    try {
      await signInWithGoogle();
      navigate('/');
    } catch (err) {

      if (err instanceof Error) {
        if (err.message.includes('auth/popup-closed-by-user')) {
          setError('Sign-in popup was closed. Please try again.');
        } else if (err.message.includes('auth/cancelled-popup-request')) {
          // User cancelled, no need to show error
          setError(null);
        } else {
          setError('Failed to sign in with Google. Please try again.');
        }
      } else {
        setError('An unexpected error occurred with Google sign-in.');
      }
    } finally {
      setGoogleLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800 p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl font-bold text-center">
            Admin Login
          </CardTitle>
          <CardDescription className="text-center">
            Sign in to access the admin dashboard
          </CardDescription>
          {isMockMode && (
            <Alert className="mt-4">
              <AlertDescription className="text-sm">
                <strong>🎭 Mock Mode Active</strong>
                <br />
                Use: admin@example.com / admin123
                <br />
                Or: test@example.com / test123
              </AlertDescription>
            </Alert>
          )}
        </CardHeader>

        <form onSubmit={handleEmailLogin}>
          <CardContent className="space-y-4">
            {error && (
              <Alert variant="destructive">
                <AlertDescription>{error}</AlertDescription>
              </Alert>
            )}

            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                placeholder="admin@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                disabled={loading || googleLoading}
              />
            </div>

            <div className="space-y-2">
              <Label htmlFor="password">Password</Label>
              <Input
                id="password"
                type="password"
                placeholder="••••••••"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                disabled={loading || googleLoading}
              />
            </div>
          </CardContent>

          <CardFooter className="flex flex-col space-y-4">
            <Button
              type="submit"
              className="w-full"
              disabled={loading || googleLoading}
            >
              {loading ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Signing in...
                </>
              ) : (
                'Sign in with Email'
              )}
            </Button>

            <div className="relative w-full">
              <div className="absolute inset-0 flex items-center">
                <span className="w-full border-t" />
              </div>
              <div className="relative flex justify-center text-xs uppercase">
                <span className="bg-card px-2 text-muted-foreground">
                  Or continue with
                </span>
              </div>
            </div>

            <Button
              type="button"
              variant="outline"
              className="w-full"
              onClick={handleGoogleLogin}
              disabled={loading || googleLoading}
            >
              {googleLoading ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Signing in...
                </>
              ) : (
                <>
                  <FcGoogle className="mr-2 h-5 w-5" />
                  Sign in with Google
                </>
              )}
            </Button>

            <p className="text-xs text-center text-muted-foreground mt-4">
              By signing in, you agree to our terms and conditions
            </p>
          </CardFooter>
        </form>
      </Card>
    </div>
  );
};

export default Login;

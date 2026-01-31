import { useEffect, useState } from 'react';
import { Box, Button, Link } from '@chakra-ui/react';
import { useGoogleLogin } from '@react-oauth/google';
import { useNavigate } from 'react-router-dom';

export const HomeScreen = ({ email, setEmail, name, setName }) => {
  const navigate = useNavigate();

  const login = useGoogleLogin({
    flow: 'auth-code',
    onSuccess: async (codeResponse) => {
      const res = await fetch('http://localhost:8080/api/auth/google', {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ code: codeResponse.code }),
      });
      const data = await res.json();
      if (!res.ok) {
        console.error('Login failed:', data);
        return;
      }
      const user = await fetch('http://localhost:8080/user-info', {
        credentials: 'include',
      });
      if (user.ok) {
        const data = await user.json();
        console.log('Fetched user info:', data);
        setEmail(data.email ?? '');
        setName(data.name ?? '');
      } else {
        console.error('Failed to fetch user-info');
        console.log(user);
      }

      navigate('/projects');
    },
  });

  const logout = async () => {
    const res = await fetch('http://localhost:8080/logout', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ yo: 'hey' }),
    });
    if (res.ok) {
      setEmail('');
      setName('');
      console.log('Logout successful');
    } else {
      console.error('Logout failed');
    }
  };

  return (
    <Box sx={{ marginLeft: '2%', marginRight: '2%' }}>
      <Box sx={{ fontSize: '100px', textAlign: 'left' }}>
        Constructifyer
        <br />
        The Leading Expense Management Tool
      </Box>
      <Box></Box>
      {/* <Box>
        
        </Box> */}
      <div>
        {email && name ? (
          <div>
            <p>Logged in as: {name}</p>
            <p>Email: {email}</p>
            <Button onClick={() => logout()}>Logout</Button>
            <Button>
              <Link href='/projects'>Your Projects</Link>
            </Button>
          </div>
        ) : (
          <div>
            <Button onClick={() => login()}>Login with Google</Button>
          </div>
        )}

        {/* <button onClick={callProtected}>Call protected route</button> */}
        {/* <button onClick={logout}>Logout</button> */}
      </div>
    </Box>
  );
};

import { useEffect, useState } from 'react';
import { Box, Button, Link } from '@chakra-ui/react';
import { useGoogleLogin } from '@react-oauth/google';

export const HomeScreen = ({ email, setEmail, name, setName, jwt, setJwt }) => {
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
      // console.log(data); Returns a success message
    },
  });

  const callProtected = async () => {
    const user = await fetch(
      'http://localhost:8080/user-info',
      // {
      //   headers: { Authorization: `Bearer ${jwt}` },
      // },
      { credentials: 'include' },
    );
    if (user.ok) {
      const data = await user.json();
      console.log('Protected data:', data);
      setEmail(data.email ?? '');
      setName(data.name ?? '');
    } else {
      console.error('Failed to fetch protected data');
      console.log(user);
    }
  };

  const logout = () => {
    setJwt(null);
    setEmail('');
    setName('');
  };

  return (
    <Box sx={{ marginLeft: '2%', marginRight: '2%' }}>
      <Box sx={{ fontSize: '100px', textAlign: 'left' }}>
        Constructifyer
        <br />
        The Leading Expense Management Tool
      </Box>
      <Box></Box>
      <Box>Ready to get started? Login!</Box>
      {/* <Box>
        <Button>
          <Link href='/projects'>Login</Link>
        </Button>
        </Box> */}
      <div>
        <p>Logged in as: {name}</p>
        <Button onClick={() => login()}>Login with Google</Button>
        <button onClick={callProtected}>Call protected route</button>
        {/* <button onClick={logout}>Logout</button> */}
      </div>
    </Box>
  );
};

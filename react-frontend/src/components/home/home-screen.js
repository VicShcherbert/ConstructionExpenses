import { useEffect } from 'react';
import { Box, Button, Link } from '@chakra-ui/react';
import { useGoogleLogin } from '@react-oauth/google';

export const HomeScreen = () => {
  // useEffect(() => {
  //     fetch('http://localhost:8080/get-projects')
  //     .then((response) => response.json)
  //     .then((responseJson) => {
  //         console.log(responseJson)
  //     })
  //     .catch((error) => console.error(error))
  // }, [])

  const login = useGoogleLogin({
  flow: 'auth-code',
  onSuccess: async (codeResponse) => {
    const res = await fetch('http://localhost:8080/api/auth/google', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ code: codeResponse.code }),
    });
    const data = await res.json();
    console.log('Login Success:', data);
  },
});

  return (
    <Box sx={{ marginLeft: '2%', marginRight: '2%' }}>
      <Box sx={{ fontSize: '100px', textAlign: 'left' }}>
        Constructifyer
        <br />
        The Leading Expense Management Tool
      </Box>
      <Box></Box>
      <Box>Ready to get started? Login!</Box>
      <Box>
        <Button>
          <Link href='/projects'>Login</Link>
        </Button>
        <Button onClick={() => login()}>Login with Google</Button>
      </Box>
    </Box>
  );
};

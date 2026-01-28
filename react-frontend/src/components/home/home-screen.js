import { useEffect, useState } from 'react';
import { Box, Button, Link } from '@chakra-ui/react';
import { useGoogleLogin } from '@react-oauth/google';

export const HomeScreen = ({ email, setEmail, name, setName, jwt, setJwt }) => {
  // const [data, setData] = useState({});
  

  console.log(email);
  console.log(name);
  console.log(jwt);

  const login = useGoogleLogin({
    flow: 'auth-code',
    onSuccess: async (codeResponse) => {
      const res = await fetch('http://localhost:8080/api/auth/google', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ code: codeResponse.code }),
      });
      const data = await res.json();
      setEmail(data.email);
      setName(data.name);
      setJwt(data.jwt);
      console.log(data);
    },
  });

  const callProtected = async () => {
    const res = await fetch('http://localhost:8080/api/protected', {
      headers: { Authorization: `Bearer ${jwt}` },
    });
    const data = await res.json();
    alert(JSON.stringify(data));
  };

  const logout = () => {
    setJwt(null);
    setEmail('');
    setName('');
  };

  // useEffect(() => {
  //   if(Object.keys(data).length === 0) return;
  //   else {
  //     setUser(data);
  //   }
  // }, [data]);

  // console.log(data);

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
        <Button onClick={() => login()}>Login with Google</Button>
      </Box> */}
      <Box>
        {jwt !== null ? (
          <div>
            <p>Logged in as: {name}</p>
            {/* <button onClick={callProtected}>Call protected route</button> */}
            <button onClick={logout}>Logout</button>
          </div>
        ) : (
          <button onClick={() => login()}>Login with Google</button>
        )}
      </Box>
    </Box>
  );
};

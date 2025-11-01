import { useEffect } from 'react';
import { Box, Button, Link } from '@chakra-ui/react';

export const HomeScreen = () => {
  // useEffect(() => {
  //     fetch('http://localhost:8080/get-projects')
  //     .then((response) => response.json)
  //     .then((responseJson) => {
  //         console.log(responseJson)
  //     })
  //     .catch((error) => console.error(error))
  // }, [])

  return (
    <Box>
      <Box>
        Welcome to Constructifyer, the world's best construction expense tracker
      </Box>
      <Box>
        <Button>
            <Link href='/projects'>
                Login
            </Link>
        </Button>
      </Box>
    </Box>
  );
}
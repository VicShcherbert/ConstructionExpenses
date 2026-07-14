import React, { useEffect, useState } from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router';
// import ReactDOM from 'react-dom/client';
import { createRoot } from 'react-dom/client';
import './index.css';
import { HomeScreen } from './components/home/home-screen';
import { Projects } from './components/projects-home/projects';
import { Project } from './components/project/project';
import { Box, Button, ChakraProvider } from '@chakra-ui/react';
import { SiteHeader } from './components/home/site-header';
import { GoogleOAuthProvider } from '@react-oauth/google';
import { RequireAuth } from './components/other/require-auth';
const GOOGLE_CLIENT_ID = process.env.REACT_APP_GOOGLE_CLIENT_ID;

export const App = () => {
  const [email, setEmail] = useState('');
  const [name, setName] = useState('');
  const [userId, setUserId] = useState(0);
  const [authLoading, setAuthLoading] = useState(false);

  console.log('userId in Home: ' + userId);
  console.log('name in Home: ' + name);
  console.log('email in Home: ' + email);

  useEffect(() => {
    (async () => {
      setAuthLoading(true);
      try {
        // On app load, check if user is already authenticated
        const res = await fetch('http://localhost:8080/user-info', {
          method: 'GET',
          credentials: 'include',
        });
        if (res.ok) {
          const data = await res.json();
          setEmail(data.email);
          setName(data.name);
          setUserId(data.user_id);
          // console.log('User is already authenticated:', data);
        } else {
          setName('');
          setEmail('');
          setUserId(null);
          console.log('User is not authenticated');
        }
      } finally {
        setAuthLoading(false);
      }
    })();
  }, []);

  return (
    <ChakraProvider>
      <Router>
        <Box textAlign={'center'} mt={'20px'}>
          {/* <SiteHeader /> */}
          <Routes>
            <Route
              path='/'
              element={
                <HomeScreen
                  email={email}
                  setEmail={setEmail}
                  name={name}
                  setName={setName}
                  userId={userId}
                  setUserId={setUserId}
                />
              }
            />
            <Route
              path='/home'
              element={
                <HomeScreen
                  email={email}
                  setEmail={setEmail}
                  name={name}
                  setName={setName}
                  userId={userId}
                  setUserId={setUserId}
                />
              }
            />
            <Route
              element={
                <RequireAuth isAuthed={!!email} authLoading={authLoading} />
              }
            >
              <Route
                path='/projects'
                element={
                  <Projects
                    email={email}
                    name={name}
                    userId={userId}
                  />
                }
              />
              <Route
                path='/projects/:project_id'
                element={
                  <Project
                    email={email}
                    name={name}
                    userId={userId}
                  />
                }
              />
            </Route>
          </Routes>
        </Box>
      </Router>
    </ChakraProvider>
  );
};

const root = createRoot(document.getElementById('root'));

root.render(
  <GoogleOAuthProvider clientId={GOOGLE_CLIENT_ID}>
    <App />
  </GoogleOAuthProvider>,
);

// ReactDOM.render(<App />, document.getElementById('root'));

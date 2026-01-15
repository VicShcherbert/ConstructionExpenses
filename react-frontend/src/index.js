import React, { useEffect } from 'react';
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

const GOOGLE_CLIENT_ID = process.env.REACT_APP_GOOGLE_CLIENT_ID;

export const App = () => {
  return (
    <ChakraProvider>
      <Router>
        <Box textAlign={'center'} mt={'20px'}>
          <SiteHeader />
          <Routes>
            <Route path='/' element={<HomeScreen />} />
            <Route path='/home' element={<HomeScreen />} />
            <Route path='/projects' element={<Projects />} />
            <Route path='/projects/:project_id' element={<Project />} />
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
  </GoogleOAuthProvider>
);

// ReactDOM.render(<App />, document.getElementById('root'));

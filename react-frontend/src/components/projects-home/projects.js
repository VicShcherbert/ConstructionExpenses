import {
  Box,
  Button,
  Heading,
  VStack,
  FormControl,
  FormLabel,
  Input,
  Flex,
  Stack,
  Link,
  Collapse,
  Spinner,
} from '@chakra-ui/react';
import { ArrowLeftIcon } from '@chakra-ui/icons';
import React, { useEffect, useState } from 'react';
import { useFormik } from 'formik';

export const Projects = () => {
  const [projects, setProjects] = useState([]);
  const [openAddProject, setOpenAddProject] = useState(false);
  const [loader, setLoader] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  console.log('Search query:', searchQuery);

  useEffect(() => {
    (async () => {
      const url = 'http://localhost:8080/get-projects';
      try {
        const response = await fetch(url);
        if (!response.ok) {
          throw new Error(`Response status: ${response.status}`);
        }
        const result = await response.json();
        setProjects(result);
      } catch (error) {
        console.error(error.message);
      }
    })();
  }, []);

  const openProjectPage = (project_id) => {
    const url = '/projects/' + project_id;
    window.open(url, '_self');
  };

  const formik = useFormik({
    initialValues: {
      project_name: '',
    },
    onSubmit: (value) => {
      setLoader(true);
      (async () => {
        await fetch('http://localhost:8080/create-project', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(value),
        })
          .then((response) => {
            if (!response.ok) {
              throw new Error(`Response status: ${response.status}`);
            }
            return response.json();
          })
          .then((data) => {
            console.log('Success:', data);
            setLoader(false);
            window.location.reload();
          })
          .catch((error) => {
            console.error('Error:', error);
          });
      })();
    },
  });

  const handleClick = () => {
    setOpenAddProject(!openAddProject);
  };

  const handleChange = (event) => {
    setSearchQuery(event.target.value);
  };

  return (
    <Box display={'flex'} alignItems={'center'} justifyContent={'center'}>
      <Box>
        <Stack direction='row' display={'inline-flex'}>
          <Button>
            <Link href='/'>
              <ArrowLeftIcon />
            </Link>
          </Button>
          <Heading>Project Dashboard</Heading>
        </Stack>
        <Box mt={'20px'}>
          <Button onClick={handleClick}>
            {openAddProject ? <>Close</> : <>Add Project</>}
          </Button>
          <Collapse in={openAddProject}>
            {!loader ? (
              <Box rounded={'md'} p={'20px'}>
                {/* <Heading size='lg'>Create project</Heading> */}
                <form onSubmit={formik.handleSubmit}>
                  <VStack>
                    {/* <Form> */}
                      <FormControl>
                        <FormLabel htmlFor='project_name'>
                          Project Name
                        </FormLabel>
                        <Input
                          id='project_name'
                          name='project_name'
                          type='project_name'
                          variant='filled'
                          onChange={formik.handleChange}
                          value={formik.values.project_name}
                        />
                      </FormControl>
                      <Button type='submit'>Submit</Button>
                    {/* </Form> */}
                  </VStack>
                </form>
              </Box>
            ) : (
              <Spinner />
            )}
          </Collapse>
        </Box>
        <Box>
          <Input value={searchQuery} onChange={handleChange} placeholder='Search for project'/>
        </Box>
        <Box mt={'20px'}>
          {projects.filter(project => project.project_name.includes(searchQuery)).map((project) => (
            <Box key={project.project_id} mb={'13px'}>
              <Button
                width={'400px'}
                onClick={() => openProjectPage(project.project_id)}
              >
                {project.project_name}
              </Button>
            </Box>
          ))}
        </Box>
      </Box>
    </Box>
  );
};

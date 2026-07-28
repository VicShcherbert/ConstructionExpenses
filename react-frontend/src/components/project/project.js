import { useState, useEffect } from 'react';
import { useParams } from 'react-router';
import {
  Spinner,
  Heading,
  Box,
  TableContainer,
  Thead,
  Table,
  Tr,
  Th,
  Tbody,
  Td,
  VStack,
  FormControl,
  FormLabel,
  Input,
  Button,
  Collapse,
  Stack,
  Link,
} from '@chakra-ui/react';
import { ArrowLeftIcon } from '@chakra-ui/icons';
import { Form, Formik } from 'formik';
import { fetchExpenses, fetchProjectInfo } from '../../api/api';

export const Project = () => {
  const params = useParams();
  const [project, setProject] = useState([]);
  const [expenses, setExpenses] = useState([]);
  const [expenseToUpdate, setExpenseToUpdate] = useState({});
  const [receipt, setReceipt] = useState(null);
  const [receiptURL, setReceiptURL] = useState('');
  const [expenseID, setExpenseID] = useState(0);
  const [expensesLoaded, setExpensesLoaded] = useState(false);
  const [loader, setLoader] = useState(false);
  const [showAddExpense, setShowAddExpense] = useState(false);
  const [totalExpenseAmount, setTotalExpenseAmount] = useState(0);
  const [searchQuery, setSearchQuery] = useState('');

  useEffect(() => {
    (async () => {
      try {
        setExpensesLoaded(false);
        
        const project = await fetchProjectInfo(params.project_id);
        setProject(project);
        const projectExpenses = await fetchExpenses(params.project_id)
        setExpenses(projectExpenses)

        setExpensesLoaded(true);
      } catch (error) {
        console.error(error.message);
      }
    })();
  }, [loader]);

  useEffect(() => {
    (async () => {
      try {
        if (expenseID !== 0 && receipt != null) {
          setLoader(true);
          const formData = new FormData();
          formData.append('file', receipt);

          const expenseUploadURL =
            'http://localhost:8080/upload-receipt/' + expenseID;
          const res = await fetch(expenseUploadURL, {
            method: 'POST',
            body: formData,
          });

          const data = await res.json();
          setReceiptURL(data.url);
        } else if (expenseID !== 0 && receipt == null) {
          setLoader(false);
        }
      } catch (error) {
        console.error(error.message);
      }
    })();
  }, [expenseID]);

  useEffect(() => {
    (async () => {
      if (receiptURL !== '' && expenseID !== 0 && expenses !== null) {
        setLoader(true);
        const updateInfo = await fetch(
          `http://localhost:8080/update-expense/${expenseID}`,
          {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
              expense_id: expenseToUpdate.expense_id,
              expense_name: expenseToUpdate.expense_name,
              expense_receipt_url: receiptURL,
              expense_cost: expenseToUpdate.expense_cost,
              project_id: expenseToUpdate.project_id,
            }),
          },
        );
        setLoader(false);
        await fetchExpenses();
      }
    })();
  }, [receiptURL]);

  useEffect(() => {
    if (expenses !== null && expenses.length !== 0) {
      setTotalExpenseAmount(
        expenses.reduce(
          (accumulator, currentValue) =>
            accumulator + currentValue.expense_cost,
          0,
        ),
      );
    }
  }, [expenses]);

  const handleClick = () => {
    setShowAddExpense(!showAddExpense);
  };

  const handleDelete = async (expenseID) => {
    setLoader(true);
    await fetch(`http://localhost:8080/delete-expense/${expenseID}`, {
      method: 'DELETE',
    })
      .then((res) => res.json())
      .then((data) => {});
    setLoader(false);
  };

  const handleChange = (event) => {
    setSearchQuery(event.target.value);
  };

  return (
    <Box display={'flex'} alignItems={'center'} justifyContent={'center'}>
      {project.length !== 0 && expensesLoaded === true ? (
        <Box>
          <Stack direction='row' display={'inline-flex'}>
            <Button>
              <Link href='/projects'>
                <ArrowLeftIcon />
              </Link>
            </Button>
            <Heading>{project[0].project_name}</Heading>
          </Stack>
          <Box mt={'20px'}>
            <Button onClick={handleClick}>
              {showAddExpense ? <>Close</> : <>Add Expense</>}
            </Button>
            <Collapse in={showAddExpense}>
              {!loader ? (
                <Box rounded={'md'} p={'20px'}>
                  <Formik
                    initialValues={{
                      expense_name: '',
                      expense_cost: 0,
                      expense_receipt_url: '',
                      project_id: project[0].project_id,
                    }}
                    onSubmit={(value) => {
                      setLoader(true);
                      (async () => {
                        await fetch('http://localhost:8080/create-expense', {
                          method: 'POST',
                          headers: {
                            'Content-Type': 'application/json',
                          },
                          body: JSON.stringify(value),
                        })
                          .then((response) => {
                            if (!response.ok) {
                              throw new Error(
                                `Response status: ${response.status}`,
                              );
                            }
                            return response.json();
                          })
                          .then((data) => {
                            setExpenseID(data.expense_id);
                            setExpenseToUpdate(data);
                          })
                          .catch((error) => {
                            console.error('Error: ', error);
                          });
                      })();
                    }}
                  >
                    {({ values, setFieldValue }) => (
                      <VStack>
                        <Form>
                          <FormControl>
                            <FormLabel htmlFor='expense_name'>
                              Expense Name
                            </FormLabel>
                            <Input
                              id='expense_name'
                              name='expense_name'
                              type='expense_name'
                              variant='filled'
                              onChange={(e) => {
                                setFieldValue('expense_name', e.target.value);
                              }}
                              value={values.expense_name}
                            />
                          </FormControl>
                          <FormControl>
                            <FormLabel htmlFor='expense_cost'>
                              Expense Cost
                            </FormLabel>
                            <Input
                              id='expense_cost'
                              name='expense_cost'
                              type='expense_cost'
                              variant='filled'
                              onChange={(e) => {
                                setFieldValue(
                                  'expense_cost',
                                  e.target.value === ''
                                    ? ''
                                    : Number(e.target.value),
                                );
                              }}
                              value={values.expense_cost}
                            />
                          </FormControl>
                          <FormControl>
                            <FormLabel htmlFor='upload_receipt'>
                              Upload Receipt
                            </FormLabel>
                            <Input
                              type='file'
                              onChange={(e) => setReceipt(e.target.files[0])}
                            />
                          </FormControl>
                          <Button
                            type='submit'
                            colorScheme='gray'
                            width='full'
                            mt={'15px'}
                          >
                            Submit
                          </Button>
                        </Form>
                      </VStack>
                    )}
                  </Formik>
                </Box>
              ) : (
                <Spinner />
              )}
            </Collapse>
          </Box>
          <Box>
            <Input
              value={searchQuery}
              onChange={handleChange}
              placeholder='Search for expenses'
              mt={'20px'}
            />
          </Box>
          <Box>
            {expenses !== null ? (
              <Box mt={'20px'}>
                <Box bg='orange' color='white' p={'10px'} rounded={'md'}>
                  <Heading size={'md'}>
                    Total Expenses: ${totalExpenseAmount}
                  </Heading>
                </Box>
                <Heading size='lg' mt={'20px'}>
                  Expenses
                </Heading>
                <TableContainer>
                  <Table variant='simple'>
                    <Thead>
                      <Tr>
                        <Th>Expense Name</Th>
                        <Th>Expense Cost</Th>
                        <Th>Receipt Link</Th>
                        <Th>Delete</Th>
                      </Tr>
                    </Thead>
                    <Tbody>
                      {expenses
                        .filter((expense) =>
                          expense.expense_name
                            .toLowerCase()
                            .includes(searchQuery.toLowerCase()),
                        )
                        .map((expense) => (
                          <Tr key={expense.expense_id}>
                            <Td>{expense.expense_name}</Td>
                            <Td>{expense.expense_cost}</Td>
                            <Td>
                              {expense.expense_receipt_url !== '' ? (
                                <Button>
                                  <a
                                    href={`${expense.expense_receipt_url}`}
                                    target='_blank'
                                    rel='noreferrer'
                                  >
                                    View Receipt
                                  </a>
                                </Button>
                              ) : null}
                            </Td>
                            <Td>
                              <Button
                                colorScheme='red'
                                onClick={() => handleDelete(expense.expense_id)}
                              >
                                Delete
                              </Button>
                            </Td>
                          </Tr>
                        ))}
                    </Tbody>
                  </Table>
                </TableContainer>
              </Box>
            ) : (
              <Box mt={'20px'}>No project expenses yet</Box>
            )}
          </Box>
        </Box>
      ) : (
        <Spinner size={'xl'} />
      )}
    </Box>
  );
};

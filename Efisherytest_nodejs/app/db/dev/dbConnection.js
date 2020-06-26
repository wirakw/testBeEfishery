import pool from './pool';

pool.on('connect', () => {
  console.log('connected to the db');
});

/**
 * Create User Table
 */
const createUserTable = () => {
  const userCreateQuery = `CREATE TABLE IF NOT EXISTS users
  (user_id SERIAL PRIMARY KEY, 
  name TEXT, 
  role TEXT,
  phone TEXT, 
  password TEXT)`;

  pool.query(userCreateQuery)
    .then((res) => {
      console.log(res);
      pool.end();
    })
    .catch((err) => {
      console.log(err);
      pool.end();
    });
};

/**
 * Create All Tables
 */
const createAllTables = () => {
  createUserTable();
};

pool.on('remove', () => {
  console.log('client removed');
  process.exit(0);
});


export {
  createAllTables,
};

require('make-runnable');

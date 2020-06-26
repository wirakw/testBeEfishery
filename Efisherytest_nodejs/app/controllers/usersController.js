/* eslint-disable camelcase */

import dbQuery from '../db/dev/dbQuery';

import {
  comparePassword,
  isEmpty,
  generateUserToken,
  generatePassword,
  validatePassword
} from '../helpers/validations';

import {
  errorMessage, successMessage, status,
} from '../helpers/status';

/**
   * Create A User
   * @param {object} req
   * @param {object} res
   * @returns {object} reflection object
   */
const createUser = async (req, res) => {
  const {
    name, role, phone
  } = req.body;
  if (isEmpty(name) || isEmpty(role) || isEmpty(phone)) {
    errorMessage.error = 'Email, password, first name and last name field cannot be empty';
    return res.status(status.bad).send(errorMessage);
  }
  const password = generatePassword(4)
  if (!validatePassword(password)) {
    errorMessage.error = 'Password must be less than four(4) characters';
    return res.status(status.bad).send(errorMessage);
  }
  const signinUserQuery = 'SELECT * FROM users WHERE name = $1';
  const { rows } = await dbQuery.query(signinUserQuery, [name]);
  const data = rows[0];
  if (data) {
    errorMessage.error = 'User with this name already exist';
    return res.status(status.conflict).send(errorMessage);
  }

  const createUserQuery = `INSERT INTO
      users(name, role, phone, password)
      VALUES($1, $2, $3, $4)
      returning *`;
  const values = [
    name,
    role,
    phone,
    password,
  ];

  try {
    const { rows } = await dbQuery.query(createUserQuery, values);
    const dbResponse = rows[0];
    const token = generateUserToken(dbResponse.name, dbResponse.role, dbResponse.password);
        // delete dbResponse.password;
    successMessage.data = dbResponse;
    successMessage.data.token = token;
    return res.status(status.created).send(successMessage);
  } catch (error) {
    if (error.routine === '_bt_check_unique') {
      errorMessage.error = 'User with that name already exist';
      return res.status(status.conflict).send(errorMessage);
    }
    errorMessage.error = 'Operation was not successful';
    return res.status(status.error).send(errorMessage);
  }
};

/**
   * Signin
   * @param {object} req
   * @param {object} res
   * @returns {object} user object
   */
const siginUser = async (req, res) => {
  const { phone, password } = req.body;
  if (isEmpty(phone) || isEmpty(password)) {
    errorMessage.error = 'Phone or Password detail is missing';
    return res.status(status.bad).send(errorMessage);
  }
  if (!validatePassword(password)) {
    errorMessage.error = 'Please enter a valid Email or Password';
    return res.status(status.bad).send(errorMessage);
  }

  const signinUserQuery = 'SELECT * FROM users WHERE phone = $1';
  try {
    const { rows } = await dbQuery.query(signinUserQuery, [phone]);
    const dbResponse = rows[0];
    if (!dbResponse) {
      errorMessage.error = 'User with this phone does not exist';
      return res.status(status.notfound).send(errorMessage);
    }

    if (dbResponse.password != password) {
      errorMessage.error = 'The password you provided is incorrect';
      return res.status(status.bad).send(errorMessage);
    }
    const token = generateUserToken(dbResponse.name, dbResponse.role, dbResponse.password);
    // delete dbResponse.password;
    successMessage.data = dbResponse;
    successMessage.data.token = token;
    return res.status(status.success).send(successMessage);
  } catch (error) {
    errorMessage.error = 'Operation was not successful';
    return res.status(status.error).send(errorMessage);
  }
};

/**
 * @params {Object} req
 * @params {Object} res
 * @returns return firstname and Lastname
 */ 

export {
  createUser,
  siginUser,
};

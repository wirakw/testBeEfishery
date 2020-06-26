/* eslint-disable camelcase */
// import bcrypt from 'bcryptjs';
import jwt from 'jsonwebtoken';
import env from '../../env';

/**
   * comparePassword
   * @param {string} hashPassword
   * @param {string} password
   * @returns {Boolean} return True or False
   */
const comparePassword = (hashedPassword, password) => {
  return password == hashedPassword;
  // return bcrypt.compareSync(password, hashedPassword);
};

/**
   * isEmpty helper method
   * @param {string, integer} input
   * @returns {Boolean} True or False
   */
const isEmpty = (input) => {
  if (input === undefined || input === '') {
    return true;
  }
  if (input.replace(/\s/g, '').length) {
    return false;
  } return true;
};

/**
   * empty helper method
   * @param {string, integer} input
   * @returns {Boolean} True or False
   */
const empty = (input) => {
  if (input === undefined || input === '') {
    return true;
  }
};

/**
   * Generate Token
   * @param {string} id
   * @returns {string} token
   */
const generateUserToken = (name, role, password) => {
  const timestamp = Date.now().toString();
  const token = jwt.sign({
    name,
    role,
    password,
    timestamp,
  },
  env.secret, { expiresIn: '1d' });
  console.log(token);

  return token;
};
/**
   * validatePassword helper method
   * @param {string} password
   * @returns {Boolean} True or False
   */
  const validatePassword = (password) => {
    if (password.length >= 5 || password === '') {
      return false;
    } return true;
  };
const generatePassword = (length) => {
  var result           = '';
  var characters       = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789';
  var charactersLength = characters.length;
  for ( var i = 0; i < length; i++ ) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
 }
 return result;
};


export {
  isEmpty,
  empty,
  generateUserToken,
  generatePassword,
  validatePassword
};

/* eslint-disable max-len */
import jwt from 'jsonwebtoken';
import dotenv from 'dotenv';
import {
  errorMessage, status,
} from '../helpers/status';

import env from '../../env';

dotenv.config();

/**
   * Verify Token
   * @param {object} req 
   * @param {object} res 
   * @param {object} next
   * @returns {object|void} response object 
   */

const verifyToken = async (req, res, next) => {
  const authHeader = req.headers.authorization;
  if (!authHeader) {
    errorMessage.error = 'Token not provided';
    return res.status(status.bad).send(errorMessage);
  }
  const token = authHeader.split(' ')[1];
  try {
    const decoded =  jwt.verify(token, process.env.SECRET);
    req.user = {
      name: decoded.name,
      role: decoded.role,
      password: decoded.password,
      timestamp: decoded.timestamp
    };
    next();
  } catch (error) {
    errorMessage.error = 'Authentication Failed';
    return res.status(status.unauthorized).send(errorMessage);
  }
};

export default verifyToken;

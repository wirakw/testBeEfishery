import {
  errorMessage, successMessage, status,
} from '../helpers/status';
import isEmpty from '../helpers/validations';
/**
   * Create A User
   * @param {object} req
   * @param {object} res
   * @returns {object} reflection object
   */
const showPrivateClaim = async (req, res) => {
  if (!req.user) {
    return res.status(status.bad).send(errorMessage);
  }
  try {
    successMessage.data = req.user;
    return res.status(status.success).send(successMessage);
  } catch (error) {
    errorMessage.error = 'Authentication Failed';
    return res.status(status.unauthorized).send(errorMessage);
  }
};

export default showPrivateClaim;

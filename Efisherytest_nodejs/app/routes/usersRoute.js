import express from 'express';

import { createUser, siginUser } from '../controllers/usersController';

const router = express.Router();

// users Routes

router.post('/registrasi', createUser);
router.post('/login', siginUser);

export default router;

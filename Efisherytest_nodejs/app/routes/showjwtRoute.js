import express from 'express';

import showPrivateClaim from '../controllers/showjwtController';
import verifyAuth from '../middlewares/verifyAuth';
const router = express.Router();

// users Routes

router.get('/secured', verifyAuth, showPrivateClaim);

export default router;

import { Pool } from 'pg';
import env from '../../../env';

const databaseConfig = { connectionString: env.database_url };
const pool = new Pool(databaseConfig);
console.log(env.database_url);
export default pool;

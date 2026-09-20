DROP TRIGGER IF EXISTS trg_prevent_default_super_admin_delete ON users;
DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
DROP TABLE IF EXISTS users;
DROP FUNCTION IF EXISTS prevent_default_super_admin_delete();
DROP FUNCTION IF EXISTS set_updated_at();

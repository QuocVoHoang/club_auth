CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS citext;

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION prevent_default_super_admin_delete()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.id = '00000000-0000-0000-0000-000000000001'::UUID THEN
        RAISE EXCEPTION 'the default super_admin cannot be deleted';
    END IF;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email CITEXT NOT NULL,
    phone VARCHAR(20),
    full_name VARCHAR(150) NOT NULL,
    role SMALLINT NOT NULL DEFAULT 2,
    password TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    last_login TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    created_by UUID,
    updated_by UUID,
    deleted_by UUID,
    CONSTRAINT users_role_check CHECK (role IN (0, 1, 2)),
    CONSTRAINT users_status_check CHECK (status IN ('active', 'inactive')),
    CONSTRAINT users_created_by_fkey FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
    CONSTRAINT users_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE SET NULL,
    CONSTRAINT users_deleted_by_fkey FOREIGN KEY (deleted_by) REFERENCES users (id) ON DELETE SET NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users (email);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_phone ON users (phone);
CREATE INDEX IF NOT EXISTS idx_users_role ON users (role);
CREATE INDEX IF NOT EXISTS idx_users_status ON users (status);
CREATE INDEX IF NOT EXISTS idx_users_role_status ON users (role, status);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);

DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
CREATE TRIGGER trg_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_prevent_default_super_admin_delete
BEFORE DELETE ON users
FOR EACH ROW
EXECUTE FUNCTION prevent_default_super_admin_delete();

COMMENT ON COLUMN users.role IS '0 = super_admin, 1 = admin, 2 = user';

-- Default credentials must be changed after the first login.
-- Password: VoquocTH2001@@
INSERT INTO users (
    id,
    email,
    full_name,
    role,
    password
)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'superadmin@superadmin.com',
    'Default Super Admin',
    0,
    crypt('VoquocTH2001@@', gen_salt('bf'))
)
ON CONFLICT (id) DO NOTHING;

CREATE TABLE IF NOT EXISTS goadmin_users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL,
    name VARCHAR(100) NOT NULL,
    avatar VARCHAR(255),
    remember_token VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS goadmin_roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS goadmin_permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    http_method VARCHAR(255),
    http_path TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS goadmin_menu (
    id SERIAL PRIMARY KEY,
    parent_id INTEGER DEFAULT 0,
    type INTEGER DEFAULT 1,
    "order" INTEGER DEFAULT 0,
    title VARCHAR(100) NOT NULL,
    plugin_name VARCHAR(100) DEFAULT '',
    header VARCHAR(100),
    icon VARCHAR(50) DEFAULT 'fa-tasks',
    uri VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS goadmin_role_users (
    role_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS goadmin_role_permissions (
    role_id INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS goadmin_role_menu (
    role_id INTEGER NOT NULL,
    menu_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS goadmin_user_permissions (
    user_id INTEGER NOT NULL,
    permission_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS goadmin_session (
    id SERIAL PRIMARY KEY,
    sid VARCHAR(100) NOT NULL UNIQUE,
    "values" TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS goadmin_operation_log (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    path VARCHAR(255) NOT NULL,
    method VARCHAR(20) NOT NULL,
    ip VARCHAR(50) NOT NULL,
    input TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS goadmin_site (
    id SERIAL PRIMARY KEY,
    key VARCHAR(100) NOT NULL UNIQUE,
    value TEXT,
    description VARCHAR(255),
    state INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO goadmin_roles (id, name, slug)
VALUES (1, 'Administrator', 'administrator'), (2, 'Operator', 'operator')
ON CONFLICT (id) DO NOTHING;

INSERT INTO goadmin_permissions (id, name, slug, http_method, http_path)
VALUES (1, 'All permission', '*', '', '*'), (2, 'Dashboard', 'dashboard', 'GET,PUT,POST,DELETE', '/')
ON CONFLICT (id) DO NOTHING;

INSERT INTO goadmin_role_users (role_id, user_id)
VALUES (1, 1), (2, 2)
ON CONFLICT DO NOTHING;

INSERT INTO goadmin_role_permissions (role_id, permission_id)
VALUES (1, 1), (1, 2), (2, 2)
ON CONFLICT DO NOTHING;

INSERT INTO goadmin_menu (id, parent_id, type, "order", title, icon, uri)
VALUES
(1, 0, 1, 2, 'Admin', 'fa-tasks', ''),
(2, 1, 1, 2, 'Users', 'fa-users', '/info/manager'),
(3, 1, 1, 3, 'Roles', 'fa-user', '/info/roles'),
(4, 1, 1, 4, 'Permission', 'fa-ban', '/info/permission'),
(5, 1, 1, 5, 'Menu', 'fa-bars', '/menu'),
(6, 1, 1, 6, 'Operation log', 'fa-history', '/info/op'),
(7, 0, 1, 1, 'Dashboard', 'fa-bar-chart', '/')
ON CONFLICT (id) DO NOTHING;

INSERT INTO goadmin_role_menu (role_id, menu_id)
VALUES (1, 1), (1, 7), (2, 7)
ON CONFLICT DO NOTHING;

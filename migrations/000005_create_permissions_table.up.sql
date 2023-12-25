CREATE TABLE permissions (
  id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
  permission VARCHAR(255) NOT NULL
);

CREATE TABLE users_permissions (
  user_id INTEGER NOT NULL,
  permission_id INTEGER NOT NULL,
  PRIMARY KEY (user_id, permission_id),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
);

INSERT INTO permissions (permission)
VALUES ('parties:read'),
       ('parties:write');

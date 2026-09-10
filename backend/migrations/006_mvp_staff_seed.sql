BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

INSERT INTO expert_groups (title)
SELECT value FROM (VALUES ('Операторы'), ('Психологи'), ('Юристы'), ('Администраторы')) AS groups(value)
WHERE NOT EXISTS (SELECT 1 FROM expert_groups WHERE lower(title) = lower(value));

INSERT INTO users (username, password_hash, role, expert_group_id, full_name, max_tickets)
SELECT seed.username, crypt(seed.password, gen_salt('bf', 10)), seed.role, groups.id, seed.full_name, seed.max_tickets
FROM (VALUES
  ('operator@otklik.local', 'operator123', 'operator', 'Операторы', 'Олег Зетник', 100),
  ('psy@otklik.local', 'expert123', 'expert', 'Психологи', 'Елена Байкова', 10),
  ('law@otklik.local', 'expert123', 'expert', 'Юристы', 'Илья Воронов', 10),
  ('admin@otklik.local', 'admin123', 'admin', 'Администраторы', 'Администратор', 100)
) AS seed(username, password, role, group_title, full_name, max_tickets)
JOIN expert_groups groups ON groups.title = seed.group_title
WHERE NOT EXISTS (SELECT 1 FROM users WHERE users.username = seed.username);

COMMIT;

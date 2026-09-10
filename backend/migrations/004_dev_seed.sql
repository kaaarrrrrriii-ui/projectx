BEGIN;

-- Development-only reference data used by docker-compose. The statements are
-- idempotent and do not alter the database schema.
DO $$
DECLARE
    category_name TEXT;
    v_category_id INTEGER;
    v_question_id INTEGER;
BEGIN
    FOREACH category_name IN ARRAY ARRAY[
        'Травля и оскорбления',
        'Кибербуллинг',
        'Конфликт с родителями',
        'Конфликт с одноклассниками',
        'Конфликт с учителем',
        'Давление и угрозы',
        'Юридический вопрос',
        'Конфликт с сестрой/братом',
        'Не знаю, как это назвать'
    ]
    LOOP
        INSERT INTO categories (name)
        SELECT category_name
        WHERE NOT EXISTS (
            SELECT 1 FROM categories WHERE lower(name) = lower(category_name)
        );

        SELECT id INTO v_category_id
        FROM categories
        WHERE lower(name) = lower(category_name)
        ORDER BY id
        LIMIT 1;

        INSERT INTO questions (cat_id, text)
        SELECT v_category_id, 'Где это происходит?'
        WHERE NOT EXISTS (
            SELECT 1 FROM questions WHERE cat_id = v_category_id AND text = 'Где это происходит?'
        );
        SELECT id INTO v_question_id FROM questions
        WHERE cat_id = v_category_id AND text = 'Где это происходит?'
        ORDER BY id LIMIT 1;
        INSERT INTO answers (question_id, text)
        SELECT v_question_id, answer_text
        FROM unnest(ARRAY['В школе', 'Онлайн', 'В другом месте']) AS answer_text
        WHERE NOT EXISTS (
            SELECT 1 FROM answers WHERE answers.question_id = v_question_id AND answers.text = answer_text
        );

        INSERT INTO questions (cat_id, text)
        SELECT v_category_id, 'Как давно это длится?'
        WHERE NOT EXISTS (
            SELECT 1 FROM questions WHERE cat_id = v_category_id AND text = 'Как давно это длится?'
        );
        SELECT id INTO v_question_id FROM questions
        WHERE cat_id = v_category_id AND text = 'Как давно это длится?'
        ORDER BY id LIMIT 1;
        INSERT INTO answers (question_id, text)
        SELECT v_question_id, answer_text
        FROM unnest(ARRAY['Недавно', 'Давно', 'Пару дней']) AS answer_text
        WHERE NOT EXISTS (
            SELECT 1 FROM answers WHERE answers.question_id = v_question_id AND answers.text = answer_text
        );

        INSERT INTO questions (cat_id, text)
        SELECT v_category_id, 'Вы уже обращались за помощью?'
        WHERE NOT EXISTS (
            SELECT 1 FROM questions WHERE cat_id = v_category_id AND text = 'Вы уже обращались за помощью?'
        );
        SELECT id INTO v_question_id FROM questions
        WHERE cat_id = v_category_id AND text = 'Вы уже обращались за помощью?'
        ORDER BY id LIMIT 1;
        INSERT INTO answers (question_id, text)
        SELECT v_question_id, answer_text
        FROM unnest(ARRAY['Нет', 'Да', 'Не уверен']) AS answer_text
        WHERE NOT EXISTS (
            SELECT 1 FROM answers WHERE answers.question_id = v_question_id AND answers.text = answer_text
        );
    END LOOP;
END $$;

-- Minimal, repeatable staff data for the Docker MVP demo. Passwords are
-- bcrypt hashes for the documented development-only credentials.
INSERT INTO expert_groups (title)
SELECT 'Операторы'
WHERE NOT EXISTS (SELECT 1 FROM expert_groups WHERE lower(title) = lower('Операторы'));

INSERT INTO expert_groups (title)
SELECT 'Психологи и консультанты'
WHERE NOT EXISTS (SELECT 1 FROM expert_groups WHERE lower(title) = lower('Психологи и консультанты'));

INSERT INTO expert_groups (title)
SELECT 'Администраторы'
WHERE NOT EXISTS (SELECT 1 FROM expert_groups WHERE lower(title) = lower('Администраторы'));

INSERT INTO cats_expert_groups (group_id, cat_id)
SELECT groups.id, categories.id
FROM expert_groups AS groups
CROSS JOIN categories
WHERE groups.title = 'Психологи и консультанты'
ON CONFLICT (group_id, cat_id) DO NOTHING;

INSERT INTO users (username, password_hash, role, expert_group_id, full_name, max_tickets)
SELECT
    'operator@otklik.local',
    '$2a$10$vhsmd4y.jrN/ntRomYWHje7nD.QY0xt6URK.zeE31YXFtPnvbcssW',
    'operator',
    id,
    'Олег Зетник',
    0
FROM expert_groups
WHERE title = 'Операторы'
ON CONFLICT (username) DO UPDATE SET
    password_hash = EXCLUDED.password_hash,
    role = EXCLUDED.role,
    expert_group_id = EXCLUDED.expert_group_id,
    full_name = EXCLUDED.full_name,
    max_tickets = EXCLUDED.max_tickets;

INSERT INTO users (username, password_hash, role, expert_group_id, full_name, max_tickets)
SELECT
    account.username,
    '$2a$10$lsjKtD/.KFxiUqqfPoYtc.PuSv.zGid0MoPtjUsynjRBf3DhydDEW',
    'expert',
    groups.id,
    account.full_name,
    10
FROM expert_groups AS groups
CROSS JOIN (VALUES
    ('psy@otklik.local', 'Елена Байкова'),
    ('law@otklik.local', 'Илья Воронов')
) AS account(username, full_name)
WHERE groups.title = 'Психологи и консультанты'
ON CONFLICT (username) DO UPDATE SET
    password_hash = EXCLUDED.password_hash,
    role = EXCLUDED.role,
    expert_group_id = EXCLUDED.expert_group_id,
    full_name = EXCLUDED.full_name,
    max_tickets = EXCLUDED.max_tickets;

INSERT INTO users (username, password_hash, role, expert_group_id, full_name, max_tickets)
SELECT
    'admin@otklik.local',
    '$2a$10$cSHIuCpdxjS.lukZXhF3QertTryOYm1GAuGho1RSmNBdupOegRBle',
    'admin',
    id,
    'Администратор',
    0
FROM expert_groups
WHERE title = 'Администраторы'
ON CONFLICT (username) DO UPDATE SET
    password_hash = EXCLUDED.password_hash,
    role = EXCLUDED.role,
    expert_group_id = EXCLUDED.expert_group_id,
    full_name = EXCLUDED.full_name,
    max_tickets = EXCLUDED.max_tickets;

COMMIT;

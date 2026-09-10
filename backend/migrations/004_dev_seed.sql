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

COMMIT;

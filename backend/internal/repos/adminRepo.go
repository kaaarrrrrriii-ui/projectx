package repos

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"example.com/german/backend/internal/models"
	"github.com/lib/pq"
)

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository { return &AdminRepository{db: db} }

type AdminDashboardRecord struct {
	NewCount, ActiveCount, UrgentCount, ReturnedCount, RoutingIssueCount, ExpertsAtCapacityCount int
}

type AdminEmployeeRecord struct {
	ID, GroupID                          int64
	Username, FullName, Role, GroupTitle string
	AverageRating                        float64
	MaxTickets, ActiveTickets            int
}

type AdminEmployeePageRecord struct {
	Items []AdminEmployeeRecord
	Total int
}

type AdminGroupRecord struct {
	ID                           int64
	Title                        string
	EmployeeCount, CategoryCount int
}

type AdminGroupPageRecord struct {
	Items []AdminGroupRecord
	Total int
}

type AdminAnswerRecord struct {
	ID   int64
	Text string
}

type AdminQuestionRecord struct {
	ID      int64
	Text    string
	Answers []AdminAnswerRecord
}

type AdminCategoryRecord struct {
	ID        int64
	Name      string
	Groups    []AdminGroupRecord
	Questions []AdminQuestionRecord
}

type AdminCategoryPageRecord struct {
	Items []AdminCategoryRecord
	Total int
}

type AdminWorkerRecord struct {
	ID, GroupID           int64
	FullName, GroupTitle  string
	IsResponsible, Actual bool
	CreatedAt             time.Time
}

type AdminTicketRecord struct {
	TrackID         string
	CategoryID      int64
	Category        string
	Status          models.TicketStatus
	ApplicantType   models.ApplicantType
	Priority        models.TicketPriority
	CreatedAt       time.Time
	ClosedAt        sql.NullTime
	ReturnCount     int
	ResponsibleID   sql.NullInt64
	ResponsibleName sql.NullString
	RoutingIssue    sql.NullString
	Workers         []AdminWorkerRecord
}

type AdminTicketFilter struct {
	Search, RoutingIssue      string
	Status                    models.TicketStatus
	Priority                  models.TicketPriority
	ApplicantType             models.ApplicantType
	CategoryID, ResponsibleID int64
	Limit, Offset             int
}

type AdminTicketPageRecord struct {
	Items []AdminTicketRecord
	Total int
}

type AdminEventRecord struct {
	ID                           int64
	EventType                    string
	ActorID                      sql.NullInt64
	ActorName, ActorRole         sql.NullString
	FromStatus, ToStatus         sql.NullInt64
	FromWorkerID, ToWorkerID     sql.NullInt64
	FromWorkerName, ToWorkerName sql.NullString
	ReasonCode                   sql.NullString
	CreatedAt                    time.Time
}

type AdminEventPageRecord struct {
	Items []AdminEventRecord
	Total int
}

type AdminEmployeeLoadRecord struct {
	ID             int64
	FullName, Role string
	ActiveTickets  int
	MaxTickets     int
	ActionCount    int
}

func (repository *AdminRepository) Dashboard(ctx context.Context) (AdminDashboardRecord, error) {
	var result AdminDashboardRecord
	err := repository.db.QueryRowContext(ctx, `
		SELECT
			COUNT(*) FILTER (WHERE t.status = $1),
			COUNT(*) FILTER (WHERE t.status IN ($2,$3,$4,$5)),
			COUNT(*) FILTER (WHERE t.priority = $6 AND t.status NOT IN ($7,$8,$9)),
			COUNT(*) FILTER (WHERE t.status = $10),
			COUNT(*) FILTER (WHERE t.status IN ($1,$10) AND (
				NOT EXISTS (SELECT 1 FROM cats_expert_groups ceg WHERE ceg.cat_id=t.category_id)
				OR NOT EXISTS (
					SELECT 1 FROM users u
					WHERE u.role='expert' AND u.expert_group_id IN (SELECT ceg.group_id FROM cats_expert_groups ceg WHERE ceg.cat_id=t.category_id)
					AND (SELECT COUNT(*) FROM tickets_workers tw JOIN tickets at ON at.id=tw.ticket_id WHERE tw.worker_id=u.id AND tw.actual AND at.status IN ($2,$3,$4,$5)) < u.max_tickets
				)
			))
		FROM tickets t`,
		models.TicketStatusNew, models.TicketStatusAssigned, models.TicketStatusInProgress,
		models.TicketStatusNeedsClarification, models.TicketStatusAnswerReady, models.TicketPriorityUrgent,
		models.TicketStatusCompleted, models.TicketStatusRejected, models.TicketStatusClosedWithoutAnswer,
		models.TicketStatusReturned,
	).Scan(&result.NewCount, &result.ActiveCount, &result.UrgentCount, &result.ReturnedCount, &result.RoutingIssueCount)
	if err != nil {
		return AdminDashboardRecord{}, fmt.Errorf("get admin dashboard: %w", err)
	}
	err = repository.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM (
			SELECT u.id FROM users u
			LEFT JOIN tickets_workers tw ON tw.worker_id=u.id AND tw.actual
			LEFT JOIN tickets t ON t.id=tw.ticket_id AND t.status IN ($1,$2,$3,$4)
			WHERE u.role='expert' GROUP BY u.id,u.max_tickets
			HAVING COUNT(t.id) >= u.max_tickets
		) overloaded`, models.TicketStatusAssigned, models.TicketStatusInProgress,
		models.TicketStatusNeedsClarification, models.TicketStatusAnswerReady,
	).Scan(&result.ExpertsAtCapacityCount)
	if err != nil {
		return AdminDashboardRecord{}, fmt.Errorf("count experts at capacity: %w", err)
	}
	return result, nil
}

func (repository *AdminRepository) GetEmployee(ctx context.Context, id int64) (AdminEmployeeRecord, error) {
	return scanAdminEmployee(repository.db.QueryRowContext(ctx, adminEmployeeSelect+` WHERE u.id=$1 GROUP BY u.id,eg.id,eg.title`, id))
}

const adminEmployeeSelect = `
	SELECT u.id,u.username,u.full_name,u.role,eg.id,eg.title,u.avg_rating,u.max_tickets,
	       COUNT(t.id) FILTER (WHERE tw.actual AND t.status IN (2,3,4,5))
	FROM users u JOIN expert_groups eg ON eg.id=u.expert_group_id
	LEFT JOIN tickets_workers tw ON tw.worker_id=u.id
	LEFT JOIN tickets t ON t.id=tw.ticket_id`

type rowScanner interface{ Scan(...any) error }

func scanAdminEmployee(row rowScanner) (AdminEmployeeRecord, error) {
	var item AdminEmployeeRecord
	err := row.Scan(&item.ID, &item.Username, &item.FullName, &item.Role, &item.GroupID, &item.GroupTitle, &item.AverageRating, &item.MaxTickets, &item.ActiveTickets)
	if errors.Is(err, sql.ErrNoRows) {
		return AdminEmployeeRecord{}, ErrAdminResourceNotFound
	}
	if err != nil {
		return AdminEmployeeRecord{}, fmt.Errorf("scan admin employee: %w", err)
	}
	return item, nil
}

func (repository *AdminRepository) ListEmployees(ctx context.Context, role, search string, limit, offset int) (AdminEmployeePageRecord, error) {
	where := ` WHERE ($1='' OR u.role=$1) AND ($2='' OR u.full_name ILIKE '%%'||$2||'%%' OR u.username ILIKE '%%'||$2||'%%')`
	var total int
	if err := repository.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users u`+where, role, search).Scan(&total); err != nil {
		return AdminEmployeePageRecord{}, fmt.Errorf("count admin employees: %w", err)
	}
	rows, err := repository.db.QueryContext(ctx, adminEmployeeSelect+where+` GROUP BY u.id,eg.id,eg.title ORDER BY u.full_name,u.id LIMIT $3 OFFSET $4`, role, search, limit, offset)
	if err != nil {
		return AdminEmployeePageRecord{}, fmt.Errorf("list admin employees: %w", err)
	}
	defer rows.Close()
	result := AdminEmployeePageRecord{Items: []AdminEmployeeRecord{}, Total: total}
	for rows.Next() {
		item, err := scanAdminEmployee(rows)
		if err != nil {
			return AdminEmployeePageRecord{}, err
		}
		result.Items = append(result.Items, item)
	}
	return result, rows.Err()
}

func (repository *AdminRepository) CreateEmployee(ctx context.Context, username, passwordHash, fullName, role string, groupID int64, maxTickets int) (AdminEmployeeRecord, error) {
	var id int64
	err := repository.db.QueryRowContext(ctx, `INSERT INTO users(username,password_hash,role,expert_group_id,full_name,max_tickets) VALUES($1,$2,$3,$4,$5,$6) RETURNING id`, username, passwordHash, role, groupID, fullName, maxTickets).Scan(&id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return AdminEmployeeRecord{}, ErrUsernameExists
		}
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			return AdminEmployeeRecord{}, ErrAdminResourceNotFound
		}
		return AdminEmployeeRecord{}, fmt.Errorf("create admin employee: %w", err)
	}
	return repository.GetEmployee(ctx, id)
}

func (repository *AdminRepository) UpdateEmployee(ctx context.Context, id int64, username, passwordHash, fullName, role string, groupID int64, maxTickets int) (AdminEmployeeRecord, error) {
	result, err := repository.db.ExecContext(ctx, `UPDATE users SET username=$2,password_hash=CASE WHEN $3='' THEN password_hash ELSE $3 END,full_name=$4,role=$5,expert_group_id=$6,max_tickets=$7 WHERE id=$1`, id, username, passwordHash, fullName, role, groupID, maxTickets)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return AdminEmployeeRecord{}, ErrUsernameExists
		}
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			return AdminEmployeeRecord{}, ErrAdminResourceNotFound
		}
		return AdminEmployeeRecord{}, fmt.Errorf("update admin employee: %w", err)
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return AdminEmployeeRecord{}, ErrAdminResourceNotFound
	}
	return repository.GetEmployee(ctx, id)
}

func (repository *AdminRepository) DeactivateEmployee(ctx context.Context, id int64) (AdminEmployeeRecord, error) {
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return AdminEmployeeRecord{}, err
	}
	defer tx.Rollback()
	var role string
	if err = tx.QueryRowContext(ctx, `SELECT role FROM users WHERE id=$1 FOR UPDATE`, id).Scan(&role); errors.Is(err, sql.ErrNoRows) {
		return AdminEmployeeRecord{}, ErrAdminResourceNotFound
	} else if err != nil {
		return AdminEmployeeRecord{}, err
	}
	if role == "expert" {
		var active int
		if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM tickets_workers tw JOIN tickets t ON t.id=tw.ticket_id WHERE tw.worker_id=$1 AND tw.actual AND t.status IN (2,3,4,5)`, id).Scan(&active); err != nil {
			return AdminEmployeeRecord{}, err
		}
		if active > 0 {
			return AdminEmployeeRecord{}, ErrEmployeeHasTickets
		}
	}
	if _, err = tx.ExecContext(ctx, `UPDATE users SET role='200' WHERE id=$1`, id); err != nil {
		return AdminEmployeeRecord{}, err
	}
	if err = tx.Commit(); err != nil {
		return AdminEmployeeRecord{}, err
	}
	return repository.GetEmployee(ctx, id)
}

func (repository *AdminRepository) ListGroups(ctx context.Context, search string, limit, offset int) (AdminGroupPageRecord, error) {
	where := ` WHERE ($1='' OR eg.title ILIKE '%%'||$1||'%%')`
	var total int
	if err := repository.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM expert_groups eg`+where, search).Scan(&total); err != nil {
		return AdminGroupPageRecord{}, err
	}
	rows, err := repository.db.QueryContext(ctx, `SELECT eg.id,eg.title,COUNT(DISTINCT u.id),COUNT(DISTINCT ceg.cat_id) FROM expert_groups eg LEFT JOIN users u ON u.expert_group_id=eg.id LEFT JOIN cats_expert_groups ceg ON ceg.group_id=eg.id`+where+` GROUP BY eg.id,eg.title ORDER BY eg.title,eg.id LIMIT $2 OFFSET $3`, search, limit, offset)
	if err != nil {
		return AdminGroupPageRecord{}, err
	}
	defer rows.Close()
	r := AdminGroupPageRecord{Items: []AdminGroupRecord{}, Total: total}
	for rows.Next() {
		var x AdminGroupRecord
		if err = rows.Scan(&x.ID, &x.Title, &x.EmployeeCount, &x.CategoryCount); err != nil {
			return AdminGroupPageRecord{}, err
		}
		r.Items = append(r.Items, x)
	}
	return r, rows.Err()
}

func (repository *AdminRepository) CreateGroup(ctx context.Context, title string) (AdminGroupRecord, error) {
	var id int64
	err := repository.db.QueryRowContext(ctx, `INSERT INTO expert_groups(title) SELECT $1 WHERE NOT EXISTS(SELECT 1 FROM expert_groups WHERE lower(title)=lower($1)) RETURNING id`, title).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return AdminGroupRecord{}, ErrNameExists
	}
	if err != nil {
		return AdminGroupRecord{}, err
	}
	return repository.getGroup(ctx, id)
}
func (repository *AdminRepository) getGroup(ctx context.Context, id int64) (AdminGroupRecord, error) {
	var x AdminGroupRecord
	err := repository.db.QueryRowContext(ctx, `SELECT eg.id,eg.title,COUNT(DISTINCT u.id),COUNT(DISTINCT ceg.cat_id) FROM expert_groups eg LEFT JOIN users u ON u.expert_group_id=eg.id LEFT JOIN cats_expert_groups ceg ON ceg.group_id=eg.id WHERE eg.id=$1 GROUP BY eg.id,eg.title`, id).Scan(&x.ID, &x.Title, &x.EmployeeCount, &x.CategoryCount)
	if errors.Is(err, sql.ErrNoRows) {
		return x, ErrAdminResourceNotFound
	}
	return x, err
}
func (repository *AdminRepository) UpdateGroup(ctx context.Context, id int64, title string) (AdminGroupRecord, error) {
	res, err := repository.db.ExecContext(ctx, `UPDATE expert_groups SET title=$2 WHERE id=$1 AND NOT EXISTS(SELECT 1 FROM expert_groups WHERE lower(title)=lower($2) AND id<>$1)`, id, title)
	if err != nil {
		return AdminGroupRecord{}, err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		var exists bool
		_ = repository.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM expert_groups WHERE id=$1)`, id).Scan(&exists)
		if exists {
			return AdminGroupRecord{}, ErrNameExists
		}
		return AdminGroupRecord{}, ErrAdminResourceNotFound
	}
	return repository.getGroup(ctx, id)
}
func (repository *AdminRepository) DeleteGroup(ctx context.Context, id int64) error {
	res, err := repository.db.ExecContext(ctx, `DELETE FROM expert_groups eg WHERE id=$1 AND NOT EXISTS(SELECT 1 FROM users WHERE expert_group_id=eg.id) AND NOT EXISTS(SELECT 1 FROM cats_expert_groups WHERE group_id=eg.id)`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		var exists bool
		_ = repository.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM expert_groups WHERE id=$1)`, id).Scan(&exists)
		if exists {
			return ErrAdminResourceInUse
		}
		return ErrAdminResourceNotFound
	}
	return nil
}

func (repository *AdminRepository) ListCategories(ctx context.Context, search string, limit, offset int) (AdminCategoryPageRecord, error) {
	where := ` WHERE ($1='' OR c.name ILIKE '%%'||$1||'%%')`
	var total int
	if err := repository.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM categories c`+where, search).Scan(&total); err != nil {
		return AdminCategoryPageRecord{}, err
	}
	rows, err := repository.db.QueryContext(ctx, `SELECT c.id,c.name FROM categories c`+where+` ORDER BY c.name,c.id LIMIT $2 OFFSET $3`, search, limit, offset)
	if err != nil {
		return AdminCategoryPageRecord{}, err
	}
	defer rows.Close()
	r := AdminCategoryPageRecord{Items: []AdminCategoryRecord{}, Total: total}
	for rows.Next() {
		var x AdminCategoryRecord
		if err = rows.Scan(&x.ID, &x.Name); err != nil {
			return r, err
		}
		r.Items = append(r.Items, x)
	}
	return r, rows.Err()
}
func (repository *AdminRepository) GetCategory(ctx context.Context, id int64) (AdminCategoryRecord, error) {
	var r AdminCategoryRecord
	if err := repository.db.QueryRowContext(ctx, `SELECT id,name FROM categories WHERE id=$1`, id).Scan(&r.ID, &r.Name); errors.Is(err, sql.ErrNoRows) {
		return r, ErrAdminResourceNotFound
	} else if err != nil {
		return r, err
	}
	r.Groups = []AdminGroupRecord{}
	r.Questions = []AdminQuestionRecord{}
	rows, err := repository.db.QueryContext(ctx, `SELECT eg.id,eg.title FROM expert_groups eg JOIN cats_expert_groups ceg ON ceg.group_id=eg.id WHERE ceg.cat_id=$1 ORDER BY eg.title,eg.id`, id)
	if err != nil {
		return r, err
	}
	for rows.Next() {
		var g AdminGroupRecord
		if err = rows.Scan(&g.ID, &g.Title); err != nil {
			rows.Close()
			return r, err
		}
		r.Groups = append(r.Groups, g)
	}
	if err = rows.Close(); err != nil {
		return r, err
	}
	qrows, err := repository.db.QueryContext(ctx, `SELECT q.id,q.text,a.id,a.text FROM questions q LEFT JOIN answers a ON a.question_id=q.id WHERE q.cat_id=$1 ORDER BY q.id,a.id`, id)
	if err != nil {
		return r, err
	}
	defer qrows.Close()
	index := map[int64]int{}
	for qrows.Next() {
		var qid int64
		var qt string
		var aid sql.NullInt64
		var at sql.NullString
		if err = qrows.Scan(&qid, &qt, &aid, &at); err != nil {
			return r, err
		}
		pos, ok := index[qid]
		if !ok {
			pos = len(r.Questions)
			index[qid] = pos
			r.Questions = append(r.Questions, AdminQuestionRecord{ID: qid, Text: qt, Answers: []AdminAnswerRecord{}})
		}
		if aid.Valid {
			r.Questions[pos].Answers = append(r.Questions[pos].Answers, AdminAnswerRecord{ID: aid.Int64, Text: at.String})
		}
	}
	return r, qrows.Err()
}
func (repository *AdminRepository) CreateCategory(ctx context.Context, name string) (AdminCategoryRecord, error) {
	var id int64
	err := repository.db.QueryRowContext(ctx, `INSERT INTO categories(name) SELECT $1 WHERE NOT EXISTS(SELECT 1 FROM categories WHERE lower(name)=lower($1)) RETURNING id`, name).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return AdminCategoryRecord{}, ErrNameExists
	}
	if err != nil {
		return AdminCategoryRecord{}, err
	}
	return repository.GetCategory(ctx, id)
}
func (repository *AdminRepository) UpdateCategory(ctx context.Context, id int64, name string) (AdminCategoryRecord, error) {
	res, err := repository.db.ExecContext(ctx, `UPDATE categories SET name=$2 WHERE id=$1 AND NOT EXISTS(SELECT 1 FROM categories WHERE lower(name)=lower($2) AND id<>$1)`, id, name)
	if err != nil {
		return AdminCategoryRecord{}, err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		var exists bool
		_ = repository.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM categories WHERE id=$1)`, id).Scan(&exists)
		if exists {
			return AdminCategoryRecord{}, ErrNameExists
		}
		return AdminCategoryRecord{}, ErrAdminResourceNotFound
	}
	return repository.GetCategory(ctx, id)
}
func (repository *AdminRepository) DeleteCategory(ctx context.Context, id int64) error {
	var name string
	if err := repository.db.QueryRowContext(ctx, `SELECT name FROM categories WHERE id=$1`, id).Scan(&name); errors.Is(err, sql.ErrNoRows) {
		return ErrAdminResourceNotFound
	} else if err != nil {
		return err
	}
	if strings.EqualFold(strings.TrimSpace(name), "Не знаю, как это назвать") {
		return ErrAdminResourceInUse
	}
	res, err := repository.db.ExecContext(ctx, `DELETE FROM categories c WHERE id=$1 AND NOT EXISTS(SELECT 1 FROM tickets WHERE category_id=c.id) AND NOT EXISTS(SELECT 1 FROM questions WHERE cat_id=c.id) AND NOT EXISTS(SELECT 1 FROM cats_expert_groups WHERE cat_id=c.id)`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrAdminResourceInUse
	}
	return nil
}

func (repository *AdminRepository) ReplaceCategoryGroups(ctx context.Context, categoryID int64, groupIDs []int64) (AdminCategoryRecord, error) {
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return AdminCategoryRecord{}, err
	}
	defer tx.Rollback()
	var exists bool
	if err = tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM categories WHERE id=$1)`, categoryID).Scan(&exists); err != nil {
		return AdminCategoryRecord{}, err
	}
	if !exists {
		return AdminCategoryRecord{}, ErrAdminResourceNotFound
	}
	if len(groupIDs) > 0 {
		var count int
		if err = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM expert_groups WHERE id=ANY($1)`, pq.Array(groupIDs)).Scan(&count); err != nil {
			return AdminCategoryRecord{}, err
		}
		if count != len(groupIDs) {
			return AdminCategoryRecord{}, ErrAdminResourceNotFound
		}
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM cats_expert_groups WHERE cat_id=$1`, categoryID); err != nil {
		return AdminCategoryRecord{}, err
	}
	for _, id := range groupIDs {
		if _, err = tx.ExecContext(ctx, `INSERT INTO cats_expert_groups(group_id,cat_id) VALUES($1,$2)`, id, categoryID); err != nil {
			return AdminCategoryRecord{}, err
		}
	}
	if err = tx.Commit(); err != nil {
		return AdminCategoryRecord{}, err
	}
	return repository.GetCategory(ctx, categoryID)
}

func (repository *AdminRepository) CreateQuestion(ctx context.Context, categoryID int64, text string, answers []string) (AdminQuestionRecord, error) {
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return AdminQuestionRecord{}, err
	}
	defer tx.Rollback()
	var qid int64
	if err = tx.QueryRowContext(ctx, `INSERT INTO questions(cat_id,text) VALUES($1,$2) RETURNING id`, categoryID, text).Scan(&qid); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			return AdminQuestionRecord{}, ErrAdminResourceNotFound
		}
		return AdminQuestionRecord{}, err
	}
	r := AdminQuestionRecord{ID: qid, Text: text, Answers: []AdminAnswerRecord{}}
	for _, answer := range answers {
		var id int64
		if err = tx.QueryRowContext(ctx, `INSERT INTO answers(question_id,text) VALUES($1,$2) RETURNING id`, qid, answer).Scan(&id); err != nil {
			return AdminQuestionRecord{}, err
		}
		r.Answers = append(r.Answers, AdminAnswerRecord{ID: id, Text: answer})
	}
	if err = tx.Commit(); err != nil {
		return AdminQuestionRecord{}, err
	}
	return r, nil
}
func (repository *AdminRepository) UpdateQuestion(ctx context.Context, id int64, text string, answers []string) (AdminQuestionRecord, error) {
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return AdminQuestionRecord{}, err
	}
	defer tx.Rollback()
	var exists, used bool
	if err = tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM questions WHERE id=$1),EXISTS(SELECT 1 FROM "QA" WHERE question_id=$1)`, id).Scan(&exists, &used); err != nil {
		return AdminQuestionRecord{}, err
	}
	if !exists {
		return AdminQuestionRecord{}, ErrAdminResourceNotFound
	}
	if used {
		return AdminQuestionRecord{}, ErrQuestionInUse
	}
	if _, err = tx.ExecContext(ctx, `UPDATE questions SET text=$2 WHERE id=$1`, id, text); err != nil {
		return AdminQuestionRecord{}, err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM answers WHERE question_id=$1`, id); err != nil {
		return AdminQuestionRecord{}, err
	}
	r := AdminQuestionRecord{ID: id, Text: text, Answers: []AdminAnswerRecord{}}
	for _, answer := range answers {
		var aid int64
		if err = tx.QueryRowContext(ctx, `INSERT INTO answers(question_id,text) VALUES($1,$2) RETURNING id`, id, answer).Scan(&aid); err != nil {
			return AdminQuestionRecord{}, err
		}
		r.Answers = append(r.Answers, AdminAnswerRecord{ID: aid, Text: answer})
	}
	if err = tx.Commit(); err != nil {
		return AdminQuestionRecord{}, err
	}
	return r, nil
}
func (repository *AdminRepository) DeleteQuestion(ctx context.Context, id int64) error {
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var exists, used bool
	if err = tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM questions WHERE id=$1),EXISTS(SELECT 1 FROM "QA" WHERE question_id=$1)`, id).Scan(&exists, &used); err != nil {
		return err
	}
	if !exists {
		return ErrAdminResourceNotFound
	}
	if used {
		return ErrQuestionInUse
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM answers WHERE question_id=$1`, id); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM questions WHERE id=$1`, id); err != nil {
		return err
	}
	return tx.Commit()
}

func adminRoutingIssueSQL(alias string) string {
	return `CASE WHEN NOT EXISTS(SELECT 1 FROM cats_expert_groups ceg WHERE ceg.cat_id=` + alias + `.category_id) THEN 'no_group' WHEN NOT EXISTS(SELECT 1 FROM users u WHERE u.role='expert' AND u.expert_group_id IN(SELECT ceg.group_id FROM cats_expert_groups ceg WHERE ceg.cat_id=` + alias + `.category_id) AND (SELECT COUNT(*) FROM tickets_workers tw JOIN tickets at ON at.id=tw.ticket_id WHERE tw.worker_id=u.id AND tw.actual AND at.status IN (2,3,4,5)) < u.max_tickets) THEN 'no_available_expert' END`
}
func adminTicketWhere(f AdminTicketFilter) (string, []any) {
	c := []string{}
	a := []any{}
	add := func(format string, v any) { a = append(a, v); c = append(c, fmt.Sprintf(format, len(a))) }
	if f.Search != "" {
		add("t.track_id ILIKE '%%'||$%d||'%%'", f.Search)
	}
	if f.Status != 0 {
		add("t.status=$%d", f.Status)
	}
	if f.Priority != 0 {
		add("t.priority=$%d", f.Priority)
	}
	if f.CategoryID != 0 {
		add("t.category_id=$%d", f.CategoryID)
	}
	if f.ApplicantType != 0 {
		add("t.morda_type=$%d", f.ApplicantType)
	}
	if f.ResponsibleID != 0 {
		add("EXISTS(SELECT 1 FROM tickets_workers tw WHERE tw.ticket_id=t.id AND tw.worker_id=$%d AND tw.actual AND tw.is_responsible)", f.ResponsibleID)
	}
	if f.RoutingIssue != "" {
		add("("+adminRoutingIssueSQL("t")+")=$%d", f.RoutingIssue)
	}
	if len(c) == 0 {
		return "", a
	}
	return " WHERE " + strings.Join(c, " AND "), a
}

func (repository *AdminRepository) ListTickets(ctx context.Context, f AdminTicketFilter) (AdminTicketPageRecord, error) {
	where, args := adminTicketWhere(f)
	var total int
	if err := repository.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM tickets t`+where, args...).Scan(&total); err != nil {
		return AdminTicketPageRecord{}, err
	}
	args = append(args, f.Limit, f.Offset)
	rows, err := repository.db.QueryContext(ctx, `SELECT t.track_id,c.id,c.name,t.status,t.morda_type,t.priority,t.created_at,t.closed_at,t.return_count,r.worker_id,r.full_name,`+adminRoutingIssueSQL("t")+` FROM tickets t JOIN categories c ON c.id=t.category_id LEFT JOIN LATERAL(SELECT u.id worker_id,u.full_name FROM tickets_workers tw JOIN users u ON u.id=tw.worker_id WHERE tw.ticket_id=t.id AND tw.actual AND tw.is_responsible LIMIT 1)r ON TRUE`+where+fmt.Sprintf(` ORDER BY CASE t.priority WHEN 2 THEN 0 WHEN 1 THEN 1 ELSE 2 END,t.created_at,t.id LIMIT $%d OFFSET $%d`, len(args)-1, len(args)), args...)
	if err != nil {
		return AdminTicketPageRecord{}, err
	}
	defer rows.Close()
	r := AdminTicketPageRecord{Items: []AdminTicketRecord{}, Total: total}
	for rows.Next() {
		var x AdminTicketRecord
		if err = rows.Scan(&x.TrackID, &x.CategoryID, &x.Category, &x.Status, &x.ApplicantType, &x.Priority, &x.CreatedAt, &x.ClosedAt, &x.ReturnCount, &x.ResponsibleID, &x.ResponsibleName, &x.RoutingIssue); err != nil {
			return r, err
		}
		r.Items = append(r.Items, x)
	}
	return r, rows.Err()
}

func (repository *AdminRepository) GetTicket(ctx context.Context, trackID string) (AdminTicketRecord, error) {
	var r AdminTicketRecord
	err := repository.db.QueryRowContext(ctx, `SELECT t.track_id,c.id,c.name,t.status,t.morda_type,t.priority,t.created_at,t.closed_at,t.return_count,r.worker_id,r.full_name,`+adminRoutingIssueSQL("t")+` FROM tickets t JOIN categories c ON c.id=t.category_id LEFT JOIN LATERAL(SELECT u.id worker_id,u.full_name FROM tickets_workers tw JOIN users u ON u.id=tw.worker_id WHERE tw.ticket_id=t.id AND tw.actual AND tw.is_responsible LIMIT 1)r ON TRUE WHERE t.track_id=$1`, trackID).Scan(&r.TrackID, &r.CategoryID, &r.Category, &r.Status, &r.ApplicantType, &r.Priority, &r.CreatedAt, &r.ClosedAt, &r.ReturnCount, &r.ResponsibleID, &r.ResponsibleName, &r.RoutingIssue)
	if errors.Is(err, sql.ErrNoRows) {
		return r, ErrTicketNotFound
	}
	if err != nil {
		return r, err
	}
	rows, err := repository.db.QueryContext(ctx, `SELECT u.id,u.full_name,eg.id,eg.title,tw.is_responsible,tw.actual,tw.created_at FROM tickets_workers tw JOIN tickets t ON t.id=tw.ticket_id JOIN users u ON u.id=tw.worker_id JOIN expert_groups eg ON eg.id=u.expert_group_id WHERE t.track_id=$1 AND tw.actual ORDER BY tw.is_responsible DESC,tw.created_at,u.id`, trackID)
	if err != nil {
		return r, err
	}
	defer rows.Close()
	r.Workers = []AdminWorkerRecord{}
	for rows.Next() {
		var x AdminWorkerRecord
		if err = rows.Scan(&x.ID, &x.FullName, &x.GroupID, &x.GroupTitle, &x.IsResponsible, &x.Actual, &x.CreatedAt); err != nil {
			return r, err
		}
		r.Workers = append(r.Workers, x)
	}
	return r, rows.Err()
}

func (repository *AdminRepository) UpdateTicketPriority(ctx context.Context, trackID string, priority models.TicketPriority, actorID int64, changedAt time.Time) (AdminTicketRecord, error) {
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return AdminTicketRecord{}, err
	}
	defer tx.Rollback()
	var id int64
	var old models.TicketPriority
	var status models.TicketStatus
	if err = tx.QueryRowContext(ctx, `SELECT id,priority,status FROM tickets WHERE track_id=$1 FOR UPDATE`, trackID).Scan(&id, &old, &status); errors.Is(err, sql.ErrNoRows) {
		return AdminTicketRecord{}, ErrTicketNotFound
	} else if err != nil {
		return AdminTicketRecord{}, err
	}
	if old != priority {
		if _, err = tx.ExecContext(ctx, `UPDATE tickets SET priority=$2 WHERE id=$1`, id, priority); err != nil {
			return AdminTicketRecord{}, err
		}
		if err = insertEvent(ctx, tx, id, actorID, "admin_priority_changed", status, status, sql.NullInt64{}, sql.NullInt64{}, changedAt); err != nil {
			return AdminTicketRecord{}, err
		}
	}
	if err = tx.Commit(); err != nil {
		return AdminTicketRecord{}, err
	}
	return repository.GetTicket(ctx, trackID)
}

func (repository *AdminRepository) UpdateTicketStatus(ctx context.Context, trackID string, status models.TicketStatus, actorID int64, changedAt time.Time) (AdminTicketRecord, error) {
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return AdminTicketRecord{}, err
	}
	defer tx.Rollback()
	var id int64
	var old models.TicketStatus
	if err = tx.QueryRowContext(ctx, `SELECT id,status FROM tickets WHERE track_id=$1 FOR UPDATE`, trackID).Scan(&id, &old); errors.Is(err, sql.ErrNoRows) {
		return AdminTicketRecord{}, ErrTicketNotFound
	} else if err != nil {
		return AdminTicketRecord{}, err
	}
	if old == status {
		if err = tx.Commit(); err != nil {
			return AdminTicketRecord{}, err
		}
		return repository.GetTicket(ctx, trackID)
	}
	if status >= models.TicketStatusAssigned && status <= models.TicketStatusAnswerReady {
		var responsible bool
		if err = tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM tickets_workers WHERE ticket_id=$1 AND actual AND is_responsible)`, id).Scan(&responsible); err != nil {
			return AdminTicketRecord{}, err
		}
		if !responsible {
			return AdminTicketRecord{}, ErrResponsibleRequired
		}
	}
	terminal := status == models.TicketStatusCompleted || status == models.TicketStatusRejected || status == models.TicketStatusClosedWithoutAnswer
	if terminal || status == models.TicketStatusNew || status == models.TicketStatusReturned {
		if _, err = tx.ExecContext(ctx, `UPDATE tickets_workers SET actual=FALSE WHERE ticket_id=$1 AND actual`, id); err != nil {
			return AdminTicketRecord{}, err
		}
	}
	if terminal {
		_, err = tx.ExecContext(ctx, `UPDATE tickets SET status=$2,closed_at=$3 WHERE id=$1`, id, status, changedAt)
	} else {
		_, err = tx.ExecContext(ctx, `UPDATE tickets SET status=$2,closed_at=NULL WHERE id=$1`, id, status)
	}
	if err != nil {
		return AdminTicketRecord{}, err
	}
	if err = insertEvent(ctx, tx, id, actorID, "admin_status_changed", old, status, sql.NullInt64{}, sql.NullInt64{}, changedAt); err != nil {
		return AdminTicketRecord{}, err
	}
	if err = tx.Commit(); err != nil {
		return AdminTicketRecord{}, err
	}
	return repository.GetTicket(ctx, trackID)
}

func (repository *AdminRepository) SetResponsibleWorker(ctx context.Context, trackID string, workerID, actorID int64, changedAt time.Time) (AdminTicketRecord, error) {
	tx, err := repository.db.BeginTx(ctx, nil)
	if err != nil {
		return AdminTicketRecord{}, err
	}
	defer tx.Rollback()
	var ticketID int64
	var status models.TicketStatus
	if err = tx.QueryRowContext(ctx, `SELECT id,status FROM tickets WHERE track_id=$1 FOR UPDATE`, trackID).Scan(&ticketID, &status); errors.Is(err, sql.ErrNoRows) {
		return AdminTicketRecord{}, ErrTicketNotFound
	} else if err != nil {
		return AdminTicketRecord{}, err
	}
	if err = lockAvailableExpert(ctx, tx, workerID, ticketID); err != nil {
		return AdminTicketRecord{}, err
	}
	var previous sql.NullInt64
	err = tx.QueryRowContext(ctx, `SELECT worker_id FROM tickets_workers WHERE ticket_id=$1 AND actual AND is_responsible FOR UPDATE`, ticketID).Scan(&previous)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return AdminTicketRecord{}, err
	}
	if previous.Valid && previous.Int64 != workerID {
		if _, err = tx.ExecContext(ctx, `UPDATE tickets_workers SET actual=FALSE WHERE ticket_id=$1 AND actual AND is_responsible`, ticketID); err != nil {
			return AdminTicketRecord{}, err
		}
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO tickets_workers(worker_id,ticket_id,is_responsible,created_at,actual) VALUES($1,$2,TRUE,$3,TRUE) ON CONFLICT(worker_id,ticket_id) DO UPDATE SET is_responsible=TRUE,actual=TRUE,created_at=EXCLUDED.created_at`, workerID, ticketID, changedAt); err != nil {
		return AdminTicketRecord{}, err
	}
	if _, err = tx.ExecContext(ctx, `UPDATE tickets SET status=$2,closed_at=NULL WHERE id=$1`, ticketID, models.TicketStatusAssigned); err != nil {
		return AdminTicketRecord{}, err
	}
	if err = insertEvent(ctx, tx, ticketID, actorID, "admin_responsible_changed", status, models.TicketStatusAssigned, previous, sql.NullInt64{Int64: workerID, Valid: true}, changedAt); err != nil {
		return AdminTicketRecord{}, err
	}
	if err = tx.Commit(); err != nil {
		return AdminTicketRecord{}, err
	}
	return repository.GetTicket(ctx, trackID)
}

func (repository *AdminRepository) ListEvents(ctx context.Context, trackID string, limit, offset int) (AdminEventPageRecord, error) {
	var ticketID int64
	if err := repository.db.QueryRowContext(ctx, `SELECT id FROM tickets WHERE track_id=$1`, trackID).Scan(&ticketID); errors.Is(err, sql.ErrNoRows) {
		return AdminEventPageRecord{}, ErrTicketNotFound
	} else if err != nil {
		return AdminEventPageRecord{}, err
	}
	var total int
	if err := repository.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM ticket_events WHERE ticket_id=$1`, ticketID).Scan(&total); err != nil {
		return AdminEventPageRecord{}, err
	}
	rows, err := repository.db.QueryContext(ctx, `SELECT e.id,e.event_type,e.actor_user_id,a.full_name,a.role,e.from_status,e.to_status,e.from_worker_id,fw.full_name,e.to_worker_id,tw.full_name,e.reason_code,e.created_at FROM ticket_events e LEFT JOIN users a ON a.id=e.actor_user_id LEFT JOIN users fw ON fw.id=e.from_worker_id LEFT JOIN users tw ON tw.id=e.to_worker_id WHERE e.ticket_id=$1 ORDER BY e.created_at DESC,e.id DESC LIMIT $2 OFFSET $3`, ticketID, limit, offset)
	if err != nil {
		return AdminEventPageRecord{}, err
	}
	defer rows.Close()
	r := AdminEventPageRecord{Items: []AdminEventRecord{}, Total: total}
	for rows.Next() {
		var x AdminEventRecord
		if err = rows.Scan(&x.ID, &x.EventType, &x.ActorID, &x.ActorName, &x.ActorRole, &x.FromStatus, &x.ToStatus, &x.FromWorkerID, &x.FromWorkerName, &x.ToWorkerID, &x.ToWorkerName, &x.ReasonCode, &x.CreatedAt); err != nil {
			return r, err
		}
		r.Items = append(r.Items, x)
	}
	return r, rows.Err()
}

func (repository *AdminRepository) Analytics(ctx context.Context, start, end time.Time) (AnalyticsRecord, error) {
	return NewOperatorRepository(repository.db).Analytics(ctx, start, end)
}
func (repository *AdminRepository) ReportTickets(ctx context.Context, start, end time.Time) ([]ReportTicketRecord, error) {
	return NewOperatorRepository(repository.db).ReportTickets(ctx, start, end)
}
func (repository *AdminRepository) EmployeeLoads(ctx context.Context, start, end time.Time) ([]AdminEmployeeLoadRecord, error) {
	rows, err := repository.db.QueryContext(ctx, `SELECT u.id,u.full_name,u.role,COUNT(DISTINCT t.id) FILTER(WHERE tw.actual AND t.status IN (2,3,4,5)),u.max_tickets,COUNT(DISTINCT e.id) FILTER(WHERE e.created_at >= $1 AND e.created_at < $2) FROM users u LEFT JOIN tickets_workers tw ON tw.worker_id=u.id LEFT JOIN tickets t ON t.id=tw.ticket_id LEFT JOIN ticket_events e ON e.actor_user_id=u.id WHERE u.role IN('operator','expert') GROUP BY u.id,u.full_name,u.role,u.max_tickets ORDER BY u.full_name,u.id`, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	r := []AdminEmployeeLoadRecord{}
	for rows.Next() {
		var x AdminEmployeeLoadRecord
		if err = rows.Scan(&x.ID, &x.FullName, &x.Role, &x.ActiveTickets, &x.MaxTickets, &x.ActionCount); err != nil {
			return nil, err
		}
		r = append(r, x)
	}
	return r, rows.Err()
}

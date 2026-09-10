package repos

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"
	"time"

	"example.com/german/backend/internal/models"
	_ "github.com/lib/pq"
)

func TestOperatorWorkflowIntegration(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("sql.Open() error = %v", err)
	}
	defer db.Close()
	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		t.Fatalf("database ping error = %v", err)
	}

	var groupID, categoryID, operatorID, firstExpertID, coExpertID, busyExpertID int64
	if err := db.QueryRowContext(ctx, `INSERT INTO expert_groups (title) VALUES ('Психологи') RETURNING id`).Scan(&groupID); err != nil {
		t.Fatalf("insert expert group: %v", err)
	}
	if err := db.QueryRowContext(ctx, `INSERT INTO categories (name) VALUES ('Кибербуллинг') RETURNING id`).Scan(&categoryID); err != nil {
		t.Fatalf("insert category: %v", err)
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO cats_expert_groups (group_id, cat_id) VALUES ($1, $2)`, groupID, categoryID); err != nil {
		t.Fatalf("insert route: %v", err)
	}
	insertUser := func(username, role, fullName string, maxTickets int) int64 {
		t.Helper()
		var id int64
		err := db.QueryRowContext(ctx, `
			INSERT INTO users (username, password_hash, role, expert_group_id, full_name, max_tickets)
			VALUES ($1, 'unused', $2, $3, $4, $5) RETURNING id`,
			username, role, groupID, fullName, maxTickets).Scan(&id)
		if err != nil {
			t.Fatalf("insert user %s: %v", username, err)
		}
		return id
	}
	operatorID = insertUser("operator", "operator", "Оператор", 10)
	firstExpertID = insertUser("expert-one", "expert", "Первый эксперт", 2)
	coExpertID = insertUser("expert-two", "expert", "Второй эксперт", 2)
	busyExpertID = insertUser("expert-busy", "expert", "Занятый эксперт", 1)

	now := time.Now().UTC().Truncate(time.Second)
	insertTicket := func(track string, status models.TicketStatus, createdAt time.Time) int64 {
		t.Helper()
		var id int64
		err := db.QueryRowContext(ctx, `
			INSERT INTO tickets (track_id, priority, morda_type, category_id, status, created_at)
			VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
			track, models.TicketPriorityUrgent, models.ApplicantTypeSchoolchild,
			categoryID, status, createdAt).Scan(&id)
		if err != nil {
			t.Fatalf("insert ticket %s: %v", track, err)
		}
		return id
	}
	busyTicketID := insertTicket("ОТК-BUSY-2345", models.TicketStatusAssigned, now.Add(-2*time.Hour))
	if _, err := db.ExecContext(ctx, `
		INSERT INTO tickets_workers (worker_id, ticket_id, is_responsible, actual, created_at)
		VALUES ($1, $2, TRUE, TRUE, $3)`, busyExpertID, busyTicketID, now.Add(-90*time.Minute)); err != nil {
		t.Fatalf("insert busy assignment: %v", err)
	}
	ticketID := insertTicket("ОТК-ABCD-2345", models.TicketStatusNew, now.Add(-time.Hour))
	_ = ticketID

	repository := NewOperatorRepository(db)
	eligible, err := repository.EligibleWorkers(ctx, "ОТК-ABCD-2345", "", 0)
	if err != nil {
		t.Fatalf("EligibleWorkers() error = %v", err)
	}
	if len(eligible.Workers) != 2 || eligible.RecommendedGroup == nil {
		t.Fatalf("eligible workers = %+v", eligible)
	}
	assignment, err := repository.SetResponsibleWorker(ctx, "ОТК-ABCD-2345", firstExpertID, operatorID, now)
	if err != nil || assignment.Status != models.TicketStatusAssigned {
		t.Fatalf("SetResponsibleWorker() = (%+v, %v)", assignment, err)
	}
	if _, err := repository.AddWorker(ctx, "ОТК-ABCD-2345", busyExpertID, operatorID, now); !errors.Is(err, ErrWorkerUnavailable) {
		t.Fatalf("AddWorker(busy) error = %v, want %v", err, ErrWorkerUnavailable)
	}
	if _, err := repository.AddWorker(ctx, "ОТК-ABCD-2345", coExpertID, operatorID, now); err != nil {
		t.Fatalf("AddWorker() error = %v", err)
	}
	if _, err := repository.RemoveWorker(ctx, "ОТК-ABCD-2345", coExpertID, operatorID, now.Add(time.Minute)); err != nil {
		t.Fatalf("RemoveWorker() error = %v", err)
	}

	page, err := repository.ListTickets(ctx, OperatorTicketFilter{Queue: "assigned", Limit: 20})
	if err != nil || page.Total != 2 || len(page.Items) != 2 {
		t.Fatalf("ListTickets(assigned) = (%+v, %v)", page, err)
	}
	returnedTicketID := insertTicket("ОТК-RETN-2345", models.TicketStatusAnswerReady, now.Add(-3*time.Hour))
	if _, err := db.ExecContext(ctx, `
		INSERT INTO tickets_workers (worker_id, ticket_id, is_responsible, actual, created_at)
		VALUES ($1, $2, TRUE, TRUE, $3)`, firstExpertID, returnedTicketID, now.Add(-2*time.Hour)); err != nil {
		t.Fatalf("insert returned ticket assignment: %v", err)
	}
	ticketRepository := NewTicketRepository(db)
	if _, err := ticketRepository.Return(ctx, "ОТК-RETN-2345", "Нужны другие рекомендации", now); err != nil {
		t.Fatalf("Return() error = %v", err)
	}
	returnedPage, err := repository.ListTickets(ctx, OperatorTicketFilter{Queue: "returned", Limit: 20})
	if err != nil || returnedPage.Total != 1 || returnedPage.Items[0].ReturnReason.String != "Нужны другие рекомендации" || !returnedPage.Items[0].PreviousWorkerID.Valid {
		t.Fatalf("ListTickets(returned) = (%+v, %v)", returnedPage, err)
	}
	analytics, err := repository.Analytics(ctx, now.Add(-24*time.Hour), now.Add(24*time.Hour))
	if err != nil || analytics.Total != 3 || analytics.UrgentCount != 3 || analytics.ReturnedCount != 1 {
		t.Fatalf("Analytics() = (%+v, %v)", analytics, err)
	}
	report, err := repository.ReportTickets(ctx, now.Add(-24*time.Hour), now.Add(24*time.Hour))
	if err != nil || len(report) != 3 {
		t.Fatalf("ReportTickets() rows = %d, error = %v", len(report), err)
	}

	rejected, err := repository.Reject(ctx, "ОТК-ABCD-2345", operatorID, "spam", "Обращение отклонено", now.Add(time.Minute))
	if err != nil || rejected.Status != models.TicketStatusRejected {
		t.Fatalf("Reject() = (%+v, %v)", rejected, err)
	}
}

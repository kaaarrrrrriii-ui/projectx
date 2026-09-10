package repos

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"example.com/german/backend/internal/models"
	_ "github.com/lib/pq"
)

func TestExpertWorkflowIntegration(t *testing.T) {
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

	suffix := time.Now().UnixNano()
	var groupID, categoryID, expertID, coworkerID int64
	if err := db.QueryRowContext(ctx, `INSERT INTO expert_groups (title) VALUES ($1) RETURNING id`, fmt.Sprintf("Эксперты %d", suffix)).Scan(&groupID); err != nil {
		t.Fatalf("insert expert group: %v", err)
	}
	if err := db.QueryRowContext(ctx, `INSERT INTO categories (name) VALUES ($1) RETURNING id`, fmt.Sprintf("Категория %d", suffix)).Scan(&categoryID); err != nil {
		t.Fatalf("insert category: %v", err)
	}
	insertUser := func(label string) int64 {
		t.Helper()
		var id int64
		if err := db.QueryRowContext(ctx, `
			INSERT INTO users (username, password_hash, role, expert_group_id, full_name, max_tickets)
			VALUES ($1, 'unused', 'expert', $2, $3, 10) RETURNING id`,
			fmt.Sprintf("%s-%d", label, suffix), groupID, label).Scan(&id); err != nil {
			t.Fatalf("insert %s: %v", label, err)
		}
		return id
	}
	expertID = insertUser("expert")
	coworkerID = insertUser("coworker")
	trackID := fmt.Sprintf("EXPERT-%d", suffix)
	createdAt := time.Now().UTC().Add(-time.Hour).Truncate(time.Second)
	var ticketID int64
	if err := db.QueryRowContext(ctx, `
		INSERT INTO tickets (track_id, priority, morda_type, category_id, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		trackID, models.TicketPriorityUrgent, models.ApplicantTypeSchoolchild,
		categoryID, models.TicketStatusAssigned, createdAt).Scan(&ticketID); err != nil {
		t.Fatalf("insert ticket: %v", err)
	}
	if _, err := db.ExecContext(ctx, `
		INSERT INTO tickets_workers (worker_id, ticket_id, is_responsible, actual, created_at)
		VALUES ($1, $2, TRUE, TRUE, $3)`, expertID, ticketID, createdAt.Add(time.Minute)); err != nil {
		t.Fatalf("assign expert: %v", err)
	}
	if _, err := db.ExecContext(ctx, `INSERT INTO messages (ticket_id, text, type, created_at) VALUES ($1, 'Описание', $2, $3)`, ticketID, models.MessageTypeApplicant, createdAt); err != nil {
		t.Fatalf("insert description: %v", err)
	}

	repository := NewExpertRepository(db)
	page, err := repository.ListTickets(ctx, ExpertTicketFilter{WorkerID: expertID, Queue: "queue", Limit: 20})
	if err != nil || page.Total != 1 || len(page.Items) != 1 {
		t.Fatalf("ListTickets(queue) = (%+v, %v)", page, err)
	}
	status, err := repository.OpenTicket(ctx, trackID, expertID, createdAt.Add(2*time.Minute))
	if err != nil || status != models.TicketStatusInProgress {
		t.Fatalf("OpenTicket() = (%v, %v)", status, err)
	}
	if _, err := repository.AddNote(ctx, trackID, expertID, "Эксперт", "Заметка", createdAt.Add(3*time.Minute)); err != nil {
		t.Fatalf("AddNote() error = %v", err)
	}
	detail, err := repository.GetTicket(ctx, trackID, expertID)
	if err != nil || detail.Description != "Описание" || len(detail.Notes) != 1 {
		t.Fatalf("GetTicket() = (%+v, %v)", detail, err)
	}
	workerRequest, err := repository.CreateWorkerRequest(ctx, trackID, expertID, "add_coworker", "Нужен коллега", createdAt.Add(4*time.Minute))
	if err != nil {
		t.Fatalf("CreateWorkerRequest() error = %v", err)
	}
	requestPage, err := repository.ListWorkerRequests(ctx, expertID, "sent", 20, 0)
	if err != nil || requestPage.Total != 1 || requestPage.Items[0].ID != workerRequest.ID {
		t.Fatalf("ListWorkerRequests(sent) = (%+v, %v)", requestPage, err)
	}
	operatorRepository := NewOperatorRepository(db)
	if _, err := operatorRepository.CompleteWorkerRequest(ctx, workerRequest.ID, coworkerID, expertID, createdAt.Add(5*time.Minute)); err != nil {
		t.Fatalf("CompleteWorkerRequest() error = %v", err)
	}
	requestPage, err = repository.ListWorkerRequests(ctx, expertID, "completed", 20, 0)
	if err != nil || requestPage.Total != 1 || requestPage.Items[0].Status != "completed" {
		t.Fatalf("ListWorkerRequests(completed) = (%+v, %v)", requestPage, err)
	}
	answer, err := repository.AddAnswer(ctx, trackID, expertID, "Ответ", createdAt.Add(6*time.Minute))
	if err != nil || answer.TicketStatus != models.TicketStatusAnswerReady {
		t.Fatalf("AddAnswer() = (%+v, %v)", answer, err)
	}
	analytics, err := repository.Analytics(ctx, expertID, createdAt.Add(-time.Hour), createdAt.Add(2*time.Hour))
	if err != nil || analytics.Total != 1 || analytics.UrgentCount != 1 {
		t.Fatalf("Analytics() = (%+v, %v)", analytics, err)
	}
}

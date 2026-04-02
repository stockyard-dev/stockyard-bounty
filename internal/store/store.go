package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct{ db *sql.DB }

type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at"`
	OpenCount   int    `json:"open_count"`
	ClosedCount int    `json:"closed_count"`
}

type Issue struct {
	ID           string   `json:"id"`
	ProjectID    string   `json:"project_id"`
	Number       int      `json:"number"`
	Title        string   `json:"title"`
	Body         string   `json:"body,omitempty"`
	Status       string   `json:"status"`
	Priority     string   `json:"priority"`
	Labels       []string `json:"labels"`
	MilestoneID  string   `json:"milestone_id,omitempty"`
	Assignee     string   `json:"assignee,omitempty"`
	TimeEstimate int      `json:"time_estimate,omitempty"`
	TimeSpent    int      `json:"time_spent,omitempty"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
	ClosedAt     string   `json:"closed_at,omitempty"`
	CommentCount int      `json:"comment_count"`
}

type Comment struct {
	ID        string `json:"id"`
	IssueID   string `json:"issue_id"`
	Author    string `json:"author,omitempty"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}

type Milestone struct {
	ID          string `json:"id"`
	ProjectID   string `json:"project_id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	DueDate     string `json:"due_date,omitempty"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	OpenCount   int    `json:"open_count"`
	ClosedCount int    `json:"closed_count"`
}

type Activity struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
	IssueID   string `json:"issue_id,omitempty"`
	Action    string `json:"action"`
	Detail    string `json:"detail,omitempty"`
	Actor     string `json:"actor,omitempty"`
	CreatedAt string `json:"created_at"`
}

type IssueFilter struct {
	ProjectID   string
	Status      string
	Priority    string
	Label       string
	MilestoneID string
	Assignee    string
	Search      string
	SortBy      string
	SortDir     string
	Limit       int
	Offset      int
}

func Open(dataDir string) (*DB, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}
	dsn := filepath.Join(dataDir, "bounty.db") + "?_journal_mode=WAL&_busy_timeout=5000"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	for _, q := range []string{
		`CREATE TABLE IF NOT EXISTS projects (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			slug TEXT UNIQUE NOT NULL,
			description TEXT DEFAULT '',
			created_at TEXT DEFAULT (datetime('now'))
		)`,
		`CREATE TABLE IF NOT EXISTS issues (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL REFERENCES projects(id),
			number INTEGER NOT NULL,
			title TEXT NOT NULL,
			body TEXT DEFAULT '',
			status TEXT DEFAULT 'open',
			priority TEXT DEFAULT 'medium',
			labels_json TEXT DEFAULT '[]',
			milestone_id TEXT DEFAULT '',
			assignee TEXT DEFAULT '',
			time_estimate INTEGER DEFAULT 0,
			time_spent INTEGER DEFAULT 0,
			created_at TEXT DEFAULT (datetime('now')),
			updated_at TEXT DEFAULT (datetime('now')),
			closed_at TEXT DEFAULT ''
		)`,
		`CREATE TABLE IF NOT EXISTS comments (
			id TEXT PRIMARY KEY,
			issue_id TEXT NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
			author TEXT DEFAULT '',
			body TEXT NOT NULL,
			created_at TEXT DEFAULT (datetime('now'))
		)`,
		`CREATE TABLE IF NOT EXISTS milestones (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL REFERENCES projects(id),
			title TEXT NOT NULL,
			description TEXT DEFAULT '',
			due_date TEXT DEFAULT '',
			status TEXT DEFAULT 'open',
			created_at TEXT DEFAULT (datetime('now'))
		)`,
		`CREATE TABLE IF NOT EXISTS activity (
			id TEXT PRIMARY KEY,
			project_id TEXT DEFAULT '',
			issue_id TEXT DEFAULT '',
			action TEXT NOT NULL,
			detail TEXT DEFAULT '',
			actor TEXT DEFAULT '',
			created_at TEXT DEFAULT (datetime('now'))
		)`,
		`CREATE INDEX IF NOT EXISTS idx_issues_project ON issues(project_id)`,
		`CREATE INDEX IF NOT EXISTS idx_issues_status ON issues(status)`,
		`CREATE INDEX IF NOT EXISTS idx_issues_milestone ON issues(milestone_id)`,
		`CREATE INDEX IF NOT EXISTS idx_comments_issue ON comments(issue_id)`,
		`CREATE INDEX IF NOT EXISTS idx_milestones_project ON milestones(project_id)`,
		`CREATE INDEX IF NOT EXISTS idx_activity_project ON activity(project_id)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_issues_project_number ON issues(project_id, number)`,
	} {
		if _, err := db.Exec(q); err != nil {
			return nil, fmt.Errorf("migrate: %w", err)
		}
	}
	return &DB{db: db}, nil
}

func (d *DB) Close() error { return d.db.Close() }
func genID() string        { return fmt.Sprintf("%d", time.Now().UnixNano()) }
func now() string          { return time.Now().UTC().Format(time.RFC3339) }

// ── Projects ──

func (d *DB) CreateProject(p *Project) error {
	p.ID = genID()
	p.CreatedAt = now()
	if p.Slug == "" {
		p.Slug = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(p.Name), " ", "-"))
	}
	_, err := d.db.Exec(`INSERT INTO projects (id,name,slug,description,created_at) VALUES (?,?,?,?,?)`,
		p.ID, p.Name, p.Slug, p.Description, p.CreatedAt)
	if err != nil {
		return err
	}
	d.log(p.ID, "", "created", "Project created: "+p.Name, "")
	return nil
}

func (d *DB) GetProject(id string) *Project {
	var p Project
	if err := d.db.QueryRow(`SELECT id,name,slug,description,created_at FROM projects WHERE id=?`, id).Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.CreatedAt); err != nil {
		return nil
	}
	d.db.QueryRow(`SELECT COUNT(*) FROM issues WHERE project_id=? AND status='open'`, id).Scan(&p.OpenCount)
	d.db.QueryRow(`SELECT COUNT(*) FROM issues WHERE project_id=? AND status='closed'`, id).Scan(&p.ClosedCount)
	return &p
}

func (d *DB) ListProjects() []Project {
	rows, err := d.db.Query(`SELECT id,name,slug,description,created_at FROM projects ORDER BY created_at DESC`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Slug, &p.Description, &p.CreatedAt); err != nil {
			continue
		}
		d.db.QueryRow(`SELECT COUNT(*) FROM issues WHERE project_id=? AND status='open'`, p.ID).Scan(&p.OpenCount)
		d.db.QueryRow(`SELECT COUNT(*) FROM issues WHERE project_id=? AND status='closed'`, p.ID).Scan(&p.ClosedCount)
		out = append(out, p)
	}
	return out
}

func (d *DB) UpdateProject(id string, p *Project) error {
	_, err := d.db.Exec(`UPDATE projects SET name=?,slug=?,description=? WHERE id=?`, p.Name, p.Slug, p.Description, id)
	return err
}

func (d *DB) DeleteProject(id string) error {
	d.db.Exec(`DELETE FROM comments WHERE issue_id IN (SELECT id FROM issues WHERE project_id=?)`, id)
	d.db.Exec(`DELETE FROM issues WHERE project_id=?`, id)
	d.db.Exec(`DELETE FROM milestones WHERE project_id=?`, id)
	d.db.Exec(`DELETE FROM activity WHERE project_id=?`, id)
	_, err := d.db.Exec(`DELETE FROM projects WHERE id=?`, id)
	return err
}

// ── Issues ──

func (d *DB) nextNumber(projectID string) int {
	var n int
	d.db.QueryRow(`SELECT COALESCE(MAX(number),0) FROM issues WHERE project_id=?`, projectID).Scan(&n)
	return n + 1
}

func (d *DB) CreateIssue(e *Issue) error {
	e.ID = genID()
	e.Number = d.nextNumber(e.ProjectID)
	e.CreatedAt = now()
	e.UpdatedAt = e.CreatedAt
	if e.Status == "" {
		e.Status = "open"
	}
	if e.Priority == "" {
		e.Priority = "medium"
	}
	if e.Labels == nil {
		e.Labels = []string{}
	}
	lj, _ := json.Marshal(e.Labels)
	_, err := d.db.Exec(`INSERT INTO issues (id,project_id,number,title,body,status,priority,labels_json,milestone_id,assignee,time_estimate,time_spent,created_at,updated_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		e.ID, e.ProjectID, e.Number, e.Title, e.Body, e.Status, e.Priority, string(lj), e.MilestoneID, e.Assignee, e.TimeEstimate, e.TimeSpent, e.CreatedAt, e.UpdatedAt)
	if err != nil {
		return err
	}
	d.log(e.ProjectID, e.ID, "created", fmt.Sprintf("#%d %s", e.Number, e.Title), "")
	return nil
}

func (d *DB) scanIssue(s interface{ Scan(...any) error }) *Issue {
	var e Issue
	var lj string
	if err := s.Scan(&e.ID, &e.ProjectID, &e.Number, &e.Title, &e.Body, &e.Status, &e.Priority, &lj, &e.MilestoneID, &e.Assignee, &e.TimeEstimate, &e.TimeSpent, &e.CreatedAt, &e.UpdatedAt, &e.ClosedAt); err != nil {
		return nil
	}
	json.Unmarshal([]byte(lj), &e.Labels)
	if e.Labels == nil {
		e.Labels = []string{}
	}
	d.db.QueryRow(`SELECT COUNT(*) FROM comments WHERE issue_id=?`, e.ID).Scan(&e.CommentCount)
	return &e
}

const issueCols = `id,project_id,number,title,body,status,priority,labels_json,milestone_id,assignee,time_estimate,time_spent,created_at,updated_at,closed_at`

func (d *DB) GetIssue(id string) *Issue {
	return d.scanIssue(d.db.QueryRow(`SELECT `+issueCols+` FROM issues WHERE id=?`, id))
}

func (d *DB) ListIssues(f IssueFilter) ([]Issue, int) {
	where := []string{"1=1"}
	args := []any{}
	if f.ProjectID != "" {
		where = append(where, "project_id=?")
		args = append(args, f.ProjectID)
	}
	if f.Status != "" && f.Status != "all" {
		where = append(where, "status=?")
		args = append(args, f.Status)
	}
	if f.Priority != "" {
		where = append(where, "priority=?")
		args = append(args, f.Priority)
	}
	if f.Label != "" {
		where = append(where, `labels_json LIKE ?`)
		args = append(args, `%"`+f.Label+`"%`)
	}
	if f.MilestoneID != "" {
		where = append(where, "milestone_id=?")
		args = append(args, f.MilestoneID)
	}
	if f.Assignee != "" {
		where = append(where, "assignee=?")
		args = append(args, f.Assignee)
	}
	if f.Search != "" {
		where = append(where, "(title LIKE ? OR body LIKE ?)")
		s := "%" + f.Search + "%"
		args = append(args, s, s)
	}
	w := strings.Join(where, " AND ")
	var total int
	d.db.QueryRow("SELECT COUNT(*) FROM issues WHERE "+w, args...).Scan(&total)

	order := "created_at"
	switch f.SortBy {
	case "updated":
		order = "updated_at"
	case "priority":
		order = "CASE priority WHEN 'critical' THEN 0 WHEN 'high' THEN 1 WHEN 'medium' THEN 2 WHEN 'low' THEN 3 END"
	case "comments":
		order = "(SELECT COUNT(*) FROM comments c WHERE c.issue_id=issues.id)"
	}
	dir := "DESC"
	if f.SortDir == "asc" {
		dir = "ASC"
	}
	if f.Limit <= 0 {
		f.Limit = 50
	}
	q := fmt.Sprintf("SELECT %s FROM issues WHERE %s ORDER BY %s %s LIMIT ? OFFSET ?", issueCols, w, order, dir)
	args = append(args, f.Limit, f.Offset)
	rows, err := d.db.Query(q, args...)
	if err != nil {
		return nil, 0
	}
	defer rows.Close()
	var out []Issue
	for rows.Next() {
		if e := d.scanIssue(rows); e != nil {
			out = append(out, *e)
		}
	}
	return out, total
}

func (d *DB) UpdateIssue(id string, e *Issue) error {
	e.UpdatedAt = now()
	lj, _ := json.Marshal(e.Labels)
	_, err := d.db.Exec(`UPDATE issues SET title=?,body=?,priority=?,labels_json=?,milestone_id=?,assignee=?,time_estimate=?,time_spent=?,updated_at=? WHERE id=?`,
		e.Title, e.Body, e.Priority, string(lj), e.MilestoneID, e.Assignee, e.TimeEstimate, e.TimeSpent, e.UpdatedAt, id)
	return err
}

func (d *DB) CloseIssue(id string) error {
	t := now()
	if _, err := d.db.Exec(`UPDATE issues SET status='closed',closed_at=?,updated_at=? WHERE id=?`, t, t, id); err != nil {
		return err
	}
	if e := d.GetIssue(id); e != nil {
		d.log(e.ProjectID, id, "closed", fmt.Sprintf("#%d closed", e.Number), "")
	}
	return nil
}

func (d *DB) ReopenIssue(id string) error {
	t := now()
	if _, err := d.db.Exec(`UPDATE issues SET status='open',closed_at='',updated_at=? WHERE id=?`, t, id); err != nil {
		return err
	}
	if e := d.GetIssue(id); e != nil {
		d.log(e.ProjectID, id, "reopened", fmt.Sprintf("#%d reopened", e.Number), "")
	}
	return nil
}

func (d *DB) DeleteIssue(id string) error {
	d.db.Exec(`DELETE FROM comments WHERE issue_id=?`, id)
	_, err := d.db.Exec(`DELETE FROM issues WHERE id=?`, id)
	return err
}

// ── Comments ──

func (d *DB) CreateComment(c *Comment) error {
	c.ID = genID()
	c.CreatedAt = now()
	_, err := d.db.Exec(`INSERT INTO comments (id,issue_id,author,body,created_at) VALUES (?,?,?,?,?)`,
		c.ID, c.IssueID, c.Author, c.Body, c.CreatedAt)
	if err != nil {
		return err
	}
	d.db.Exec(`UPDATE issues SET updated_at=? WHERE id=?`, c.CreatedAt, c.IssueID)
	if e := d.GetIssue(c.IssueID); e != nil {
		d.log(e.ProjectID, c.IssueID, "commented", fmt.Sprintf("#%d comment by %s", e.Number, c.Author), c.Author)
	}
	return nil
}

func (d *DB) ListComments(issueID string) []Comment {
	rows, err := d.db.Query(`SELECT id,issue_id,author,body,created_at FROM comments WHERE issue_id=? ORDER BY created_at ASC`, issueID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.IssueID, &c.Author, &c.Body, &c.CreatedAt); err != nil {
			continue
		}
		out = append(out, c)
	}
	return out
}

func (d *DB) DeleteComment(id string) error {
	_, err := d.db.Exec(`DELETE FROM comments WHERE id=?`, id)
	return err
}

// ── Milestones ──

func (d *DB) CreateMilestone(m *Milestone) error {
	m.ID = genID()
	m.CreatedAt = now()
	if m.Status == "" {
		m.Status = "open"
	}
	_, err := d.db.Exec(`INSERT INTO milestones (id,project_id,title,description,due_date,status,created_at) VALUES (?,?,?,?,?,?,?)`,
		m.ID, m.ProjectID, m.Title, m.Description, m.DueDate, m.Status, m.CreatedAt)
	return err
}

func (d *DB) GetMilestone(id string) *Milestone {
	var m Milestone
	if err := d.db.QueryRow(`SELECT id,project_id,title,description,due_date,status,created_at FROM milestones WHERE id=?`, id).Scan(&m.ID, &m.ProjectID, &m.Title, &m.Description, &m.DueDate, &m.Status, &m.CreatedAt); err != nil {
		return nil
	}
	d.db.QueryRow(`SELECT COUNT(*) FROM issues WHERE milestone_id=? AND status='open'`, id).Scan(&m.OpenCount)
	d.db.QueryRow(`SELECT COUNT(*) FROM issues WHERE milestone_id=? AND status='closed'`, id).Scan(&m.ClosedCount)
	return &m
}

func (d *DB) ListMilestones(projectID string) []Milestone {
	rows, err := d.db.Query(`SELECT id,project_id,title,description,due_date,status,created_at FROM milestones WHERE project_id=? ORDER BY due_date ASC, created_at DESC`, projectID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []Milestone
	for rows.Next() {
		var m Milestone
		if err := rows.Scan(&m.ID, &m.ProjectID, &m.Title, &m.Description, &m.DueDate, &m.Status, &m.CreatedAt); err != nil {
			continue
		}
		d.db.QueryRow(`SELECT COUNT(*) FROM issues WHERE milestone_id=? AND status='open'`, m.ID).Scan(&m.OpenCount)
		d.db.QueryRow(`SELECT COUNT(*) FROM issues WHERE milestone_id=? AND status='closed'`, m.ID).Scan(&m.ClosedCount)
		out = append(out, m)
	}
	return out
}

func (d *DB) UpdateMilestone(id string, m *Milestone) error {
	_, err := d.db.Exec(`UPDATE milestones SET title=?,description=?,due_date=?,status=? WHERE id=?`,
		m.Title, m.Description, m.DueDate, m.Status, id)
	return err
}

func (d *DB) DeleteMilestone(id string) error {
	d.db.Exec(`UPDATE issues SET milestone_id='' WHERE milestone_id=?`, id)
	_, err := d.db.Exec(`DELETE FROM milestones WHERE id=?`, id)
	return err
}

// ── Activity ──

func (d *DB) log(projectID, issueID, action, detail, actor string) {
	d.db.Exec(`INSERT INTO activity (id,project_id,issue_id,action,detail,actor,created_at) VALUES (?,?,?,?,?,?,?)`,
		genID(), projectID, issueID, action, detail, actor, now())
}

func (d *DB) ListActivity(projectID string, limit int) []Activity {
	if limit <= 0 {
		limit = 50
	}
	q := `SELECT id,project_id,issue_id,action,detail,actor,created_at FROM activity`
	var args []any
	if projectID != "" {
		q += ` WHERE project_id=?`
		args = append(args, projectID)
	}
	q += ` ORDER BY created_at DESC LIMIT ?`
	args = append(args, limit)
	rows, err := d.db.Query(q, args...)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []Activity
	for rows.Next() {
		var a Activity
		if err := rows.Scan(&a.ID, &a.ProjectID, &a.IssueID, &a.Action, &a.Detail, &a.Actor, &a.CreatedAt); err != nil {
			continue
		}
		out = append(out, a)
	}
	return out
}

// ── Stats ──

type Stats struct {
	Projects     int            `json:"projects"`
	OpenIssues   int            `json:"open_issues"`
	ClosedIssues int            `json:"closed_issues"`
	Comments     int            `json:"comments"`
	Milestones   int            `json:"milestones"`
	ByPriority   map[string]int `json:"by_priority"`
}

func (d *DB) Stats() Stats {
	var s Stats
	d.db.QueryRow(`SELECT COUNT(*) FROM projects`).Scan(&s.Projects)
	d.db.QueryRow(`SELECT COUNT(*) FROM issues WHERE status='open'`).Scan(&s.OpenIssues)
	d.db.QueryRow(`SELECT COUNT(*) FROM issues WHERE status='closed'`).Scan(&s.ClosedIssues)
	d.db.QueryRow(`SELECT COUNT(*) FROM comments`).Scan(&s.Comments)
	d.db.QueryRow(`SELECT COUNT(*) FROM milestones`).Scan(&s.Milestones)
	s.ByPriority = map[string]int{}
	rows, _ := d.db.Query(`SELECT priority, COUNT(*) FROM issues WHERE status='open' GROUP BY priority`)
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var p string
			var c int
			rows.Scan(&p, &c)
			s.ByPriority[p] = c
		}
	}
	return s
}

func (d *DB) AllLabels() []string {
	rows, err := d.db.Query(`SELECT DISTINCT labels_json FROM issues WHERE labels_json != '[]'`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	seen := map[string]bool{}
	for rows.Next() {
		var j string
		rows.Scan(&j)
		var labels []string
		json.Unmarshal([]byte(j), &labels)
		for _, l := range labels {
			seen[l] = true
		}
	}
	var out []string
	for l := range seen {
		out = append(out, l)
	}
	return out
}

func (d *DB) AllAssignees() []string {
	rows, err := d.db.Query(`SELECT DISTINCT assignee FROM issues WHERE assignee != '' ORDER BY assignee`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var a string
		rows.Scan(&a)
		out = append(out, a)
	}
	return out
}

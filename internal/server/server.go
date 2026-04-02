package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/stockyard-dev/stockyard-bounty/internal/store"
)

type Server struct {
	db  *store.DB
	mux *http.ServeMux
}

func New(db *store.DB) *Server {
	s := &Server{db: db, mux: http.NewServeMux()}

	// Projects
	s.mux.HandleFunc("GET /api/projects", s.listProjects)
	s.mux.HandleFunc("POST /api/projects", s.createProject)
	s.mux.HandleFunc("GET /api/projects/{id}", s.getProject)
	s.mux.HandleFunc("PUT /api/projects/{id}", s.updateProject)
	s.mux.HandleFunc("DELETE /api/projects/{id}", s.deleteProject)

	// Issues
	s.mux.HandleFunc("GET /api/issues", s.listIssues)
	s.mux.HandleFunc("POST /api/issues", s.createIssue)
	s.mux.HandleFunc("GET /api/issues/{id}", s.getIssue)
	s.mux.HandleFunc("PUT /api/issues/{id}", s.updateIssue)
	s.mux.HandleFunc("DELETE /api/issues/{id}", s.deleteIssue)
	s.mux.HandleFunc("POST /api/issues/{id}/close", s.closeIssue)
	s.mux.HandleFunc("POST /api/issues/{id}/reopen", s.reopenIssue)

	// Comments
	s.mux.HandleFunc("GET /api/issues/{id}/comments", s.listComments)
	s.mux.HandleFunc("POST /api/issues/{id}/comments", s.createComment)
	s.mux.HandleFunc("DELETE /api/comments/{id}", s.deleteComment)

	// Milestones
	s.mux.HandleFunc("GET /api/projects/{id}/milestones", s.listMilestones)
	s.mux.HandleFunc("POST /api/milestones", s.createMilestone)
	s.mux.HandleFunc("GET /api/milestones/{id}", s.getMilestone)
	s.mux.HandleFunc("PUT /api/milestones/{id}", s.updateMilestone)
	s.mux.HandleFunc("DELETE /api/milestones/{id}", s.deleteMilestone)

	// Activity & Meta
	s.mux.HandleFunc("GET /api/activity", s.listActivity)
	s.mux.HandleFunc("GET /api/stats", s.stats)
	s.mux.HandleFunc("GET /api/labels", s.allLabels)
	s.mux.HandleFunc("GET /api/assignees", s.allAssignees)
	s.mux.HandleFunc("GET /api/health", s.health)

	// Dashboard
	s.mux.HandleFunc("GET /ui", s.dashboard)
	s.mux.HandleFunc("GET /ui/", s.dashboard)
	s.mux.HandleFunc("GET /", s.root)

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) { s.mux.ServeHTTP(w, r) }

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

func (s *Server) root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "/ui", http.StatusFound)
}

// ── Projects ──

func (s *Server) listProjects(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"projects": orEmpty(s.db.ListProjects())})
}

func (s *Server) createProject(w http.ResponseWriter, r *http.Request) {
	var p store.Project
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeErr(w, 400, "invalid json")
		return
	}
	if p.Name == "" {
		writeErr(w, 400, "name required")
		return
	}
	if err := s.db.CreateProject(&p); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 201, p)
}

func (s *Server) getProject(w http.ResponseWriter, r *http.Request) {
	p := s.db.GetProject(r.PathValue("id"))
	if p == nil {
		writeErr(w, 404, "not found")
		return
	}
	writeJSON(w, 200, p)
}

func (s *Server) updateProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	existing := s.db.GetProject(id)
	if existing == nil {
		writeErr(w, 404, "not found")
		return
	}
	var p store.Project
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeErr(w, 400, "invalid json")
		return
	}
	if p.Name == "" {
		p.Name = existing.Name
	}
	if p.Slug == "" {
		p.Slug = existing.Slug
	}
	if err := s.db.UpdateProject(id, &p); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, s.db.GetProject(id))
}

func (s *Server) deleteProject(w http.ResponseWriter, r *http.Request) {
	if err := s.db.DeleteProject(r.PathValue("id")); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, map[string]string{"deleted": "ok"})
}

// ── Issues ──

func (s *Server) listIssues(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	f := store.IssueFilter{
		ProjectID:   q.Get("project_id"),
		Status:      q.Get("status"),
		Priority:    q.Get("priority"),
		Label:       q.Get("label"),
		MilestoneID: q.Get("milestone_id"),
		Assignee:    q.Get("assignee"),
		Search:      q.Get("search"),
		SortBy:      q.Get("sort"),
		SortDir:     q.Get("dir"),
		Limit:       limit,
		Offset:      offset,
	}
	if f.Status == "" {
		f.Status = "open"
	}
	issues, total := s.db.ListIssues(f)
	writeJSON(w, 200, map[string]any{"issues": orEmpty(issues), "total": total})
}

func (s *Server) createIssue(w http.ResponseWriter, r *http.Request) {
	var e store.Issue
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		writeErr(w, 400, "invalid json")
		return
	}
	if e.Title == "" {
		writeErr(w, 400, "title required")
		return
	}
	if e.ProjectID == "" {
		writeErr(w, 400, "project_id required")
		return
	}
	if s.db.GetProject(e.ProjectID) == nil {
		writeErr(w, 404, "project not found")
		return
	}
	if err := s.db.CreateIssue(&e); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 201, e)
}

func (s *Server) getIssue(w http.ResponseWriter, r *http.Request) {
	e := s.db.GetIssue(r.PathValue("id"))
	if e == nil {
		writeErr(w, 404, "not found")
		return
	}
	writeJSON(w, 200, e)
}

func (s *Server) updateIssue(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	existing := s.db.GetIssue(id)
	if existing == nil {
		writeErr(w, 404, "not found")
		return
	}
	var e store.Issue
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		writeErr(w, 400, "invalid json")
		return
	}
	if e.Title == "" {
		e.Title = existing.Title
	}
	if e.Body == "" {
		e.Body = existing.Body
	}
	if e.Priority == "" {
		e.Priority = existing.Priority
	}
	if e.Labels == nil {
		e.Labels = existing.Labels
	}
	if e.MilestoneID == "" {
		e.MilestoneID = existing.MilestoneID
	}
	if e.Assignee == "" {
		e.Assignee = existing.Assignee
	}
	if e.TimeEstimate == 0 {
		e.TimeEstimate = existing.TimeEstimate
	}
	if e.TimeSpent == 0 {
		e.TimeSpent = existing.TimeSpent
	}
	if err := s.db.UpdateIssue(id, &e); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, s.db.GetIssue(id))
}

func (s *Server) deleteIssue(w http.ResponseWriter, r *http.Request) {
	if err := s.db.DeleteIssue(r.PathValue("id")); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, map[string]string{"deleted": "ok"})
}

func (s *Server) closeIssue(w http.ResponseWriter, r *http.Request) {
	if err := s.db.CloseIssue(r.PathValue("id")); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, s.db.GetIssue(r.PathValue("id")))
}

func (s *Server) reopenIssue(w http.ResponseWriter, r *http.Request) {
	if err := s.db.ReopenIssue(r.PathValue("id")); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, s.db.GetIssue(r.PathValue("id")))
}

// ── Comments ──

func (s *Server) listComments(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if s.db.GetIssue(id) == nil {
		writeErr(w, 404, "issue not found")
		return
	}
	writeJSON(w, 200, map[string]any{"comments": orEmpty(s.db.ListComments(id))})
}

func (s *Server) createComment(w http.ResponseWriter, r *http.Request) {
	issueID := r.PathValue("id")
	if s.db.GetIssue(issueID) == nil {
		writeErr(w, 404, "issue not found")
		return
	}
	var c store.Comment
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		writeErr(w, 400, "invalid json")
		return
	}
	if c.Body == "" {
		writeErr(w, 400, "body required")
		return
	}
	c.IssueID = issueID
	if err := s.db.CreateComment(&c); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 201, c)
}

func (s *Server) deleteComment(w http.ResponseWriter, r *http.Request) {
	if err := s.db.DeleteComment(r.PathValue("id")); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, map[string]string{"deleted": "ok"})
}

// ── Milestones ──

func (s *Server) listMilestones(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"milestones": orEmpty(s.db.ListMilestones(r.PathValue("id")))})
}

func (s *Server) createMilestone(w http.ResponseWriter, r *http.Request) {
	var m store.Milestone
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		writeErr(w, 400, "invalid json")
		return
	}
	if m.Title == "" {
		writeErr(w, 400, "title required")
		return
	}
	if m.ProjectID == "" {
		writeErr(w, 400, "project_id required")
		return
	}
	if err := s.db.CreateMilestone(&m); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 201, m)
}

func (s *Server) getMilestone(w http.ResponseWriter, r *http.Request) {
	m := s.db.GetMilestone(r.PathValue("id"))
	if m == nil {
		writeErr(w, 404, "not found")
		return
	}
	writeJSON(w, 200, m)
}

func (s *Server) updateMilestone(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	existing := s.db.GetMilestone(id)
	if existing == nil {
		writeErr(w, 404, "not found")
		return
	}
	var m store.Milestone
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		writeErr(w, 400, "invalid json")
		return
	}
	if m.Title == "" {
		m.Title = existing.Title
	}
	if err := s.db.UpdateMilestone(id, &m); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, s.db.GetMilestone(id))
}

func (s *Server) deleteMilestone(w http.ResponseWriter, r *http.Request) {
	if err := s.db.DeleteMilestone(r.PathValue("id")); err != nil {
		writeErr(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, map[string]string{"deleted": "ok"})
}

// ── Activity & Meta ──

func (s *Server) listActivity(w http.ResponseWriter, r *http.Request) {
	pid := r.URL.Query().Get("project_id")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	writeJSON(w, 200, map[string]any{"activity": orEmpty(s.db.ListActivity(pid, limit))})
}

func (s *Server) stats(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, s.db.Stats())
}

func (s *Server) allLabels(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"labels": orEmpty(s.db.AllLabels())})
}

func (s *Server) allAssignees(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"assignees": orEmpty(s.db.AllAssignees())})
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	st := s.db.Stats()
	writeJSON(w, 200, map[string]any{
		"status":   "ok",
		"service":  "bounty",
		"projects": st.Projects,
		"issues":   st.OpenIssues + st.ClosedIssues,
	})
}

func orEmpty[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

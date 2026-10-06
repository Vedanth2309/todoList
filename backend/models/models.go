package models

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Base struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}

func (b *Base) B() *Base { return b }

type SearchDoc struct {
	Type, Title, Content string
	Tags                 []string
	Date                 string
}

// Doc is implemented by every user-owned entity handled by the generic CRUD controller.
type Doc interface {
	B() *Base
	Validate() string // returns an error message, "" if valid; may set defaults
	Search() *SearchDoc // nil if not indexed in Elasticsearch
}

func blank(s string) bool { return strings.TrimSpace(s) == "" }
func def(s *string, v string) {
	if *s == "" {
		*s = v
	}
}

type Task struct {
	Base        `bson:",inline"`
	Title       string     `bson:"title" json:"title"`
	Description string     `bson:"description" json:"description"`
	Status      string     `bson:"status" json:"status"`
	Priority    string     `bson:"priority" json:"priority"`
	DueDate     string     `bson:"dueDate" json:"dueDate"`
	DueTime     string     `bson:"dueTime" json:"dueTime"`
	Category    string     `bson:"category" json:"category"`
	Tags        []string   `bson:"tags" json:"tags"`
	CompletedAt *time.Time `bson:"completedAt,omitempty" json:"completedAt"`
}

func (t *Task) Validate() string {
	if blank(t.Title) {
		return "Title is required"
	}
	def(&t.Status, "TODO")
	def(&t.Priority, "MEDIUM")
	return ""
}
func (t *Task) Search() *SearchDoc {
	d := t.DueDate
	if d == "" {
		d = t.CreatedAt.Format("2006-01-02")
	}
	return &SearchDoc{"task", t.Title, t.Description, t.Tags, d}
}

type Habit struct {
	Base        `bson:",inline"`
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	Category    string `bson:"category" json:"category"`
	Frequency   string `bson:"frequency" json:"frequency"`
	Target      int    `bson:"target" json:"target"`
	StartDate   string `bson:"startDate" json:"startDate"`
	Color       string `bson:"color" json:"color"`
	Icon        string `bson:"icon" json:"icon"`
}

func (h *Habit) Validate() string {
	if blank(h.Name) {
		return "Name is required"
	}
	def(&h.Frequency, "daily")
	return ""
}
func (h *Habit) Search() *SearchDoc { return nil }

type HabitCheckin struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	UserID    primitive.ObjectID `bson:"userId" json:"userId"`
	HabitID   primitive.ObjectID `bson:"habitId" json:"habitId"`
	Date      string             `bson:"date" json:"date"`
	Completed bool               `bson:"completed" json:"completed"`
	Note      string             `bson:"note" json:"note"`
}

type Event struct {
	Base          `bson:",inline"`
	Title         string `bson:"title" json:"title"`
	Description   string `bson:"description" json:"description"`
	Date          string `bson:"date" json:"date"`
	StartDateTime string `bson:"startDateTime" json:"startDateTime"`
	EndDateTime   string `bson:"endDateTime" json:"endDateTime"`
	Location      string `bson:"location" json:"location"`
	Category      string `bson:"category" json:"category"`
	Reminder      string `bson:"reminder" json:"reminder"`
	Color         string `bson:"color" json:"color"`
}

func (e *Event) Validate() string {
	if blank(e.Title) {
		return "Title is required"
	}
	if e.Date == "" && len(e.StartDateTime) >= 10 {
		e.Date = e.StartDateTime[:10]
	}
	return ""
}
func (e *Event) Search() *SearchDoc { return nil }

type Reminder struct {
	Base           `bson:",inline"`
	Title          string `bson:"title" json:"title"`
	Description    string `bson:"description" json:"description"`
	Date           string `bson:"date" json:"date"`
	Time           string `bson:"time" json:"time"`
	Completed      bool   `bson:"completed" json:"completed"`
	RelatedTaskID  string `bson:"relatedTaskId" json:"relatedTaskId"`
	RelatedEventID string `bson:"relatedEventId" json:"relatedEventId"`
}

func (r *Reminder) Validate() string {
	if blank(r.Title) {
		return "Title is required"
	}
	return ""
}
func (r *Reminder) Search() *SearchDoc { return nil }

type Diary struct {
	Base    `bson:",inline"`
	Title   string   `bson:"title" json:"title"`
	Content string   `bson:"content" json:"content"`
	Date    string   `bson:"date" json:"date"`
	Mood    string   `bson:"mood" json:"mood"`
	Tags    []string `bson:"tags" json:"tags"`
}

func (d *Diary) Validate() string {
	if blank(d.Title) {
		return "Title is required"
	}
	def(&d.Date, time.Now().Format("2006-01-02"))
	return ""
}
func (d *Diary) Search() *SearchDoc { return &SearchDoc{"diary", d.Title, d.Content, d.Tags, d.Date} }

type Note struct {
	Base    `bson:",inline"`
	Title   string   `bson:"title" json:"title"`
	Content string   `bson:"content" json:"content"`
	Tags    []string `bson:"tags" json:"tags"`
}

func (n *Note) Validate() string {
	if blank(n.Title) {
		return "Title is required"
	}
	return ""
}
func (n *Note) Search() *SearchDoc {
	return &SearchDoc{"note", n.Title, n.Content, n.Tags, n.CreatedAt.Format("2006-01-02")}
}

type Settings struct {
	Theme           string `bson:"theme" json:"theme"`
	Timezone        string `bson:"timezone" json:"timezone"`
	DefaultPriority string `bson:"defaultPriority" json:"defaultPriority"`
	Notifications   bool   `bson:"notifications" json:"notifications"`
}

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
	Settings  Settings           `bson:"settings" json:"settings"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}

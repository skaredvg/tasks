package postgres

import (
	"context"
	"strconv"
	"strings"
	"tasks/dbflow"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(connstr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}
	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) NewTask(t dbflow.Task) (int, error) {
	sql := `INSERT INTO tasks.public.tasks (author_id, assigned_id, title, content)
			VALUES ($1, $2, $3, $4) RETURNING id`
	id := 0
	err := s.db.QueryRow(context.Background(), sql, t.Author.ID(), t.Assign.ID(), t.Title, t.Content).Scan(&id)
	t.SetID(id)
	return id, err
}

func (s *Storage) UpdateTask(t dbflow.Task) error {
	sql := `UPDATE tasks.public.tasks
			SET author_id = $1,
				assigned_id = $2,
				title = $3,
				content = $4
			WHERE id = $5`

	_, err := s.db.Exec(context.Background(), sql, t.Author.ID(), t.Assign.ID(), t.Title, t.Content)
	return err
}

func (s *Storage) DeleteTask(id int) error {
	sql := "DELETE FROM tasks.public.tasks WHERE id = $1"

	_, err := s.db.Exec(context.Background(), sql, id)
	return err
}

func convToSliceLabel(strlbl string) []dbflow.Label {
	sl := make([]dbflow.Label, 3)
	for _, v := range strings.Split(strlbl, ";") {
		v1 := strings.Split(v, "|")
		lo := dbflow.Label{Name: v1[1]}
		id, _ := strconv.Atoi((v1[0]))
		lo.SetID(id)
		sl = append(sl, lo)
	}
	return sl
}

func tasksByRows(s *Storage, sql string, par ...any) (dbflow.Tasks, error) {
	rows, err := s.db.Query(context.Background(), sql, par...)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	if err != nil {
		return make(dbflow.Tasks, 0), err
	}

	m := make(dbflow.Tasks, 0)

	for rows.Next() {
		t := dbflow.Task{}

		id := 0
		opened := 0
		closed := 0
		aut_id := 0
		aut_name := ""
		ass_id := 0
		ass_name := ""
		title := ""
		content := ""
		lbl_str := ""
		err := rows.Scan(&id, &opened, &closed, &aut_id, &ass_id, &title, &content, &aut_name, &ass_name, &lbl_str)

		if err != nil {
			return make(dbflow.Tasks, 0), err
		}

		t.SetID(id)
		t.SetOpened(opened)
		t.SetClosed(closed)
		t.Author = dbflow.User{Name: aut_name}
		t.Author.SetID(aut_id)
		t.Assign = dbflow.User{Name: ass_name}
		t.Assign.SetID(ass_id)
		t.Title = title
		t.Content = content

		t.Labels = convToSliceLabel(lbl_str)

		m = append(m, t)
	}
	return m, nil
}

func (s *Storage) TasksByAuthor(id int) (dbflow.Tasks, error) {
	sql := `SELECT t.id,
					t.opened,
					t.closed,
					t.author_id,
					t.assigned_id,
					t.title,
					t.content,
					(SELECT u.name FROM tasks.public.users u WHERE u.id = t.author_id) AS author_name,
					(SELECT u.name FROM tasks.public.users u WHERE u.id = t.assigned_id) AS assigned_name,
					(SELECT string_agg(l.id ||'|'||l.name, ';' order by id) FROM
					tasks.public.tasks_labels tl ON t1.tasks_id = t.id
					JOIN tasks.public.labels l ON l.id = t1.label_id
					WHERE tl.task_id = t.id) AS labels
				   FROM tasks.public.tasks t 
			WHERE author_id = $1`
	return tasksByRows(s, sql, id)
}

func (s *Storage) TasksByLabel(id int) (dbflow.Tasks, error) {
	sql := `SELECT t.*,
			(SELECT u.name FROM tasks.public.users u WHERE u.id = t.author_id) AS author_name,
			(SELECT u.name FROM tasks.public.users u WHERE u.id = t.assigned_id) AS assigned_name,
			(SELECT string_agg(l.id ||'|'||l.name, ';' order by id) FROM
			 tasks.public.tasks_labels tl ON t1.tasks_id = t.id
			 JOIN tasks.public.labels l ON l.id = t1.label_id
			 WHERE tl.task_id = t.id) AS labels
			FROM tasks.public.tasks t
			JOIN tasks.public.tasks_labels tl on tl.task_id = t.id
			WHERE tl.label_id = $1`
	return tasksByRows(s, sql, id)
}

func (s *Storage) Tasks(idb int, ide int) (dbflow.Tasks, error) {
	sql := `SELECT t.*,
			(SELECT u.name FROM tasks.public.users u WHERE u.id = t.author_id) AS author_name,
			(SELECT u.name FROM tasks.public.users u WHERE u.id = t.assigned_id) AS assigned_name,
			(SELECT string_agg(l.id ||'|'||l.name, ';' order by id) FROM
			 tasks.public.tasks_labels tl ON t1.tasks_id = t.id
			 JOIN tasks.public.labels l ON l.id = t1.label_id
			 WHERE tl.task_id = t.id) AS labels
			FROM tasks.public.tasks t
			WHERE ($1 = 0 or t.id >= $1)
				  AND $2 = or t.id <= $2`
	return tasksByRows(s, sql, idb, ide)
}

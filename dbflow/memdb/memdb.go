package memdb

import (
	"fmt"
	"sort"
	"tasks/dbflow"
)

type Storage struct {
	id int
	db dbflow.Tasks
}

func New(connstr string) (*Storage, error) {
	db := make(dbflow.Tasks, 0, 10)
	return &Storage{
		id: 0,
		db: db,
	}, nil
}

func (s *Storage) NewTask(t dbflow.Task) (int, error) {
	s.id++
	t.SetID(s.id)
	s.db = append(s.db, t)
	sort.Slice(s.db, func(i, j int) bool { return s.db[i].ID() < s.db[j].ID() })
	return s.id, nil
}

func (s *Storage) UpdateTask(t dbflow.Task) error {
	for k, v := range s.db {
		if v.ID() == t.ID() {
			s.db[k] = t
			return nil
		}
	}

	return fmt.Errorf("Not found task id = %d", t.ID())
}

func (s *Storage) DeleteTask(id int) error {
	idk := -1
	for k, v := range s.db {
		if v.ID() == id {
			idk = k
			break
		}
	}
	if idk == -1 {
		return fmt.Errorf("Not fund task id = %d", id)
	}

	s.db = append(s.db[:idk], s.db[idk+1:]...)
	return nil
}

func (s *Storage) TasksByAuthor(id int) (dbflow.Tasks, error) {
	sl := make(dbflow.Tasks, 0, 10)

	for _, v := range s.db {
		if v.Author.ID() == id {
			sl = append(sl, v)
		}
	}
	return sl, nil
}

func (s *Storage) TasksByLabel(id int) (dbflow.Tasks, error) {
	sl := make(dbflow.Tasks, 0, 10)

	for _, v := range s.db {
		if v.Assign.ID() == id {
			sl = append(sl, v)
		}
	}
	return sl, nil
}

func (s *Storage) Tasks(idb int, ide int) (dbflow.Tasks, error) {
	sl := make(dbflow.Tasks, 0, 10)

	for _, v := range s.db {
		if (idb == 0 || v.ID() >= idb) && (ide == 0 || v.ID() <= ide) {
			sl = append(sl, v)
		}
	}

	return sl, nil
}

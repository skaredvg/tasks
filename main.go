package main

import (
	"fmt"
	"log"
	"tasks/dbflow"
	"time"

	stg "tasks/dbflow/memdb"
	_ "tasks/dbflow/postgres"
)

var ti *dbflow.TaskInterface

func main() {
	pwd := ""
	connstr := fmt.Sprintf("postgres://postgres:%s@localhost/tasks", pwd)
	ti, _ := stg.New(connstr)

	for i := 1; i < 10; i++ {
		t := dbflow.Task{Title: "Задача " + string(i), Content: "Сложная задача " + string(i)}
		_, err := ti.NewTask(t)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	sl, _ := ti.Tasks(1, 5)
	log.Printf("sl1=%+v\n", sl)

	ti.DeleteTask(5)

	log.Println("--------------------------")
	sl, _ = ti.Tasks(1, 5)
	log.Printf("%+v", sl)

	log.Println("--------------------------")
	t := sl[len(sl)-2]
	t.SetOpened(int(time.Now().Unix()))
	ti.UpdateTask(t)
	sl, _ = ti.Tasks(1, 0)
	log.Printf("sl2=%+v", sl)
}

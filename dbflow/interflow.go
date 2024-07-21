package dbflow

type TaskInterface interface {
	NewTask(Task) (int, error)
	UpdateTask(Task) error
	DeleteTask(int) error
	TasksByAuthor(int) (Tasks, error)
	TasksByLabel(string) (Tasks, error)
	Tasks(int, int) (Tasks, error)
}

package data

import (
	"database/sql"

	"github.com/dddbliss/flowman/internal/config"
	"github.com/dddbliss/flowman/internal/models"
	_ "modernc.org/sqlite"
)

var Db *sql.DB

func Open() {
	conn, err := sql.Open("sqlite", "./data.db")
	if err != nil {
		panic("could not open database file.")
	}

	conn.Exec("CREATE TABLE IF NOT EXISTS tasks (id integer primary key autoincrement, title varchar(32) not null, description varchar(32) not null, status integer not null)")
	_, err = conn.Exec("INSERT INTO tasks (title, description, status) values('test 1', 'desc 2', 1)")
	if err != nil {
		panic(err)
	}

	Db = conn
}

func GetTasksByStatus(status config.Status) ([]models.Task, error) {
	var Task models.Task
	var Tasks []models.Task

	rows, err := Db.Query("SELECT id, title, description, status FROM tasks where status = ?", status)
	if err != nil {
		return nil, err
	}

	var id, state int
	var title, desc string
	for rows.Next() {
		if err := rows.Scan(&id, &title, &desc, &state); err != nil {
			return nil, err
		}

		Task = models.NewTask(id, config.Status(state), title, desc)
		Tasks = append(Tasks, Task)
	}

	rows.Close() //good habit to close

	return Tasks, nil
}

func CreateTask(t models.Task) (int, error) {
	stmt, err := Db.Prepare("INSERT INTO tasks(title, description, status) VALUES(?, ?, ?)")
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(t.Title(), t.Description(), t.Status())
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil

}

func UpdateStatus(id int, status config.Status) (bool, error) {
	stmt, err := Db.Prepare("update tasks set status=? where id=?")
	if err != nil {
		return false, err
	}

	res, err := stmt.Exec(status, id)
	if err != nil {
		return false, err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	if affect > 0 {
		return true, nil
	}

	return false, nil
}

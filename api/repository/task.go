package repository

import (
	"database/sql"
	model "flutter_task_app_server/api/model"
)

type TaskRepo struct {
	db *sql.DB
}

func NewTaskRepo(db *sql.DB) *TaskRepo {
	return &TaskRepo{
		db,
	}
}

func (repo *TaskRepo) NewTodo(todo model.CreateTaskBody) error {
	_, insertErrr := repo.db.Exec("INSERT INTO tasks(title,description,userId,timestamp)VALUES(?,?,?, now())", todo.Title, todo.Description, todo.UserId)
	return insertErrr
}

func (repo *TaskRepo) GetTodosByUserId(userId int) ([]model.Task, error) {
	listOfTodos := []model.Task{}
	stmt, prepErr := repo.db.Prepare("Select id, title, description, isCompleted, timestamp from tasks where userId=?")
	if prepErr != nil {
		return listOfTodos, prepErr
	}
	rows, queryErr := stmt.Query(userId)

	if queryErr != nil {
		return listOfTodos, queryErr
	}

	for rows.Next() {
		var task model.Task
		rows.Scan(&task.Id, &task.Title, &task.Description, &task.IsCompleted, &task.Timestamp)
		listOfTodos = append(listOfTodos, task)
	}
	// _, insertErrr := repo.db.Exec("INSERT INTO tasks(title,description,userId,created_at)VALUES(?,?,?, now())", todo.Title, todo.Description, todo.UserId)
	return listOfTodos, nil
}

func (repo *TaskRepo) UpdateTaskStatus(taskId, userId, status int) error {
	_, updateStatusErr := repo.db.Exec("Update tasks set isCompleted = ? Where userId = ? AND id = ?", status, userId, taskId)
	return updateStatusErr
}
func (repo *TaskRepo) UpdateTask(task model.UpdateTaskBody, userId, taskId int) error {
	_, updateStatusErr := repo.db.Exec("Update tasks set title = ?, description = ? Where userId = ? AND id = ?", task.Title, task.Description, userId, taskId)
	return updateStatusErr
}

// func (repo *TaskRepo) FindUserByEmail(email string) (model.LoginUserBody, error) {
// 	row := repo.db.QueryRow("Select id,email,password from users where email = ?", email)
// 	user := model.LoginUserBody{}
// 	selectErr := row.Scan(&user.Id, &user.Email, &user.Password)
// 	fmt.Println(selectErr)
// 	return user, selectErr
// }

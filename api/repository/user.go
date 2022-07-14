package repository

import (
	"database/sql"
	model "flutter_task_app_server/api/model"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db,
	}
}

func (repo *UserRepo) RegisterUser(user *model.User) error {
	_, insertErrr := repo.db.Exec("INSERT INTO users(email,name,password)VALUES(?,?,?)", user.Email, user.Name, user.Password)
	return insertErrr
}
func (repo *UserRepo) FindUserByEmail(email string) (model.LoginUserBody, error) {
	row := repo.db.QueryRow("Select id,email,password from users where email = ?", email)
	user := model.LoginUserBody{}
	selectErr := row.Scan(&user.Id, &user.Email, &user.Password)
	return user, selectErr
}
func (repo *UserRepo) DeleteTaskByIdAndUserId(taskId, userId int) error {
	_, deleteErr := repo.db.Exec("Delete from tasks where id=? AND userId=?", taskId, userId)
	return deleteErr
}

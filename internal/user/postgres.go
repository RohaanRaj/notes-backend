package user

import (
	"fmt"
	"database/sql"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository{
	return &PostgresRepository{
		db: db,
	}
}


func(r *PostgresRepository) DeleteUser(Id string) error{
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		`
		delete from user_notes
		where user_id = $1
		`,Id, 
	)
	if err != nil{
		return err
	}
	_, err = tx.Exec(
		`
		delete from users_table
		where user_id = $1
		`, Id,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}


func(r *PostgresRepository) UpdateNotes(updatedNotes Notes, Id string) error{

	_, err := r.db.Exec(
		`
		update user_notes
		set (title, myspace) = ($1, $2)
		where user_id = $3
		`,updatedNotes.Title, updatedNotes.MySpace, Id, 
	)
	if err != nil {
		return err
	}
	return nil
}
func(r *PostgresRepository) DeleteNotes(Id string) error{

	_, err := r.db.Exec(
		`
		delete from user_notes
		where user_id = $1
		`, Id,
	)
	if err != nil{
		return err
	}
	return nil

}

func (r *PostgresRepository) GetNotes(Id string) ([]Notes, error){
	rows, err := r.db.Query(
		`
		select title, myspace
		from user_notes
		where user_id = $1
		order by created_at desc
		limit 10
		offset 0
		`,Id,
	)	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []Notes
	for rows.Next(){
		var n Notes
		if err := rows.Scan(&n.Title, &n.MySpace); err != nil{
			return nil, err
		}
		notes = append(notes, n)
	}
	if len(notes)>0{
		return notes, nil
	}
	return nil, sql.ErrNoRows
}

func(r *PostgresRepository) MyNotes(Id string, notes Notes) error{
	_, err := r.db.Exec(
		`
		insert into user_notes(user_id, title, myspace)
		values($1, $2, $3)
		`,Id, notes.Title, notes.MySpace,
	)
	if err != nil {
		return err 
	}
	return nil
}

func (r *PostgresRepository) SaveRefreshToken(email, refeshToken string) error{
	_, err := r.db.Exec(
		`
		insert into refreshtoken(email, refreshToken)
		values ($1, $2)
		`, email, refeshToken,
	)
	if err!= nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) GetRefreshToken(email string) (string, error) {
	var refreshToken string
	err := r.db.QueryRow(
		`
		select refreshToken
		from refreshtoken
		where email = $1
		`, email,
	).Scan(&refreshToken)
	if err != nil {
		return "", err 
	}
	return refreshToken, nil
}

func (r *PostgresRepository) GetByEmail(email string) (*User, error){

	var u User
	err := r.db.QueryRow(
		`
		select user_id, email, password
		from users_table 
		where email = $1
		`,email,
	).Scan(&u.Id, &u.Email, &u.Password)

	if err != nil {
		fmt.Println(err)
		return &User{}, err 
	}
	return &u, nil 
}

func (r *PostgresRepository) Create(email, password string) error{
	_, err := r.db.Exec(
		`
		insert into users_table (email, password)
		values ($1, $2)
		`,email, password,
	)
	if err != nil {
		fmt.Println(err)
		return err 
	}
	return nil
}


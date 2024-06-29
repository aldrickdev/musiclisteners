package db

import (
	"fmt"

	"github.com/aldricdev/musiclisteners/internals/types"
	"github.com/jmoiron/sqlx"
)

type GetAllUsersQueryResult struct {
	Users []types.User
	Err   error
}

type GetAllUsersQuery struct {
	SQL    []string
	Result chan GetAllUsersQueryResult
}

func NewGetAllUsersQuery(result chan GetAllUsersQueryResult) *GetAllUsersQuery {
	return &GetAllUsersQuery{
		SQL:    []string{SelectAllUsers},
		Result: result,
	}
}

func (q *GetAllUsersQuery) GetQuery() []string {
	return q.SQL
}

func (q *GetAllUsersQuery) Execute(dbConnection *sqlx.DB) {
	allUsers := []types.User{}
	singleUser := types.User{}

	// rows, err := dbConnection.Queryx(SelectAllUsers)
	rows, err := dbConnection.Queryx(q.SQL[0])
	if err != nil {
		q.Result <- GetAllUsersQueryResult{
			Users: allUsers,
			Err:   fmt.Errorf("Failed to query for all users: %q", err),
		}
		return
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.StructScan(&singleUser)
		if err != nil {
			q.Result <- GetAllUsersQueryResult{
				Users: allUsers,
				Err:   fmt.Errorf("Failed to scan for all users: %q", err),
			}
			return
		}

		allUsers = append(allUsers, singleUser)
	}
	q.Result <- GetAllUsersQueryResult{
		Users: allUsers,
		Err:   nil,
	}
}


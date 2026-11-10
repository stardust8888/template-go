package repo

import (
	"context"

	"template-go/db/mysql"
	"template-go/model"
	"template-go/utils/x/slicex"
)

func GetUsersByID(ctx context.Context, usersID []int) (map[int]*model.User, error) {
	users := make([]*model.User, 0)
	usersID = slicex.RemoveDuplicates(usersID)
	err := mysql.Find(ctx, mysql.GetUserDB(), &users, map[string]any{"id": usersID})
	if err != nil {
		return nil, err
	}
	result := make(map[int]*model.User, len(users))
	for i := range users {
		result[users[i].ID] = users[i]
	}
	return result, nil
}

package service

type mockDAO struct {
	users []int
}

func (d *mockDAO) CreateUser() error {
	d.users = append(d.users, 1)
	return nil
}

func (d *mockDAO) CreatePlayer() error {
	return nil
}

func (d *mockDAO) CreateTask() error {
	return nil
}

func (d *mockDAO) GetTask(id string) Task {
	return nil
}

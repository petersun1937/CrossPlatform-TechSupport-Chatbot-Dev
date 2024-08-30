package database

// type MockDB struct {
// 	mock.Mock
// }

// func (m *MockDB) Create(value interface{}) error {
// 	args := m.Called(value)
// 	return args.Error(0)
// }

// func (m *MockDB) Where(query interface{}, args ...interface{}) Database {
// 	//argsCalled := m.Called(query, args)
// 	//return argsCalled.Get(0).(Database)
// 	m.Called(query, args)
// 	return m
// }

// // Mock First method (for both User and Item types)
// func (m *MockDB) First(out interface{}, where ...interface{}) error {
// 	args := m.Called(out)
// 	switch v := out.(type) {
// 	case *models.User:
// 		*v = models.User{
// 			Model:        gorm.Model{ID: 1},
// 			UserID:       "12345", // Mock UserID
// 			FirstName:    "Peter",
// 			LastName:     "Sun",
// 			UserName:     "testuser",
// 			LanguageCode: "en",
// 			//Password: "testuserpassword",
// 			//Role: "user",
// 		}
// 		/*case *models.Item:
// 		*v = models.Item{
// 			Model:    gorm.Model{ID: 1},
// 			Title:    "testitem",
// 			Category: "testcategory",
// 		}*/
// 	}
// 	return args.Error(0)
// }

// func (m *MockDB) Save(value interface{}) error {
// 	args := m.Called(value)
// 	return args.Error(0)
// }

// // Mock the Model method
// func (m *MockDB) Model(value interface{}) Database {
// 	m.Called(value)
// 	return m
// }

// // Mock the Take method
// func (m *MockDB) Take(out interface{}, where ...interface{}) error {
// 	args := m.Called(out)
// 	return args.Error(0)

// 	/*args := m.Called(out)
// 	user := out.(*models.User)
// 	user.ID = 1
// 	user.Username = "testuser"
// 	user.Password = "testuserpassword"
// 	user.Role = "user"
// 	return args.Error(0)*/
// }

// func (m *MockDB) Delete(value interface{}, where ...interface{}) error {
// 	args := m.Called(value, where)
// 	return args.Error(0)
// }

// func (m *MockDB) Find(out interface{}, where ...interface{}) error {
// 	args := m.Called(out)
// 	return args.Error(0)
// }

// func (m *MockDB) Updates(values interface{}) error {
// 	args := m.Called(values)
// 	return args.Error(0)
// }

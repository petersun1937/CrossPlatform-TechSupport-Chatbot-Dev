package database

// type Database2 interface {
// 	Init() error
// 	GetDB() *gorm.DB
// }

// type database2 struct {
// 	user string
// 	pwd  string
// 	db   *gorm.DB
// }

// func NewDatabase2(config *config.Config) Database2 {
// 	/*return &database2{
// 		user: config.GetDBUser(),
// 		pwd:  config.GetDBPwd(),
// 	}*/
// }

// func (db2 *database2) Init() error {
// 	dbstr := fmt.Sprintf("host=localhost user=%s password=%s dbname=chatbot port=5432 sslmode=disable", db2.user, db2.pwd)
// 	db, err := gorm.Open(postgres.Open(dbstr), &gorm.Config{})
// 	if err != nil {
// 		return err
// 	}

// 	// Auto migrate the User and Item schemas
// 	// Write in two use for array
// 	if err := db.AutoMigrate(&models.User{}); err != nil {
// 		return err
// 	}

// 	db2.db = db
// 	return nil
// }

// func (db2 *database2) GetDB() *gorm.DB {
// 	return db2.db
// }

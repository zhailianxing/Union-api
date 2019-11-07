package media

import (
	"union/realapi/packet"
	"union/realapi/share"
)

//const (
//	INSERT_BOOK_SQL            = "INSERT INTO book (name) VALUES(?)"
//	SELECT_VALID_TOKEN_SQL       = "SELECT * FROM book WHERE id = ? "
//	DELETE_BOOK_BY_ID_SQL = "UPDATE book SET name = ? WHERE userid = ? LIMIT 1"
//)
//
//func AddBook(bookName string) (bool, error) {
//	_, err := share.TestDb.Exec(INSERT_BOOK_SQL, bookName)
//	if err!=nil {
//		return true, nil
//	}
//	return false, err
//}
//
//func DeleteBook(id int32, bookName string) (bool, error) {
//	_, err := share.TestDb.Exec(DELETE_BOOK_BY_ID_SQL, bookName, id)
//	if err!=nil {
//		return true, nil
//	}
//	return false, err
//}
//
//func GetBook(id uint32) (*packet.Book, error) {
//	//var ret *sql.Rows
//	//var err error
//	ret, err := share.TestDb.Query(SELECT_VALID_TOKEN_SQL, id)
//	if err != nil {
//		return nil, err
//	}
//	book := &packet.Book{}
//	if ret.Next() {
//		ret.Scan(&book.Id, &book.Name)
//	}
//	return book, nil
//}


const (
	SELECT_ALL_MEDIA_SQL       = "SELECT * FROM  media"
	INSERT_ONE_MEDIA_SQL       = "INSERT INTO media(name,status,user_id) values(?,?,101) "
)

//func GetBook(id uint32) (*packet.Book, error) {
//	//var ret *sql.Rows
//	//var err error
//	ret, err := share.TestDb.Quer
//	y(SELECT_VALID_TOKEN_SQL, id)
//	if err != nil {
//		return nil, err
//	}
//	book := &packet.Book{}
//	if ret.Next() {
//		ret.Scan(&book.Id, &book.Name)
//	}
//	return book, nil
//}

func GetMedia() ([]*packet.Meida, error) {
	ret, err := share.TestDb.Query(SELECT_ALL_MEDIA_SQL)
	if err != nil {
		return nil, err
	}
	media := make([]*packet.Meida, 0)

	for ret.Next(){
		item := new(packet.Meida)
		ret.Scan(&item.Id, &item.Name, &item.Status, &item.UserId, &item.Type, &item.Config, &item.ExpConfig,
			&item.Experiment, &item.Comment, &item.CreateTime, &item.ModifiedTime, &item.MediaClass, &item.DownloadUrl,
			&item.CategoryIds, &item.Rate, &item.IsRate, &item.CityClassBlackWhite, &item.AccClass, &item.BlackClass,
			&item.ApkId, &item.DomainId, &item.WhiteCityMaterialLevel, &item.MediaBdRateMax, &item.MediaMdRateMin)
		media = append(media, item)
	}
	return media, nil
}


func AddMedia(name string, status int32) (int32, error) {
	_, err := share.TestDb.Query(INSERT_ONE_MEDIA_SQL, name, status)
	if err != nil {
		return 0, err
	}
	return 1, nil
}



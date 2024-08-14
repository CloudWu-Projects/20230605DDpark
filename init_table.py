from utils.sqlite_dao import SqliteUtil
from conf.config import carplate_tabel

sql = """
CREATE TABLE PICInfo (
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	create_time DATETIME DEFAULT CURRENT_TIMESTAMP,	
	orderid TEXT NOT NULL,
	picUrl TEXT NOT NULL,	
	UNIQUE(orderid)
);



"""

dao = SqliteUtil()
def init_table():    
    dao.updateBySql(sql)
        

if __name__ == "__main__":
    init_table()
    print("ok!")

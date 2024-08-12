from utils.sqlite_dao import SqliteUtil
from conf.config import carplate_tabel

sql = """
CREATE TABLE IF NOT EXISTS TABLE_NAME (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    carnumber TEXT NOT NULL,
    vin TEXT NOT NULL,
    motor TEXT,
    pfjd TEXT ,
    oss TEXT ,
    url TEXT ,
    UNIQUE(carnumber),
    UNIQUE(vin)
);

"""
p_sql = """
CREATE TABLE park (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    park_id TEXT NOT NULL,
    campany TEXT ,
    UNIQUE(park_id)
);
"""

dao = SqliteUtil()
def init_table():
    for k in carplate_tabel:
        dao.updateBySql(sql.replace("TABLE_NAME",carplate_tabel[k]))
        
    dao.updateBySql(p_sql)    
if __name__ == "__main__":
    init_table()
    print("ok!")

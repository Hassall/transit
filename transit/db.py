import psycopg2

class Database:
    def __init__(self, host, user, password, database):
        self.conn = Database._connect(host, user, password, database)
        self.cur = self.conn.cursor()

    def query(self, query):
        self.cur.execute(query)

    def results(self):
        return self.cur.fetchall()

    def _connect(host, user, password, database):
        return psycopg2.connect(dbname=database, user=user, password=password, host=host)

    def close(self):
        self.cur.close()
        self.conn.close()
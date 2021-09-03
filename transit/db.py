import psycopg2

class Database:
    def __init__(self, host, user, password, database):
        self.conn = Database._connect(host, user, password, database)
        self.cur = self.conn.cursor()

    # make query
    def query(self, query):
        self.cur.execute(query)

    # collect all results of a query (throws exception if no rows)
    def results(self):
        return self.cur.fetchall()

    # connect to db
    def _connect(host, user, password, database):
        return psycopg2.connect(dbname=database, user=user, password=password, host=host)

    # close cursor and db connection
    def close(self):
        self.cur.close()
        self.conn.close()
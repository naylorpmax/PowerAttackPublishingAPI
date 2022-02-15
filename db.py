import os

import pandas as pd
from pandas import read_csv
from sqlalchemy import create_engine


class DB:
	def __init__(self, host: str, port: str, database: str, user: str, password: str):
		self.engine = create_engine(f"postgresql://{user}:{password}@{host}:{port}/{database}")

	@classmethod
	def from_env(cls):
		host = os.environ.get("POSTGRES_HOST", "homebrew-db")
		port = os.environ.get("POSTGRES_PORT", "5432")
		database = os.environ.get("POSTGRES_DB", "homebrew")
		user = os.environ.get("POSTGRES_USER", "max")
		password = os.environ["POSTGRES_PASSWORD"]
		return cls(host, port, database, user, password)


def main():
	db = DB.from_env()

	data_path = os.environ["DATA_PATH"]
	data = read_csv(data_path)

	table_name = os.environ["TABLE_NAME"]
	data.to_sql(
	    name=table_name,
	    con=db.engine,
	    if_exists="replace"
	)
	print(f"wrote data at: {data_path} to table: {table_name}")


main()

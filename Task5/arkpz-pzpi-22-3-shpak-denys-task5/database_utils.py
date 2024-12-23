import psycopg2
import subprocess
import os

# Define the `ensure_database_exists` function
# This function ensures that the specified database exists
# If the database does not exist, it creates a new database
def ensure_database_exists(dbname, user, password, host, port):
    try:
        default_connection = psycopg2.connect(
            dbname="postgres", user=user, password=password, host=host, port=port
        )
        default_connection.autocommit = True

        with default_connection.cursor() as cursor:
            cursor.execute("SELECT 1 FROM pg_database WHERE datname = %s", (dbname,))
            exists = cursor.fetchone()
            
            if not exists:
                cursor.execute(f"CREATE DATABASE \"{dbname}\"")
                print(f"База даних {dbname} створена успішно.")
            else:
                print(f"База даних {dbname} вже існує.")
    except psycopg2.Error as e:
        raise Exception(f"Помилка перевірки/створення бази даних: {e}")
    finally:
        if 'default_connection' in locals():
            default_connection.close()

# Define the `truncate_tables` function
# This function truncates all tables in the specified database
def truncate_tables(db_name, db_user, db_password, host, port):
    truncate_cmd = [
        "psql",
        f"--dbname={db_name}",
        f"--username={db_user}",
        f"--host={host}",
        f"--port={port}",
        "-c",
        "TRUNCATE TABLE roles, users, companies, routes, deliveries, product_categories, products, waypoints, sensor_data, user_companies RESTART IDENTITY CASCADE;",
    ]
    env = os.environ.copy()
    env["PGPASSWORD"] = db_password
    subprocess.run(truncate_cmd, env=env, check=True)

# Define the `restore_database` function
# This function restores the specified database from a backup file
# It uses the `pg_restore` command to restore the database
def restore_database(latest_file, db_name, db_user, db_password, host, port):
    restore_cmd = [
        "pg_restore",
        "--no-owner",
        "--role=postgres",
        f"--dbname={db_name}",
        "--format=c",
        "-v",
        "--clean",
        "--if-exists",
        f"--host={host}",
        f"--port={port}",
        f"--username={db_user}",
        latest_file,
    ]
    env = os.environ.copy()
    env["PGPASSWORD"] = db_password
    subprocess.run(restore_cmd, env=env, check=True)

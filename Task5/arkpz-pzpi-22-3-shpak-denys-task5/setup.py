import getpass
from config import write_config
from user_management import create_user, fetch_user_data
from system_utils import install_dependencies, add_go_path
from database_utils import ensure_database_exists, truncate_tables, restore_database
from encryption_utils import decrypt_file, encrypt_data, decrypt_data
from pathlib import Path
from helper import write_file, get_latest_backup_file

# Function to run the setup script
# This function will be called when the script is run
def main():
    print("=== Запуск скрипта налаштування ===")

    print("1. Створити нового користувача")
    print("2. Ввести дані користувача")

    choice = input("Виберіть дію: ")

    username = input("Введіть ім'я користувача: ").strip()
    password = getpass.getpass("Введіть пароль: ").strip()

    try:
        if choice == "1":            
            user_data = create_user(username, password)
        elif choice == "2":
            user_data = fetch_user_data(username, password)
        else:
            print("Невірний вибір.")
            return

        config_path = "./server/config/config.yaml"
        write_config(config_path, user_data)
        print(f"Конфігурація записана в {config_path}")

        install_dependencies(user_data['db_password'])

        ensure_database_exists(
            dbname=user_data['dbname'],
            user=user_data['dbuser'],
            password=user_data['db_password'],
            host="localhost",
            port=5432
        )

        migrations_folder = "./server/migrations"
        backup_path = get_latest_backup_file(migrations_folder)
        if not backup_path:
            return
            
        print(f"Знайдено файл міграції: {backup_path}")

        decrypted_sql = decrypt_file(backup_path, user_data['encryption_key'].encode())
     
        write_file(backup_path, decrypted_sql)

        truncate_tables(
            db_name=user_data['dbname'],
            db_user=user_data['dbuser'],
            db_password=user_data['db_password'],
            host="localhost",
            port=5432
        )

        restore_database(
            latest_file=backup_path,
            db_name=user_data['dbname'],
            db_user=user_data['dbuser'],
            db_password=user_data['db_password'],
            host="localhost",
            port=5432
        )

        encrypted_data = encrypt_data(decrypted_sql, user_data['encryption_key'])
        write_file(backup_path, encrypted_data)
         
        print("Міграція успішно виконана.")

    except Exception as e:
        print(f"Помилка: {e}")

if __name__ == "__main__":
    main()
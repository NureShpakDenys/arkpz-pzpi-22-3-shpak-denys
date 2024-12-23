import os

# Define the `write_file` function
# The function writes the data to the file
def write_file(path, data):
    with open(path, "w") as file:
        file.write(data)

# Define the `get_latest_backup_file` function
# The function returns the latest backup file
# from the specified folder
def get_latest_backup_file(migrations_folder):
    files = [
        os.path.join(migrations_folder, f) for f in os.listdir(migrations_folder)
        if os.path.isfile(os.path.join(migrations_folder, f))
    ]
    if not files:
        print("Нема файлів міграції в папці.")
        return False
    
    return max(files, key=os.path.getmtime)

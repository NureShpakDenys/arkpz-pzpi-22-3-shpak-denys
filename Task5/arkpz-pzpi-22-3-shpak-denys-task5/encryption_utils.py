from cryptography.fernet import Fernet

# Function to decrypt a file
# The function decrypts the file using the encryption key
# and returns the decrypted data
# The function raises an exception if an error occurs
# The function takes the file path and the encryption key as arguments
def decrypt_file(file_path, encryption_key):
    try:
        with open(file_path, 'rb') as encrypted_file:
            encrypted_data = encrypted_file.read()

        cipher = Fernet(encryption_key)
        decrypted_data = cipher.decrypt(encrypted_data)
        return decrypted_data.decode('utf-8')
    except Exception as e:
        raise Exception(f"Помилка розшифрування файла: {e}")

# Function to encrypt data
# The function encrypts the data using the encryption key
# and returns the encrypted data
# The function takes the data and the encryption key as arguments
def encrypt_data(data, encryption_key):
    fernet = Fernet(encryption_key)
    return fernet.encrypt(data.encode()).decode()

# Function to decrypt data
# The function decrypts the data using the encryption key
# and returns the decrypted data
# The function takes the data and the encryption key as arguments
def decrypt_data(data, encryption_key):
    fernet = Fernet(encryption_key)
    return fernet.decrypt(data.encode()).decode()

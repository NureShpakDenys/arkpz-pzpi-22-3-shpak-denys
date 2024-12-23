import json
import yaml
import requests
import bcrypt
import copy

# The user enters their username and password
def get_user_credentials():
    username = input("Enter your username: ")
    password = input("Enter your password: ")
    return username, password

# Hash the password
def hash_password(password):
    hashed_password = bcrypt.hashpw(password.encode('utf-8'), bcrypt.gensalt())
    return hashed_password.decode('utf-8')

# Convert YAML data to JSON
def convert_yaml_to_json(yaml_file_path):
    
    with open(yaml_file_path, 'r') as yaml_file:
        config_data = yaml.safe_load(yaml_file)
    
    return config_data

# Send JSON data to the server
def send_json_to_server(json_data, server_url):
    headers = {'Content-Type': 'application/json'}
    
    cleaned_data = copy.deepcopy(json_data)
    print(cleaned_data)
    
    response = requests.post(server_url, json=cleaned_data, headers=headers)

    if response.status_code == 200 or response.status_code == 201:
        print("Config data successfully sent to the server.")
    else:
        print(f"Failed to send data: {response.status_code} - {response.text}")

# Main function
def main():
    username, password = get_user_credentials()
    
    hashed_password = hash_password(password)
    
    config_yaml_path = ".\\server\\config\\config.yaml"  
    config_data = convert_yaml_to_json(config_yaml_path)
    
    new_user = {
        "username": username,             
        "password": hashed_password,      
        "configuration": config_data  
    }

    server_url = "http://localhost:8080/save-json"  
    send_json_to_server(new_user, server_url)


if __name__ == "__main__":
    main()

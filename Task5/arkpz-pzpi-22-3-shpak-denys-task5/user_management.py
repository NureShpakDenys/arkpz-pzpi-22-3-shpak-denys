import requests

# Function to create a new user
def create_user(username, password):
    url = "http://localhost:8080/create"
    response = requests.post(url, json={"username": username, "password": password})

    if response.status_code == 201:
        return response.json()
    else:
        raise Exception(f"Помилка створення користувача: {response.text}")

# Function to fetch user data
def fetch_user_data(username, password):
    url = f"http://localhost:8080/find?username={username}&password={password}"
    response = requests.get(url)

    if response.status_code == 200:
        return response.json()
    else:
        raise Exception(f"Помилка аутентифікації: {response.text}")

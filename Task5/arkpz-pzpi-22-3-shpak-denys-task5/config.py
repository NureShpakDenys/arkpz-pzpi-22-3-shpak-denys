import os
import yaml

# Define the `write_config` function
# The function writes the configuration data to the file
def write_config(config_path, config_data):
    os.makedirs(os.path.dirname(config_path), exist_ok=True)
    with open(config_path, 'w') as file:
        yaml.dump(config_data, file)
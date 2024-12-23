import os
import subprocess
import platform
import time

# Define the `is_program_installed` function
def is_program_installed(command, check_string):
    try:
        result = subprocess.check_output(command, shell=True, stderr=subprocess.STDOUT, text=True)
        return check_string in result
    except subprocess.CalledProcessError:
        return False
    except FileNotFoundError:
        return False

# Define the `add_go_path` function
# This function adds the Go path to the PATH environment variable
def add_go_path():
    go_path = os.path.join(os.environ['USERPROFILE'], 'go', 'bin')

    current_path = os.environ['PATH']

    if go_path not in current_path:
        new_path = f"{current_path};{go_path}"
        
        subprocess.run(["setx", "PATH", new_path], check=True)
        print(f"Go path ({go_path}) has been added to PATH successfully.")
    else:
        print("Go path is already present in PATH.")

# Define the `install_dependencies` function
# This function installs the required dependencies for the setup script
# It installs Go, PostgreSQL, and Taskfile
def install_dependencies(db_password):
    try:
        os_type = platform.system()

        go_installed = is_program_installed("go version", "go version")
        postgres_installed = is_program_installed("psql --version", "psql")
        taskfile_installed = is_program_installed("task --version", "Task")

        if os_type == "Windows":
            choco_path = r"C:\\ProgramData\\chocolatey\\bin\\choco.exe"
            choco_installed =  os.path.exists(choco_path)
            
            if not choco_installed:
                print("Chocolatey not found. Installing Chocolatey...")
                subprocess.run([
                    "powershell", "-Command",
                    "Set-ExecutionPolicy Bypass -Scope Process -Force; "
                    "[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.SecurityProtocolType]::Tls12; "
                    "iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))"
                ], check=True)

                print("Chocolatey installation completed.")
                subprocess.run([
                    "powershell", "-Command",
                    "Start-Process powershell -ArgumentList '-NoExit; [System.Environment]::SetEnvironmentVariable(\"PATH\", $env:PATH + \";C:\\ProgramData\\chocolatey\\bin\")' -Verb RunAs"
                ], check=True)
            
                
            if choco_installed:
                print("Installing Go, PostgreSQL, and Taskfile using Chocolatey...")
                if not go_installed:
                    subprocess.run(["choco", "install", "golang", "-y"], check=True)

                if not postgres_installed:
                    subprocess.run(["choco", "install", "postgresql16", "--params", f"'/password:{db_password}'"], check=True)

                if not taskfile_installed:
                    subprocess.run(["choco", "install", "go-task", "-y"], check=True)
            else:
                raise Exception("Chocolatey installation failed.")


        go_path = os.path.join(os.environ['USERPROFILE'], 'go', 'bin') if os_type == "Windows" else os.path.expanduser('~/go/bin')
        taskfile_path = os.path.join(os.environ['USERPROFILE'], 'taskfile', 'bin') if os_type == "Windows" else "/usr/local/bin"

        if go_path not in os.environ['PATH']:
            add_go_path()

        if taskfile_path not in os.environ['PATH']:
            os.environ['PATH'] += os.pathsep + taskfile_path

        print("Dependencies installed successfully.")

    except Exception as e:
        raise Exception(f"Error installing dependencies: {e}")


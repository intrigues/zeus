# Zeus - A Devops Automation Tool
With this tool we are trying generalize and minimize devops reperating task while trying to encourage shared responsibility model acorss devloper teams.

# Configurations
- `RUN_ENV=PRODUCTION`: Enables optimization specific to the production environment (default: nil)
- `HOSTNAME=0.0.0.0`: To run server on internet routable address (default: localhost)
- `PORT=8080`: To change the application port (default: 8080)
- `CREATE_ADMIN_USER=TRUE`: Creates admin user on startup
- `DEFAULT_ADMIN_USERNAME`: Username for admin user (created only if - `CREATE_ADMIN_USER` env is set to `TRUE`) (default: admin)
- `DEFAULT_ADMIN_EMAIL`: Email for admin user (created only if - `CREATE_ADMIN_USER` env is set to `TRUE`) (default: admin@example.com)
- `DEFAULT_ADMIN_PASSWORD`: Password for admin user (created only if - `CREATE_ADMIN_USER` env is set to `TRUE`) (default: nil -> It will be printed in logs)
- `DATA_DIR=./data`: Persistent data storage (make sure this directory has required permission)

# Usage Instructions
- Add a template
- Create new automation
- Select template that you want to automate
- Provide git credentials and url to where you want to push your code
- Fetch branches by pressing on fetch button besides branch dropdown
- Select the branch
- Add the credentials
- Click on automate

## Creating templates
- Add project name and Technology name
- Add number of files required
    - Add filename
    - Add file content
        - Templatized files: File should contain @@VARIABLE_NAME@@ where you want dynamic varialbes
        - Add mapping for each dynamic variable

### Examples for files
File content example
```
def projectName = "@@PROJECT_NAME@@"
def serviceName = "@@SERVICE_NAME@@"
...
...
end of code
```

File variable mapping example
```
[
    {
        "variable": "PROJECT_NAME",
        "hint": "Project Name",
        "placeholder": "Project Name"
    },
    {
        "variable": "SERVICE_NAME",
        "hint": "Service Name",
        "placeholder": "Service Name"
    }
]
```

# Developer Guidelines
Install required dependancies using
```
go mod tidy
```

Run the application
```
./run.sh
```
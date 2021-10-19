# Zeus - A Devops Automation Tool
With this tool we are trying generalize and minimize devops reperating task while trying to encourage shared responsibility model acorss devloper teams.

# Configurations
`RUN_ENV=PRODUCTION`: Enables optimization specific to the production environment (default: nil)
`HOSTNAME=0.0.0.0`: To run server on internet routable address (default: localhost)
`PORT=8080`: To change the application port (default: 8080)
`CREATE_ADMIN_USER=TRUE`: Creates admin user on startup
`DEFAULT_ADMIN_USERNAME`: Username for admin user (created only if `CREATE_ADMIN_USER` env is set to `TRUE`) (default: admin)
`DEFAULT_ADMIN_EMAIL`: Email for admin user (created only if `CREATE_ADMIN_USER` env is set to `TRUE`) (default: admin@example.com)
`DEFAULT_ADMIN_PASSWORD`: Password for admin user (created only if `CREATE_ADMIN_USER` env is set to `TRUE`) (default: nil -> It will be printed in logs)

# Developer Guidelines
Install required dependancies using
```
go mod tidy
```

Run the application
```
./run.sh
```
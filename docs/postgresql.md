# PostgreSQL Setup Guide (Windows)

This guide explains how to install, configure, and run PostgreSQL on Windows using Chocolatey and the command line. It is intended for local development with the go-vibe microservice.

---

## 1. Install PostgreSQL with Chocolatey

Chocolatey is a popular Windows package manager. If you don't have it, install from https://chocolatey.org/install.

Open **PowerShell as Administrator** and run:

```powershell
choco install postgresql
```

This installs PostgreSQL binaries.
- The default install location is `C:\Program Files\PostgreSQL\<version>`.
- choco adds PostgreSQL CLI tools (like `psql` and `pg_ctl`) to your PATH. Simply open new PS after install.
- The installer adds a windows service and starts it as well. go to "services", find the service, stop it and set it's run type to "manual"
---


## 2. Initialize and Configure PostgreSQL

### Initialize the Database Cluster

If not already initialized, 
- change to the root folder of your project, 
- create a folder "data" 
- and run the following
- Note: encoding UTF8 is essential for go to go ;-)

```powershell
initdb -D ".\data" -U myapp -A password -W --encoding=UTF8 --locale=en_US.UTF-8
```

- Right at the beginning you are going to be asked about a password: for dev purpose it's safe to use ``myapp``
- This creates the data directory and sets up the initial database.

### Configure PostgreSQL to Listen on localhost:5432

Edit `postgresql.conf` (usually in `C:\Program Files\PostgreSQL\<version>\data`):

- Set:
  - `listen_addresses = 'localhost'`
  - `port = 5432`

You can edit with Notepad or PowerShell:

```powershell
notepad "C:\Program Files\PostgreSQL\<version>\data\postgresql.conf"
```

---

## 3. Start PostgreSQL Server

Use the `pg_ctl` command to start the server. 
Change into the folder of your project!

```powershell
pg_ctl -D ".\data" start
```

- The server will run on `localhost:5432`.

To stop the server:

```powershell
pg_ctl -D ".\data" stop
```

## 4. Just once - Create the myapp DB in the postgresql server

use psql to connect to a freshly setup db and just once create the DB "myapp"

```powershell
psql -U myapp  postgres
```
Type ``myapp`` when asked for the password

then at the prompt enter
```
postgres=# CREATE DATABASE myapp OWNER myapp;
```

After this verify that the new db exists
```powershell
psql -U myapp  myapp
```

Again type ``myapp``
Note: Now you get a different prompt
```
myapp=#
```

---

## 5. Connect to PostgreSQL

To connect to the project DB using the CLI:

```powershell
psql -U myapp -h localhost -p 5432
```

---

## 6. Configure your deployment stage for running the application locally

Don't forget to configure your application configuration.
Open File ``.\config\development.yaml`` (development is the default stage, see docu about selecting the stage elsewhere)
It should read
```
database:
  url: "postgres://myapp:myapp@localhost:5432/myapp?sslmode=disable"
```

---

## 7. Troubleshooting

- If port 5432 is in use, change the `port` setting in `postgresql.conf`.
- Ensure Windows Firewall allows connections to port 5432.
- For more help, see the official docs: https://www.postgresql.org/docs/

---

## Summary

- Install with Chocolatey
- Configure to listen on localhost:5432
- Start/stop with `pg_ctl`
- Connect with `psql`

This setup is ideal for local development with the go-vibe microservice.

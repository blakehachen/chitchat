# ChitChat

A clone of the forum concept with user tracking. Users have the ability to like posts, threads and view their interactions within a user dashboard private to them.

## DB Setup Instructions (Linux)

1. Navigate to command-line, install postgreSQL and initialize database Run `initdb -D /var/lib/postgres/data`
2. Run `psql -U postgres`
3. Next, we will need to create a role matching our user account Run `createuser --interactive` when prompted to enter the name of the role make sure it matches your local user account name.
4. After re-entering the postgres cli as your new user Run `createdb chitchat` then exit the cli interface. The database now exists in postgres, all that is left is to set the tables up with the setup.sql file.
5. In order to import the schema Run `psql -U <NEW_USER> -d chitchat -f ~/path/to/setup.sql` This will create a database schema and grant the program functionality.

## Building and serving to a local machine

1. Now that the database is built you can compile the go program. Navigate to the root of the repository and Run `go build && ./chitchat`

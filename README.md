# Gator, a RSS feed aggrator for the command line

Small CLI tool written in Go. Fetches RSS feeds. Multiple users can register, add feeds, subsribe to feeds, and browse recent posts from the feeds they follow.

This is heavily based on the guided Blog aggregator project on boot.dev.

## Installation

These instructions are for Linux (WSL/Ubuntu)

### 1. Install Go

gator requires Golang toolchain for the installation process.

**Option 1**: [The webi installer](https://webinstall.dev/golang/) is the simplest way for most people. Just run this in your terminal:

```bash
curl -sS https://webi.sh/golang | sh
```
*Read the output of the command and follow any instructions.*

**Option 2**: Use the [official installation instructions](https://go.dev/doc/install).

Run `go version` on your command line to make sure the installation worked. If it did, move on to step 2.

### 2. Install PostgreSQL

```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
```

Run `psql --version` to make sure the installation worked. psql is the default front-end (client) to PostgreSQL database (server).

### 3. Install/build gator

**Option 1**: Install gator using `go install`. This is the easiest way if you just want to use gator.

```bash
go install github.com/hemukka/gator@latest
```
This downloads, compiles and installs it automatically. After, you should be able to use `gator` anywhere.

**Option 2**: Download the source code and compile it with `go build`. This is good if you want to modify the source code.

Clone with git:
```bash
mkdir gator
git clone https://github.com/hemukka/gator.git
cd gator
```
or download ZIP from https://github.com/hemukka/gator and extract.

Build the binary:
```bash
go build
```
This creates an executable in the current directory. You run it with `./gator` instead of `gator` in all later commands.

## Setup

### 4. Create PostgreSQL database

For a local PostgreSQL server, start one with:
```bash
sudo service postgresql start
```
To stop it or check the status:
```bash
sudo service postgresql stop
sudo service postgresql status
```
The default admin user is `postgres`, you may need to set the password for this account with:
```bash
sudo passwd postgres
```

Once the server is up, access it with psql:
```bash
sudo -u postgres psql
```
Create a new database:
```sql
CREATE DATABASE gator;
```
Connect to the database:
```sql
\c gator
```

This database has a default user called postgres, set its password with:
```sql
ALTER USER postgres PASSWORD 'postgres';
```

Exit psql with `\q` or `exit` or Ctrl+D.

Now the connection string to this db should be this URL:
```
postgres://postgres:postgres@localhost:5432/gator
```

### 5. Perform database migrations

Required migrations are in the `sql/schema`directory.

**Option 1**: Input them manually using psql. Open gator database direclt with:
```bash
psql "postgres://postgres:postgres@localhost:5432/gator"
```
and apply all the up migrations.

**Option 2**: Use Goose database migration tool. Install it with Go:
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```
Run the up migrations:
```bash
cd sql/schema
goose postgres <connection_string> up
cd ../..
```
You can use `\dt` in psql to list all tables.

### 6. Create a config file

Create a `.gatorconfig.json` file in your home directory and put the connection string as "db_url":
```bash
echo "{"db_url":"postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"}" > ~/.gatorconfig.json
```

### 7. Run the aggregator

Start the long-running aggregator that scrapes the RSS feeds regularly:
```bash
gator agg <time_between_requests>
```
The time can be for example `1h`, `30m`, `40s` or `10m30s`.

Leave this running in one terminal.


## Usage
### Users

Register as a user:
```bash
gator register <user_name>
```
This will also log you in. To log into a different user:
```bash
gator login <user_name>
```
List all users and which one is current user:
```bash
gator users
```

### Feeds

Add RSS feed to the aggregator:
```bash
gator addfeed <feed_name> <feed_url>
```
this will also make the current user follow the feed.

List all feeds added to gator:
```bash
gator feeds
```

### Follows

List all feeds followed by the current user:
```bash
gator following
```

Follow a feed a that has already been added to gator:
```bash
gator follow <feed_url>
```

Unfollow a feed:
```bash
gator unfollow <feed_url>
```

### Browse

Browse most recent posts from the feeds followed by the current user:
```bash
gator browse <post_limit>
```
post_limit specifies how many posts are displayed, default is 2.


## Reset

Reset all users with `gator reset`. This permanently deletes all users, feeds, follows, and posts.
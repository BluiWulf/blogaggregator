
# RSS Aggregator Project
---
Welcome to my RSS Aggregator for the *Boot.dev* guided project, ***Blog Aggregator***.

This project requires *Postgres* and *Go* installed in order to use the *rssagg* CLI.

**NOTE:** It is assumed you are using this tool on a Linux / WSL (Debian) system.  All installation instructions are written as such.

## Installation and Configuration

### Install PostgreSQL

[PostgreSQL](https://www.postgresql.org/) (or Postgres) is production-ready, open-source database tool operating as a SQL server.  It listens for requests and responds accordingly.  To interact with Postgres will require a client.  For this project, we will use [psql](https://www.postgresql.org/docs/current/app-psql.html#:~:text=psql%20is%20a%20terminal%2Dbased,or%20from%20command%20line%20arguments.) as a client.

1. Install Postgres

    ```bash
    sudo apt update
    sudo apt install postgresql postgresql-contrib
    ```

2. Verify installation

    ```bash
    psql --version
    ```

3. Update postgres password

    ```bash
    sudo passwd postgres
    ```
    - Recommend an easy password to remember just for this project, such as `postgres`

4. Start Postgres server

    ```bash
    sudo service postgresql start
    ```

5. Enter `psql` shell

    ```bash
    sudo -u postgres psql
    ```

    - Should see the following prompt

    ```sql
    postgres=#
    ```

6. Create new database named `gator` at

    ```sql
    postgres=# CREATE DATABASE gator;
    ```

7. Connect to new `gator` database

    ```sql
    postgres=# \c gator
    ```

    - Should see the following prompt

    ```sql
    gator=#
    ```

8. Set the user password

    ```sql
    gator=# ALTER USER postgres PASSWORD 'postgres';
    ```
    - Recommend an easy password to remember just for this project, such as `postgres`

9. Exit out of `psql`

    ```sql
    gator=# exit
    ```

10. Test the database connection URL

    - It will be in this format:

    ```bash
    protocol://username:password@host:port/database?sslmode=disable
    ```

    - This will be the connection URL for this project

    ```bash
    psql "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"
    ```

    - You should once again see the `gator` prompt
    - You can exit out of `psql` after seeing the prompt

    **NOTE:** This needs to be saved for later setting up a configuration file

### Install Go

[Go](https://go.dev/) is an open-source programming language supported by Google with a robust library.  It is the language used to build and run this RSS aggregator tool.

The [Webi Installer](https://webinstall.dev/golang/) is usually the simplest method for installing `Go` on Linux/WSL.

1. Enter the installation command in the terminal

    ```bash
    curl -sS https://webi.sh/golang | sh
    ```

    - Read the output of the command and follow any additional instructions
    - This same command can be used in order to update go

2. Verify the installation was successful

    ```bash
    go version
    ```

    #### Troubleshooting

    If you see a *"command not found"* error after installation, then the directory containing the `go` program is most likely missing from the system `PATH` environment variable.

    1. Local *where* the `go` program is located (possibly in the locations below)

        - `~/.local/opt/go/bin` (Webi Install)
        - `/usr/local/go/bin` (Official Golang Install)

    2. Run the following commands in the terminal

    ```bash
    echo `export PATH=$PATH:$HOME/.local/opt/go/bin` >> ~/.bashrc
    source ~/.bashrc
    ```

### Install `rssagg` Tool

Time to install the RSS Aggregator CLI tool.

1. Build and install the tool from the source code

    ```bash
    go install .
    ```

2. Make sure the '~/go/bin/' directory is added to `PATH`

    ```bash
    echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc && source ~/.bashrc
    ```

3. Add an alias for the tool

    ```bash
    echo 'alias rssagg="blogaggregator"' >> ~/.bashrc && source ~/.bashrc
    ```

    - This isn't necessary, but it makes for a simpler and shorter command usage with the tool
    - Feel free to use a different alias if you prefer or no alias at all

### Create Configuration File

A configuration file is needed to maintain the connection URL and the current user of the tool.  This will be stored in the users `HOME` directory where the tool will be expecting it.

1. Create the `~/.gatorconfig.json` configuration file in the `HOME` directory with the following content:

    ```json
    {
	    "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
	    "current_user_name": ""
    }
    ```

    - An empty string will be sufficient for the current username as this will be updated appropriately when the `register` and/or `login` commands are used

## Commands and Usage

This tool is designed to run directly in the terminal.  It can run in the background aggregating posts from feeds and storing them in the database for perusing later.  Another terminal can be used to browse the latest posts stored.

This section will provide a description of each available command as well as the input parameters required to use them.

### register

#### Description

The `register` command is used add a new user to the database.  It will also update the configuration file to make the new user the current user.

#### Arguments

***`username`***: User to be registered
- Must be unique
- Cannot contain whitespace characters

### login

#### Description

The `login` command is used to change which user is currently using the tool.  It will also update the configuration file to the username provided as an argument.

#### Arguments

***`username`***: User to be logged in
- Must be unique
- Cannot contain whitespace characters

### users

#### Description

The `users` command is used to list all registered users in the database.  It take no additional arguments.

### agg

#### Description

The `agg` command continuously runs until terminated scraping one feed at a time for all new posts and adding them to the database for viewing.

#### Arguments

***`interval`***: Time interval between scraping a single feed for new posts
- Must be a positive integer
- Must have an appropriate time unit (i.e., `s`, `m`, `h`)
- For example, `30s`

### addfeed

#### Description

The `addfeed` command registers a new feed for the current user.

#### Arguments

***`name`***: The name of the RSS feed
- Recommend using quotes around the name to allow for whitespace characters

***`url`***: The URL of the RSS feed
- Must be unique

### feeds

#### Description

The `feeds` command is used to list all registered feeds in the database.  It take no additional arguments.

### follow

#### Description

The `follow` command registers a new follow for the current user to a feed already registered in the database.  If the provided feed URL isn't registered, then it will return an error.

#### Arguments

***`url`***: The URL of the RSS feed
- Must be unique
- Must be a URL already registered in the database

### following

#### Description

The `following` command is used to list all feeds the current user is currently following.  It take no additional arguments.

### unfollow

#### Description

The `unfollow` command unfollows the provided feed URL for the current user.  If the provided feed URL isn't registered or it isn't currently followed by the current user, then it will return an error.

#### Arguments

***`url`***: The URL of the RSS feed
- Must be unique
- Must be a URL already registered in the database
- Must be a URL the current user is actively following

### browse

#### Description

The `browse` command lists a number of the latest posts for the current user's followed feeds.  It takes an optional parameter to define a specific number of posts to provide.

#### Arguments

***`count`***: The number of posts to display
- Defaults to `2` if not provided
- Must be a positive integer

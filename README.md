# Readcommend

Readcommend is a book recommendation web app for the true book aficionados and disavowed human-size bookworms. It allows to search for book recommendations with best ratings, based on different search criteria.

# Instructions

The front-end single-page app has already been developed using Node/TypeScript/React (see `/app`) and a PostgreSQL database with sample data is also provided (see `/migrate.*`). Your mission - if you accept it - is to implement a back-end microservice that will fetch the data from the database and serve it to the front-end app.

- In the `service` directory, write a back-end microservice, using the language of your choice, that listens on `http://localhost:5001`. **NOTE: You will be evaluated on the mastery of your chosen language, you should opt to use the language you are most proficient with**.
- Write multiple REST endpoints, all under the `/api/v1` route, as specified in the `/open-api.yaml` swagger file.
- The most important endpoint, `/books`, must return book search results in order of descending ratings (from 5.0 to 1.0 stars) and filtered according to zero, one or multiple user selected criteria: author(s), genre(s), min/max pages, start/end publication date (the "era"). A maximum number of results can also be specified.
- It's OK to use libraries for http handling/routing and SQL (ie: query builders), but try to refrain from relying heavily on end-to-end frameworks (ie: Django) and ORMs that handle everything and leave little room to showcase your coding skills! ;)
- Write some documentation (ie: at the end of this file) to explain how to deploy and run your service. *Please assume that the reviewer should be able to run your solution without needing to install extra dependencies on his machine.*
- Keep your code simple, clean and well-organized.
- During development commit and push regularly to the CodeSubmit git repository and, when you are done, submit your assignment in CodeSubmit.
- Don't hesitate to come back to us with any questions along the way. We prefer that you ask questions, rather than assuming and misinterpreting requirements.
- You have no time limit to complete this exercise, but the more time you take, the higher our expectations in terms of quality and completeness.
- You will be evaluated mainly based on how well you respect the above instructions. Please refer to the **Expectations** section below for more details on how you will be evaluated.

# Expectations

nesto recruits for all proficiency levels. Depending on your experience and the position you are applying for, what we expect to see in the test will differ. We do, however, understand that you may have a life (some people do). If you don't have the time to respect all the instructions, simply do your best and focus on what you deem most important.

**Junior Submission**

- Provides a solution that is bug free, fast and meets all the requirements listed above.
- Code provided shows a good proficiency in the chosen language.
- Shows good reflexes with regards to 3rd party dependencies. You should add dependencies that add value to your solution without solving everything for you.
- Basic testing of the core functionalities.

**Senior Submission**

We expect everything from a *Junior Submission* and that you treat this as production ready code without necessarily providing full test coverage
- Solution is architectured properly to allow for adequate test coverage and future extensions.
- Each layer of the provided solution must be protected against typical attack vectors. You should **NOT** rely on your API contract to protect you from malicious inputs.
- Showcases the type of tests needed to validate your architectural choices and to provide confidence that your solution satisfies the 'business' and performance requirements.

# Development environment

## Docker Desktop

Make sure you have the latest version of Docker Desktop installed, with sufficient memory allocated to it, otherwise you might run into errors such as:

```
app_1         | Killed
app_1         | npm ERR! code ELIFECYCLE
app_1         | npm ERR! errno 137.
```

If that happens, first try running the command again, but if it doesn't help, try increasing the amount of memory allocated to Docker in Preferences > Resources.

## Starting front-end app and database

In this repo's root dir, run this command to start the front-end app (on port 8080) and PostgreSQL database (on port 5432):

```bash
$ docker-compose up --build
$ docker-compose up -d
```

(later you can press Ctrl+C to stop this docker composition when you no longer need it)

Wait for everything to build and start properly.

## Connecting to database

During development, you can connect to and experiment with the PostgreSQL database by running this command:

```bash
$ ./psql.sh
```

Then, on the psql CLI, test as follows:

```psql
readcommend=# \dt
```

If everything went well, you should get this result:

```psql
    List of relations
 Schema |  Name  | Type  |  Owner
--------+--------+-------+----------
 public | author | table | postgres
 public | book   | table | postgres
 public | era    | table | postgres
 public | genre  | table | postgres
 public | size   | table | postgres
(5 rows)
```

To exit the PostgreSQL session, type `\q` and press `ENTER`.

## Accessing front-end app

Point your browser to http://localhost:8080

Be patient, the first time it might take up to 1 or 2 minutes for parcel to build and serve the front-end app.

You should see the front-end app appear, with all components displaying error messages because the back-end service does not exist yet.

# Deploying and running back-end microservice

First run should run the front-end service. After this, it should be located in the `./service` directory where all the backend service information is located. There you can configure your credentials for your postgres instance in the `config/config.json` file.

To run the service run
`$ go run main.go`

This will run the service on port `http://localhost:5001`, just make sure you have the database running on the postgres instance.

To run the tests, we can use the command:
`$ go test ./...`


---

## Debugging PostgreSQL in Docker

### 1. List Active Containers

If you're unsure of the exact name of your PostgreSQL container, list all active containers with:

```bash
docker ps
```

Look for the container corresponding to your PostgreSQL service, for example, something like `book-recommendations-postgresql-1`.

### 2. Access the PostgreSQL Container

Once you have the container name, access it using the following command:

```bash
docker exec -it <container_name> psql -U postgres -d readcommend
```

Replace `<container_name>` with the actual name of your PostgreSQL container. For example:

```bash
docker exec -it book-recommendations-postgresql-1 psql -U postgres -d readcommend
```

This command opens the `psql` client, connecting you to the `readcommend` database with the `postgres` user.

### 3. Execute Commands in PostgreSQL

While inside the `psql` console, you can run SQL commands to check the state of your database. Here are some useful commands:

- **List Databases**: To view all available databases, use:

  ```sql
  \l
  ```

  or

  ```sql
  \list
  ```

- **Connect to a Database**: To connect to a specific database, use:

  ```sql
  \c database_name
  ```

- **Show All Tables**: To list all tables in the current database, use:

  ```sql
  \dt
  ```

- **Show Table Structure**: To view the structure of a specific table, including column names and data types, use:

  ```sql
  \d table_name
  ```

- **View All Data in a Table**: To see all data from a table, use the following SQL command:

  ```sql
  SELECT * FROM table_name;
  ```

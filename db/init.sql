CREATE TABLE IF NOT EXISTS "user"
(
    id serial NOT NULL PRIMARY KEY,
    login varchar(16) NOT NULL UNIQUE,
    password text NOT NULL,
    is_admin boolean DEFAULT false
);

CREATE TABLE IF NOT EXISTS movie
(
    id serial NOT NULL PRIMARY KEY,
    name varchar(150) NOT NULL,
    description varchar(1000),
    release_date date NOT NULL,
    rating int,
    CHECK (rating >= 0 AND rating <= 10)
);

CREATE TABLE IF NOT EXISTS actor
(
    id serial NOT NULL PRIMARY KEY,
    name varchar(16) NOT NULL,
    surname varchar(16) NOT NULL,
    gender char(1) NOT NULL,
    birth_date date NOT NULL,
    CHECK ( gender in ('F', 'M') )
);

CREATE TABLE IF NOT EXISTS movie_actor
(
    id serial NOT NULL PRIMARY KEY,
    movie_id int,
    actor_id int,
    FOREIGN KEY (movie_id) REFERENCES movie(id) ON DELETE CASCADE,
    FOREIGN KEY (actor_id) REFERENCES actor(id) ON DELETE CASCADE
);

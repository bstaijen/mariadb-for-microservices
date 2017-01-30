CREATE SCHEMA ProfileService;
CREATE SCHEMA CommentService;
CREATE SCHEMA PhotoService;
CREATE SCHEMA VoteService;

CREATE TABLE IF NOT EXISTS ProfileService.users (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, username varchar(255) NOT NULL UNIQUE, email varchar(255) NOT NULL UNIQUE, password varchar(255) NOT NULL,createdAt timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, updatedAt timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP);

CREATE TABLE IF NOT EXISTS PhotoService.photos (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, title varchar(255) NOT NULL, user_id INT NOT NULL, filename varchar(255) NOT NULL UNIQUE, contentType varchar(255),photo MEDIUMBLOB,createdAt timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, updatedAt timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP);

CREATE TABLE IF NOT EXISTS VoteService.votes (user_id INT NOT NULL, photo_id INT NOT NULL,upvote boolean DEFAULT false,downvote boolean DEFAULT false ,createdAt timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, updatedAt timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, CONSTRAINT constraint_key PRIMARY KEY(user_id, photo_id) );

CREATE TABLE IF NOT EXISTS CommentService.comments (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, user_id INT NOT NULL, photo_id INT NOT NULL,comment TEXT NOT NULL,createdAt timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, updatedAt timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP);

CREATE USER 'authentication_service'@'%' IDENTIFIED BY 'password';
GRANT ALL ON ProfileService.* TO 'authentication_service'@'%';

CREATE USER 'profile_service'@'%' IDENTIFIED BY 'password';
GRANT ALL ON ProfileService.* TO 'profile_service'@'%';

CREATE USER 'comment_service'@'%' IDENTIFIED BY 'password';
GRANT ALL ON CommentService.* TO 'comment_service'@'%';

CREATE USER 'photo_service'@'%' IDENTIFIED BY 'password';
GRANT ALL ON PhotoService.* TO 'photo_service'@'%';

CREATE USER 'vote_service'@'%' IDENTIFIED BY 'password';
GRANT ALL ON VoteService.* TO 'vote_service'@'%';
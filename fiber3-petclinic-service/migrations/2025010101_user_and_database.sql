-- steps to create a new user and database for the pet clinic service
CREATE USER pet_user WITH PASSWORD 'pet_user';
CREATE DATABASE pet_clinic owner pet_user;
GRANT ALL PRIVILEGES ON DATABASE pet_clinic TO pet_user;




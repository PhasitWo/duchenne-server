
CREATE TABLE appointment (
  id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
  create_at bigint NOT NULL,
  date bigint NOT NULL,
  patient_id int NOT NULL,
  doctor_id int NOT NULL
);


CREATE TABLE question (
  id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
  topic varchar (50) NOT NULL,
  question varchar(700) NOT NULL,
  create_at bigint NOT NULL,
  answer varchar(500),
  answer_at bigint,
  patient_id int NOT NULL,
  doctor_id int 
);


CREATE TABLE doctor (
  id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
  first_name varchar(255) NOT NULL,
  middle_name varchar(255),
  last_name varchar(255) NOT NULL,
  username varchar(20) UNIQUE NOT NULL,
  password varchar(255) NOT NULL,
  role varchar(10) NOT NULL
);

CREATE TABLE patient (
  id int PRIMARY KEY AUTO_INCREMENT,
  hn varchar(15) UNIQUE NOT NULL,
  first_name varchar(255) NOT NULL,
  middle_name varchar(255),
  last_name varchar(255) NOT NULL,
  email varchar(255),
  phone varchar(15),
  verified bool NOT NULL DEFAULT 0
);

CREATE TABLE device (
  id int PRIMARY KEY AUTO_INCREMENT,
  login_at bigint NOT NULL,
  device_name varchar(30) NOT NULL,
  expo_token varchar(100) NOT NULL,
  patient_id int NOT NULL
);

ALTER TABLE appointment ADD CONSTRAINT appointment_doctor_id_fk FOREIGN KEY (doctor_id) REFERENCES doctor (id) ON DELETE CASCADE;
ALTER TABLE appointment ADD CONSTRAINT appointment_patient_id_fk FOREIGN KEY (patient_id) REFERENCES patient (id) ON DELETE CASCADE;
ALTER TABLE question ADD CONSTRAINT question_doctor_id_fk FOREIGN KEY (doctor_id) REFERENCES doctor (id) ON DELETE SET NULL;
ALTER TABLE question ADD CONSTRAINT question_patient_id_fk FOREIGN KEY (patient_id) REFERENCES patient (id) ON DELETE CASCADE;
ALTER TABLE device ADD CONSTRAINT device_patient_id_fk FOREIGN KEY (patient_id) REFERENCES patient (id) ON DELETE CASCADE;

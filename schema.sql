
CREATE TABLE appointment (
  id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
  create_at timestamp,
  date timestamp,
  patient_id int NOT NULL,
  doctor_id int NOT NULL
);


CREATE TABLE ask (
  id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
  topic varchar (30) NOT NULL,
  question varchar(500) NOT NULL,
  answer varchar(500),
  patient_id int NOT NULL,
  doctor_id int NOT NULL
);


CREATE TABLE doctor (
  id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
  doctor_info varchar(500),
  username varchar(20),
  password varchar(100),
  role varchar(10) NOT NULL
);

CREATE TABLE patient (
  id int PRIMARY KEY AUTO_INCREMENT,
  hn varchar(15) UNIQUE NOT NULL,
  first_name varchar(100) NOT NULL,
  middle_name varchar(100),
  last_name varchar(100) NOT NULL,
  email varchar(255),
  phone varchar(15),
  verified bool NOT NULL DEFAULT 0,
);

ALTER TABLE appointment ADD CONSTRAINT appointment_doctor_id_fk FOREIGN KEY (doctor_id) REFERENCES doctor (id);
ALTER TABLE appointment ADD CONSTRAINT appointment_patient_id_fk FOREIGN KEY (patient_id) REFERENCES patient (id);
ALTER TABLE ask ADD CONSTRAINT ask_doctor_id_fk FOREIGN KEY (doctor_id) REFERENCES doctor (id);
ALTER TABLE ask ADD CONSTRAINT ask_patient_id_fk FOREIGN KEY (patient_id) REFERENCES patient (id);

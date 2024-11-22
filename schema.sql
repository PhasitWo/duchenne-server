
CREATE TABLE appointment (
  id int NOT NULL PRIMARY KEY,
  create_at timestamp,
  date timestamp,
  patient_id NOT NULL int,
  doctor_id NOT NULL int
);


CREATE TABLE ask (
  id int NOT NULL PRIMARY KEY,
  question varchar(500),
  answer varchar(500),
  patient_id NOT NULL int,
  doctor_id NOT NULL int
);


CREATE TABLE doctor (
  id int NOT NULL PRIMARY KEY,
  doctor_info varchar(500),
  username varchar(20),
  password varchar(20),
  role NOT NULL varchar(10)
);

CREATE TABLE patient (
  id int PRIMARY KEY,
  hn varchar(15),
  patient_info varchar(500),
  email varchar(255),
  phone varchar(15),
  password varchar(20),
  verified bool,
  verification_token varchar(500)
);

ALTER TABLE appointment ADD CONSTRAINT appointment_doctor_id_fk FOREIGN KEY (doctor_id) REFERENCES doctor (id);
ALTER TABLE appointment ADD CONSTRAINT appointment_patient_id_fk FOREIGN KEY (patient_id) REFERENCES patient (id);
ALTER TABLE ask ADD CONSTRAINT ask_doctor_id_fk FOREIGN KEY (doctor_id) REFERENCES doctor (id);
ALTER TABLE ask ADD CONSTRAINT ask_patient_id_fk FOREIGN KEY (patient_id) REFERENCES patient (id);

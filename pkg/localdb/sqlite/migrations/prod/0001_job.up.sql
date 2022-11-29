create table job (
  id varchar(255) PRIMARY KEY,
  created timestamp,
  jobdata json
);

create table job_event (
  id SERIAL PRIMARY KEY,
  job_id varchar(255),
  created timestamp,
  eventdata json,
  FOREIGN KEY(job_id) REFERENCES job(id)
);

create table local_event (
  id SERIAL PRIMARY KEY,
  job_id varchar(255),
  created timestamp,
  eventdata json,
  FOREIGN KEY(job_id) REFERENCES job(id)
);
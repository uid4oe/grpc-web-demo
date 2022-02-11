create table if not exists catalog (
  id text not null,
  name text not null,
  details text not null,
	total_orders integer,
	average_cost integer,
  primary key (id)
);

insert into catalog (id,name,details,total_orders,average_cost) values ('f5cc1d93-ef74-4bee-a545-013bee9d2f2a', 'Time Travel', 'A very complex way to save time', 10,100),
('71984813-8f1b-4380-8f41-1909c237cf61', 'Recruit Superheros', 'No need to worry about the job fit', 20,200),
('5d903ae4-5f45-450f-a2fe-a6489a4efd92', 'Clone HR People', 'Boost candidate assessment speed', 30,300);

create table if not exists orders (
  id text not null,
  offer_id text not null,
  rating float,
  primary key (id)
);

create table if not exists order_details (
  step integer,
  catalog_id text,
  detail text
);

insert into order_details (step,catalog_id,detail) values 
(1,'5d903ae4-5f45-450f-a2fe-a6489a4efd92','Order Received! Waiting partner''s approval.'),
(2,'5d903ae4-5f45-450f-a2fe-a6489a4efd92','Confirmed! Checking cloning machine updates.'),
(3,'5d903ae4-5f45-450f-a2fe-a6489a4efd92','Filling 3D cartridge.'),
(4,'5d903ae4-5f45-450f-a2fe-a6489a4efd92','Almost done, hang on!'),
(5,'5d903ae4-5f45-450f-a2fe-a6489a4efd92','Oh no! Angry Minya and Linnea noises!!!'),
(6,'5d903ae4-5f45-450f-a2fe-a6489a4efd92','Done, please rate the service.'),
(1,'71984813-8f1b-4380-8f41-1909c237cf61','Order Received! Waiting partner''s approval.'),
(2,'71984813-8f1b-4380-8f41-1909c237cf61','Confirmed! Looking for The Avengers contact info.'),
(3,'71984813-8f1b-4380-8f41-1909c237cf61','Interviewing Thanos. (He offers 50% reduction in the workload, interesting..)'),
(4,'71984813-8f1b-4380-8f41-1909c237cf61','Answering questions about the Skatteverket.'),
(5,'71984813-8f1b-4380-8f41-1909c237cf61','Agreed on terms, one infinity stone annualy.'),
(5,'71984813-8f1b-4380-8f41-1909c237cf61','Done, please rate the service.'),
(1,'f5cc1d93-ef74-4bee-a545-013bee9d2f2a','Order Received! Waiting partner''s approval.'),
(2,'f5cc1d93-ef74-4bee-a545-013bee9d2f2a','Confirmed! Asking for the Doctor Who''s assitance.'),
(3,'f5cc1d93-ef74-4bee-a545-013bee9d2f2a','Setting up the TARDIS to the time where you found the best candidate.'),
(4,'f5cc1d93-ef74-4bee-a545-013bee9d2f2a','Wow! Looks like it took you 3 months to find the best candidate.'),
(5,'f5cc1d93-ef74-4bee-a545-013bee9d2f2a','Sucessfully teleported.'),
(6,'f5cc1d93-ef74-4bee-a545-013bee9d2f2a','Done, please rate the service.');
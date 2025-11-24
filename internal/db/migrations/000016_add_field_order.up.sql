alter table orders
add column payment_status varchar(50) not null default 'UNPAID',
add column transaction_id varchar(255),
add column delivered_at timestamp;
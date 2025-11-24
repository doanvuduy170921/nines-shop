create table order_status_history(
    id serial primary key ,
    order_id integer not null references orders(id),
    status VARCHAR(50) NOT NULL,
    note text,
    created_at timestamp default current_timestamp not null
);
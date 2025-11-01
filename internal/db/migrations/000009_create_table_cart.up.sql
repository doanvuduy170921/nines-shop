-- create table cart
create table if not exists cart(
    id BIGSERIAL primary key,
    user_id int not null references users(id) on delete cascade,
    product_id int not null references products(id) on delete cascade,
    quantity int default 1 check(quantity >0),
    created_at timestamp default CURRENT_TIMESTAMP,
    updated_at timestamp default CURRENT_TIMESTAMP
);


-- create index for better performance
CREATE INDEX idx_cart_user_id ON cart(user_id);
CREATE INDEX idx_cart_product_id ON cart(product_id);


-- create function update_updated on column
create or replace function update_cart_updated_at()
returns trigger as $$
begin
    NEW.updated_at = CURRENT_TIMESTAMP;
return NEW;
end;
$$ language plpgsql ;

 create trigger trg_cart_updated_at
     before update on cart
for each row
execute function update_cart_updated_at();

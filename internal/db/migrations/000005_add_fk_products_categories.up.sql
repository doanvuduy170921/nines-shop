alter table products
add constraint fk_product_category foreign key (category_id) references categories(id) on delete set null on update cascade ;
CREATE TABLE IF NOT EXISTS product_images (
                                id          SERIAL PRIMARY KEY,
                                product_id  BIGINT  NOT NULL REFERENCES products(id) ON DELETE CASCADE,
                                image_url   TEXT NOT NULL
);

-- create index for better query performance
create index  idx_product_image_url on product_images(image_url);
create index  idx_product_id on product_images(product_id);
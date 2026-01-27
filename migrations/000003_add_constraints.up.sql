ALTER TABLE cargo_type
ADD CONSTRAINT cargo_type_title_uq UNIQUE (title);

ALTER TABLE operation
ADD CONSTRAINT operation_title_uq UNIQUE (title);
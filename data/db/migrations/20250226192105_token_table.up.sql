CREATE TABLE tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    access_token VARCHAR(1024) NOT NULL UNIQUE,
    refresh_token VARCHAR(1024) NOT NULL UNIQUE,
    user_id UUID NOT NULL,
    CONSTRAINT fk_user
      FOREIGN KEY(user_id)
        REFERENCES users(id)
          ON DELETE CASCADE
);

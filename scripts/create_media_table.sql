CREATE TABLE IF NOT EXISTS media (
    id SERIAL PRIMARY KEY,
    account_id VARCHAR(255) NOT NULL,
    folder_id VARCHAR(255) NOT NULL,
    media_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    is_deleted BOOLEAN,
    CONSTRAINT fk_account FOREIGN KEY (account_id) REFERENCES users(account_id),
    CONSTRAINT fk_folder FOREIGN KEY (folder_id) REFERENCES folders(folder_id)
);
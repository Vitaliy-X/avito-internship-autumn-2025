CREATE INDEX idx_users_team_active
    ON users (team_name, is_active);

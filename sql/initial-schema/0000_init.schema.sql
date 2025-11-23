CREATE TABLE teams
(
    team_name TEXT PRIMARY KEY
);

CREATE TABLE users
(
    user_id   TEXT PRIMARY KEY,
    username  TEXT    NOT NULL,
    team_name TEXT    NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    FOREIGN KEY (team_name)
        REFERENCES teams (team_name)
        ON DELETE CASCADE
);

CREATE TABLE pull_requests
(
    pull_request_id    TEXT PRIMARY KEY,
    pull_request_name  TEXT NOT NULL,
    author_id          TEXT NOT NULL REFERENCES users (user_id),
    status             TEXT NOT NULL CHECK (status IN ('OPEN', 'MERGED')),
    assigned_reviewers TEXT[] NOT NULL DEFAULT '{}',
    created_at         TIMESTAMPTZ DEFAULT NOW(),
    merged_at          TIMESTAMPTZ
);

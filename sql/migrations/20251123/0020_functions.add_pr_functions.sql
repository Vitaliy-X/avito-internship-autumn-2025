CREATE OR REPLACE FUNCTION get_active_team_members(team TEXT)
    RETURNS TEXT[]
    LANGUAGE plpgsql
AS
'
DECLARE
    ids TEXT[];
BEGIN
    SELECT ARRAY(
        SELECT user_id FROM users
        WHERE team_name = team AND is_active = TRUE
        ORDER BY user_id
    )
    INTO ids;
    RETURN COALESCE(ids, ARRAY[]::TEXT[]);
END';

CREATE OR REPLACE FUNCTION pick_random_active_member(team TEXT, exclude_ids TEXT[])
    RETURNS TEXT
    LANGUAGE plpgsql
AS
'
DECLARE
    result TEXT;
BEGIN
    SELECT user_id INTO result FROM users
    WHERE team_name = team
      AND is_active = TRUE
      AND NOT (user_id = ANY (exclude_ids))
    ORDER BY random()
        LIMIT 1;
    RETURN result;
END';

CREATE OR REPLACE FUNCTION assign_reviewers_for_pr(pr_id TEXT)
    RETURNS TEXT[]
    LANGUAGE plpgsql
AS
'
DECLARE
    author      TEXT;
    team_name   TEXT;
    members     TEXT[];
    reviewers   TEXT[] := ARRAY[]::TEXT[];
    candidate   TEXT;
BEGIN
    SELECT author_id
    INTO author
    FROM pull_requests
    WHERE pull_request_id = pr_id;

    IF author IS NULL THEN
        RAISE EXCEPTION ''PR not found: %'', pr_id;
    END IF;

    SELECT team_name
    INTO team_name
    FROM users
    WHERE user_id = author;

    IF team_name IS NULL THEN
        RAISE EXCEPTION ''Team for author % not found'', author;
    END IF;

    members := get_active_team_members(team_name);

    members := ARRAY(
        SELECT unnest(members)
        WHERE unnest <> author
    );

    SELECT pick_random_active_member(team_name, reviewers || author)
    INTO candidate;
    IF candidate IS NOT NULL THEN
        reviewers := reviewers || candidate;
    END IF;

    SELECT pick_random_active_member(team_name, reviewers || author)
    INTO candidate;
    IF candidate IS NOT NULL THEN
        reviewers := reviewers || candidate;
    END IF;

    UPDATE pull_requests
    SET assigned_reviewers = reviewers
    WHERE pull_request_id = pr_id;

    RETURN reviewers;
END';

CREATE OR REPLACE FUNCTION reassign_reviewer(pr_id TEXT, old_user TEXT)
    RETURNS TEXT
    LANGUAGE plpgsql
AS
'
DECLARE
    pr_status   TEXT;
    reviewers   TEXT[];
    new_user    TEXT;
    old_team    TEXT;
BEGIN
    SELECT status, assigned_reviewers
    INTO pr_status, reviewers
    FROM pull_requests
    WHERE pull_request_id = pr_id;

    IF pr_status IS NULL THEN
        RAISE EXCEPTION ''PR not found'';
    END IF;

    IF pr_status = ''MERGED'' THEN
        RAISE EXCEPTION ''merged'';
    END IF;

    IF NOT (old_user = ANY (reviewers)) THEN
        RAISE EXCEPTION ''not_assigned'';
    END IF;

    SELECT team_name
    INTO old_team
    FROM users
    WHERE user_id = old_user;

    IF old_team IS NULL THEN
        RAISE EXCEPTION ''user_not_found'';
    END IF;

    SELECT pick_random_active_member(old_team, reviewers)
    INTO new_user;
    IF new_user IS NULL THEN
        RAISE EXCEPTION ''no_candidate'';
    END IF;

    reviewers := ARRAY(
        SELECT CASE WHEN r = old_user THEN new_user ELSE r END
        FROM unnest(reviewers) AS r
    );

    UPDATE pull_requests
    SET assigned_reviewers = reviewers
    WHERE pull_request_id = pr_id;

    RETURN new_user;
END';

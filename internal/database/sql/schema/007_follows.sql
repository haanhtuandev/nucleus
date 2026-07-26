-- +goose Up
CREATE TABLE follows (
    follower_id UUID NOT NULL,
    followee_id UUID NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_follower FOREIGN KEY (follower_id) REFERENCES users(id),
    CONSTRAINT fk_followee FOREIGN KEY (followee_id) REFERENCES users(id),
    CONSTRAINT pk_follow    PRIMARY KEY (follower_id, followee_id)
);


-- +goose Down
DROP TABLE follows;
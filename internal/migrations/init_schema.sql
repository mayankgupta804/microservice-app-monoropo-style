-- +migrate Up
CREATE TABLE incidents (
    id int NOT NULL AUTO_INCREMENT,
    message text,
    status varchar(10),
    ack varchar(3),
    PRIMARY KEY (id)
);

CREATE TABLE comments (
    id INT NOT NULL AUTO_INCREMENT,
    incident_id INT NOT NULL,
    comment text,
    PRIMARY KEY(id),
    FOREIGN KEY (incident_id) REFERENCES incidents(id)
);

-- +migrate Down
DROP TABLE comments;
DROP TABLE incidents;

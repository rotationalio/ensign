-- ULID: 01GKHJRF01YXHZ51YMMKV3RCMK
INSERT INTO organizations (id, name, domain, created, modified) VALUES
    (x'0184e32c3c01f763f287d4a4f63c3293', 'Testing', 'example.com', '2022-12-05T16:43:57.825256Z', '2022-12-05T16:43:57.825256Z')
;

-- ULID: 01GKHJSK7CZW0W282ZN3E9W86Z
-- Password: theeaglefliesatmidnight
INSERT INTO users (id, name, email, password, last_login, created, modified) VALUES
    (x'0184e32cccecff01c1205fa8dc9e20df', 'Jannel P. Hudson', 'jannel@example.com', '$argon2id$v=19$m=65536,t=1,p=2$Ujy6FI2NBqRIUHmqH0YcQA==$f1lwLv4DpE4OTkMq3sTShZS3NHADg9UvnZNHtuUOmZ8=', '2022-12-13T01:22:39Z', '2022-12-05T16:44:34.924036Z', '2022-12-05T16:44:34.924036Z')
;

INSERT INTO organization_users (organization_id, user_id, created, modified) VALUES
    (x'0184e32c3c01f763f287d4a4f63c3293', x'0184e32cccecff01c1205fa8dc9e20df', '2022-12-05T16:44:35.00123Z', '2022-12-05T16:44:35.00123Z')
;

INSERT INTO user_roles (user_id, role_id, created, modified) VALUES
    (x'0184e32cccecff01c1205fa8dc9e20df', 1, '2022-12-05T16:44:35.00123Z', '2022-12-05T16:44:35.00123Z')
;
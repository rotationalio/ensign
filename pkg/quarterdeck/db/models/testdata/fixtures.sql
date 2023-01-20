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

-- ULID: 01GME02TJP2RRP39MKR525YDQ6
-- Client Secret: wAfRpXLTiWn7yo7HQzOCwxMvveqiHXoeVJghlSIK2YbMqOMCUiSVRVQOLT0ORrVS
-- ProjectID: 01GQ7P8DNR9MR64RJR9D64FFNT
INSERT INTO api_keys (id, key_id, secret, name, organization_id, project_id, created_by, created, modified) VALUES
    (x'01851c016a56163161a693c1445f36e6', 'DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa', '$argon2id$v=19$m=65536,t=1,p=2$5tE7XLSdqM36DUmzeSppvA==$eTfRYSCuBssAcuxxFv/eh92CyL1NuNqBPkhlLoIAVAw=', 'Eagle Publishers', x'0184e32c3c01f763f287d4a4f63c3293', x'0185cf6436b84d306262584b4c47beba', NULL, '2022-12-05T16:48:21.123332Z', '2022-12-05T16:48:21.123332Z')
;

INSERT INTO api_key_permissions (api_key_id, permission_id, created, modified) VALUES
    (x'01851c016a56163161a693c1445f36e6', 11, '2022-12-05T16:48:21.123332Z', '2022-12-05T16:48:21.123332Z'),
    (x'01851c016a56163161a693c1445f36e6', 14, '2022-12-05T16:48:21.123332Z', '2022-12-05T16:48:21.123332Z'),
    (x'01851c016a56163161a693c1445f36e6', 15, '2022-12-05T16:48:21.123332Z', '2022-12-05T16:48:21.123332Z'),
    (x'01851c016a56163161a693c1445f36e6', 16, '2022-12-05T16:48:21.123332Z', '2022-12-05T16:48:21.123332Z'),
    (x'01851c016a56163161a693c1445f36e6', 17, '2022-12-05T16:48:21.123332Z', '2022-12-05T16:48:21.123332Z')
;

-- ULID example.com: 01GKHJRF01YXHZ51YMMKV3RCMK
-- ULID checkers.io: 01GQFQ14HXF2VC7C1HJECS60XX
INSERT INTO organizations (id, name, domain, created, modified) VALUES
    (x'0184e32c3c01f763f287d4a4f63c3293', 'Testing', 'example.com', '2022-12-05T16:43:57.825256Z', '2022-12-05T16:43:57.825256Z'),
    (x'0185df70923d78b6c3b03193999303bd', 'Checkers', 'checkers.io', '2023-01-23T16:22:54.781762Z', '2023-01-23T16:22:54.781762Z')
;

-- Jannel ULID: 01GKHJSK7CZW0W282ZN3E9W86Z
-- Jannel Password: theeaglefliesatmidnight
-- Edison ULID: 01GQFQ4475V3BZDMSXFV5DK6XX
-- Edison Password: supersecretssquirrel
INSERT INTO users (id, name, email, password, last_login, created, modified) VALUES
    (x'0184e32cccecff01c1205fa8dc9e20df', 'Jannel P. Hudson', 'jannel@example.com', '$argon2id$v=19$m=65536,t=1,p=2$Ujy6FI2NBqRIUHmqH0YcQA==$f1lwLv4DpE4OTkMq3sTShZS3NHADg9UvnZNHtuUOmZ8=', '2022-12-13T01:22:39Z', '2022-12-05T16:44:34.924036Z', '2022-12-05T16:44:34.924036Z'),
    (x'0185df7210e5d8d7f6d33d7ecad99bbd', 'Edison Edgar Franklin', 'eefrank@checkers.io', '$argon2id$v=19$m=65536,t=1,p=2$x4Zh4ARSD4wK7uZFaauyjg==$eCkUszypW+rLvQ+D9lpfTgVwqPSKH13rCdmzV9vZ8cQ=', '2023-02-14T14:48:08Z', '2023-01-23T16:24:32.741955Z', '2023-01-23T16:24:32.741955Z')
;

INSERT INTO organization_users (organization_id, user_id, created, modified) VALUES
    (x'0184e32c3c01f763f287d4a4f63c3293', x'0184e32cccecff01c1205fa8dc9e20df', '2022-12-05T16:44:35.00123Z', '2022-12-05T16:44:35.00123Z'),
    (x'0185df70923d78b6c3b03193999303bd', x'0185df7210e5d8d7f6d33d7ecad99bbd', '2023-01-23T16:24:32.741955Z', '2023-01-23T16:24:32.741955Z')
;

INSERT INTO user_roles (user_id, role_id, created, modified) VALUES
    (x'0184e32cccecff01c1205fa8dc9e20df', 1, '2022-12-05T16:44:35.00123Z', '2022-12-05T16:44:35.00123Z'),
    (x'0185df7210e5d8d7f6d33d7ecad99bbd', 1, '2023-01-23T16:24:32.741955Z', '2023-01-23T16:24:32.741955Z')
;

-- Row, ULID, Secret, ProjectID, OrgName
-- 0, 01GME02TJP2RRP39MKR525YDQ6, wAfRpXLTiWn7yo7HQzOCwxMvveqiHXoeVJghlSIK2YbMqOMCUiSVRVQOLT0ORrVS, 01GQ7P8DNR9MR64RJR9D64FFNT, Example (Bird)
-- 1, 01GQFQC3F7B42YJWVRXYBE80QB, W0ZROSP0FMy2fvv5GisoeivZdQlyXeZRspD8IzfjUZpVSDEWRE9UOP11IYyFhO9w, 01GQFQCFC9P3S7QZTPYFVBJD7F, Checkers
-- 2, 01GQFQPH67WPJJET1HFNCXDZ3A, 6VHrsfUyRH1ODNWjnJwco4g1MMauaqNYkHY1fPft9xoaRoMiQdPzZjxBMxM6i2A4, 01GQFQCFC9P3S7QZTPYFVBJD7F, Checkers
-- 3, 01GQFQT8YGR9MQSGY1VD5DEMFE, kwqSUyMPQYrzjAw9RR8y46MhuPEEwhv3dVx6aECxXYo129a2jn7dd47gUsxpg9K1, 01GQFQCFC9P3S7QZTPYFVBJD7F, Checkers
-- 4, 01GQFRCEX7HMJSX5Q0G0DVHHPS, 60nUF37uk9uFBNW93mhbji7rQ7ujPBrmCB2GlQQmqnTrOFBYUokewnoYM2DZfd0N, 01GQFR0KM5S2SSJ8G5E086VQ9K, Example (Project)
-- 5, 01GQFRG53XKBJDY50FXR0107ZQ, 1yLncMfrI6yHZYYvC0OG2DW4LzO2issNfm0sj6Rb6DBqCNbkRu3Ob4ODfTPKKJNm, 01GQFR0KM5S2SSJ8G5E086VQ9K, Example (Project)
-- 6, 01GQFRJ5VEDG1XSPXQAVFEJRAB, SakdbTabXJ2qbyezZLPwEbuaI83hH2Wb1IJJrnwrzKCpuDlpvLUptuSAerpA80ce, 01GQFR0KM5S2SSJ8G5E086VQ9K, Example (Project)
-- 7, 01GQFRKWTKZWHSAX2GK16F5Q1V, rQntyIwNdkQ7ZpuXnP2SONAx8CCrIVVaXg2bLtgIl1Iow0p0qwVLP41iCKPsD3d4, 01GQFR0KM5S2SSJ8G5E086VQ9K, Example (Project)
-- 8, 01GQFRP3SWM5H5EFXSA5Z33YZ5, VFoAAG69UUQ7Pkobyr3V391kiNBERiP4NJ5dZlKwG90XCJkHqGzySyY6HElnCDGn, 01GQFR0KM5S2SSJ8G5E086VQ9K, Example (Project)
-- 9, 01GQFRR52DGJMG2JTJV2WTQGMZ, n4qTNIRHcrXIo44nVzisaQtoBW48QDZ7eY94rWSAqNF4B6ta7XJkzbE8my1HAvsm, 01GQFR0KM5S2SSJ8G5E086VQ9K, Example (Project)
-- 10, 01GQFRT0EGEZYJAMKZ8PE0P5R8, KLTlCeVy38etl2YBM6KZrTenz0fe2P1ZGsF9XXVOYzsO0R4BrhYWpf7ASvkeWDyo, 01GQFR0KM5S2SSJ8G5E086VQ9K, Example (Project)
-- 11, 01GQFRVWVH76934TC9WHVMB2W5, 3qiq8HJyuViqk7aR5c6OyCaGGLxqvRaZ0v6bCP7xobQLPnxaflD0TfUVhxSDbMrg, 01GQFR0KM5S2SSJ8G5E086VQ9K, Example (Project)
-- 12, 01GQFRXM7DA2JG7K1MS4JX92ZF, ToCSvbjz3SwB00lOxMJj0OEd1bn8gfZSRG6Hlo5VRSm4HIS8gQtBz8hAZG5inS8c, 01GQFR0KM5S2SSJ8G5E086VQ9K, Example (Project)
-- 13, 01GQFRZAG98EAEE3Y8D65F0YDR, btRCwtw18w2Ty0WMp9oB9INsWIK6JkEbpI7iUlR8lwlZ02ZIBWO7YAcJbjwnMUlv, 01GQ7P8DNR9MR64RJR9D64FFNT, Example (Bird)
INSERT INTO api_keys (id, key_id, secret, name, organization_id, project_id, created_by, source, user_agent, last_used, created, modified) VALUES
    (x'01851c016a56163161a693c1445f36e6', 'DbIxBEtIUgNIClnFMDmvoZeMrLxUTJVa', '$argon2id$v=19$m=65536,t=1,p=2$5tE7XLSdqM36DUmzeSppvA==$eTfRYSCuBssAcuxxFv/eh92CyL1NuNqBPkhlLoIAVAw=', 'Eagle Publishers', x'0184e32c3c01f763f287d4a4f63c3293', x'0185cf6436b84d306262584b4c47beba', x'0184e32cccecff01c1205fa8dc9e20df', 'Beacon UI', 'Quarterdeck API/v1', '2023-01-22T13:26:25.394129Z', '2022-12-05T16:48:21.123332Z', '2022-12-05T16:48:21.123332Z'),
    (x'0185df760de75905e97378ef96e402eb', 'quYyMDzuNTswqGPIvDhwvzrSkVyCTTfp', '$argon2id$v=19$m=65536,t=1,p=2$8KR6dbPgM6TCSjfAbSw70A==$7UC5eTFIHHHJB9baYHYai2C/63e5H4TiX6NtZMbZoBY=', 'Checkers Publishers', x'0185df70923d78b6c3b03193999303bd', x'0185df763d89b0f27bff56f3f6b934ef', x'0185df7210e5d8d7f6d33d7ecad99bbd', 'Beacon UI', 'Quarterdeck API/v1', NULL, '2023-01-23T16:29:06.313194Z', '2023-01-23T16:29:06.313194Z'),
    (x'0185df7b44c7e5a52768317d59d6fc6a', 'OtOuXvOGqOzxkrXEttPDgCXuDrJMELua', '$argon2id$v=19$m=65536,t=1,p=2$u2rSUKFqtBhICLNN0mp3/Q==$iGtzfRfQrVGl3U581F+5NGasgi/1xx1iL5vg7wjlmuw=', 'Checkers Subscribers', x'0185df70923d78b6c3b03193999303bd', x'0185df763d89b0f27bff56f3f6b934ef', x'0185df7210e5d8d7f6d33d7ecad99bbd', 'Beacon UI', 'Quarterdeck API/v1', NULL, '2023-01-23T16:34:35.847672Z', '2023-01-23T16:34:35.847672Z'),
    (x'0185df7d23d0c2697cc3c1db4ad751ee', 'jVZugZpfCvYQiJBQESDIBpCNidVVjnSO', '$argon2id$v=19$m=65536,t=1,p=2$KPPD+fNEHB70F1vmzmguPQ==$Xu9sjGDP5/bL5yYvuviTYS+izSY/iB2L/FbFvPfLIGY=', 'Checkers Topic Manager', x'0185df70923d78b6c3b03193999303bd', x'0185df763d89b0f27bff56f3f6b934ef', x'0185df7210e5d8d7f6d33d7ecad99bbd', 'Beacon UI', 'Quarterdeck API/v1', NULL, '2023-01-23T16:36:38.480154Z', '2023-01-23T16:36:38.480154Z'),
    (x'0185df863ba78d259e96e0801bb8c6d9', 'fExRbbTgnLWlxyucKpuYOmTCIMqNnGMU', '$argon2id$v=19$m=65536,t=1,p=2$rHWSWs4F0FTr0ELFIqkaOA==$kmA914OuPG6KV/URmchZPjk3pQVDRpvb0gEoA9xzZcI=', 'Project Key 1', x'0184e32c3c01f763f287d4a4f63c3293', x'0185df804e85c8b399220570106ddd33', x'0184e32cccecff01c1205fa8dc9e20df', 'Beacon UI', 'Quarterdeck API/v1', NULL, '2023-01-23T16:46:34.407159Z', '2023-01-23T16:46:34.407159Z'),
    (x'0185df88147d9ae4df140fee00101ff7', 'TljJVeQqfymUbhFbhTbZiqnXhnSLNJle', '$argon2id$v=19$m=65536,t=1,p=2$Vqv9T6E7c1NwvnVeBEYS6w==$ytwwA6yTSEgxBlNaUHjEH9hFHF38BvckKhLuNTZ+Nqo=', 'Project Key 2', x'0184e32c3c01f763f287d4a4f63c3293', x'0185df804e85c8b399220570106ddd33', x'0184e32cccecff01c1205fa8dc9e20df', 'Beacon UI', 'Quarterdeck API/v1', NULL, '2023-01-23T16:48:35.453808Z', '2023-01-23T16:48:35.453808Z'),
    (x'0185df89176e6c03dcdbb756dee9614b', 'NokCWMzDNTZQBmhEMqOiBmIRYYsFYOJV', '$argon2id$v=19$m=65536,t=1,p=2$1ymdXHNraeiMQCAkbkNjUA==$+saAD5/6DHY9ZwzdvOlwZagPh13uVLruWhe6FsRgXf8=', 'Project Key 3', x'0184e32c3c01f763f287d4a4f63c3293', x'0185df804e85c8b399220570106ddd33', x'0184e32cccecff01c1205fa8dc9e20df', 'Beacon UI', 'Quarterdeck API/v1', NULL, '2023-01-23T16:49:41.743053Z', '2023-01-23T16:49:41.743053Z'),
    (x'0185df89f353ff23957450984cf2dc3b', 'LitfLhZmDedaunvRUoiiPNIjuepUGiId', '$argon2id$v=19$m=65536,t=1,p=2$AjVFpspFXe/uWSoPw3vv2g==$Ifb6zx8D9bIuLbS47GsjggS3kgJlG6TD9sEIoVWor/I=', 'Project Key 4', x'0184e32c3c01f763f287d4a4f63c3293', x'0185df804e85c8b399220570106ddd33', x'0184e32cccecff01c1205fa8dc9e20df', 'Beacon UI', 'Quarterdeck API/v1', NULL, '2023-01-23T16:50:38.035102Z', '2023-01-23T16:50:38.035102Z'),
    (x'0185df8b0f3ca162573fb9517e31fbe5', 'aMoOGzLiORKmTNTQtePhLNRflzyNwyLr', '$argon2id$v=19$m=65536,t=1,p=2$odHs5ybjrbxdZ2D9QF5TNg==$nh8Sev6yRt5NiGVvshWAMUEh9X4hhGfgVQbpmC2ltqA=', 'Project Key 5', x'0184e32c3c01f763f287d4a4f63c3293', x'0185df804e85c8b399220570106ddd33', x'0184e32cccecff01c1205fa8dc9e20df', 'Beacon UI', 'Quarterdeck API/v1', NULL, '2023-01-23T16:51:50.71668Z', '2023-01-23T16:51:50.71668Z'),
    (x'0185df8c144d84a9014b52d8b9abc29f', 'UpECadqhKTpVYSmZJYpajYaOAuOohBhF', '$argon2id$v=19$m=65536,t=1,p=2$NWvOmRiJHH2+/fPBENlufA==$vMcJ1vfsmjNPXdYnZvTjBzx0lagyjpzvaB8Gc53v3JM=', 'Project Key 6', x'0184e32c3c01f763f287d4a4f63c3293', x'0185df804e85c8b399220570106ddd33', x'0184e32cccecff01c1205fa8dc9e20df', 'Beacon UI', 'Quarterdeck API/v1', NULL, '2023-01-23T16:52:57.549172Z', '2023-01-23T16:52:57.549172Z'),
    (x'0185df8d01d077fd25527f459c0b1708', 'NdMTPpYFTcqefZcwxPNuEHvLBjauWBZf', '$argon2id$v=19$m=65536,t=1,p=2$qg1tOPZCI8fxRNZO8gD4Gw==$AsrvtypadtHWkbZiEHzvil6dTcn6WbPTeIrPiv8hRns=', 'Project Key 7', x'0184e32c3c01f763f287d4a4f63c3293', x'0185df804e85c8b399220570106ddd33', x'0184e32cccecff01c1205fa8dc9e20df', 'Beacon UI', 'Quarterdeck API/v1', NULL, '2023-01-23T16:53:58.352279Z', '2023-01-23T16:53:58.352279Z'),
    (x'0185df8df3713992326989e477458b85', 'JvDjsgQWrCIAqkmIMyVrxYwVuKIxuwog', '$argon2id$v=19$m=65536,t=1,p=2$7xgMKRYv1TBodhS0KafBNg==$BqnflEqi8j6e08XFDj/b4aqckpNkgfZSWgHD8UMfMNU=', 'Project Key 8', x'0184e32c3c01f763f287d4a4f63c3293', x'0185df804e85c8b399220570106ddd33', x'0184e32cccecff01c1205fa8dc9e20df', 'Beacon UI', 'Quarterdeck API/v1', NULL, '2023-01-23T16:55:00.209135Z', '2023-01-23T16:55:00.209135Z'),
    (x'0185df8ed0ed50a503cc34c925d48bef', 'aDbUYAqUlGqXCmTaHExRnTsAclVsxKXq', '$argon2id$v=19$m=65536,t=1,p=2$WGhcFEyqnUrJlTJ/B3484w==$kWbhWn5Dg8ZaAYGxhbznoUoVrhhgPwSpCCX/7YvQBSU=', 'Project Key 9', x'0184e32c3c01f763f287d4a4f63c3293', x'0185df804e85c8b399220570106ddd33', x'0184e32cccecff01c1205fa8dc9e20df', 'Beacon UI', 'Quarterdeck API/v1', NULL, '2023-01-23T16:55:56.9096Z', '2023-01-23T16:55:56.9096Z'),
    (x'0185df8faa094394e70fc8698af079b8', 'HFqkNtqwBXKGZGpowsLeMvbJNKZhQaBS', '$argon2id$v=19$m=65536,t=1,p=2$HWdlhUg6oBlhclXY1wmeWg==$ndABLT3EJniwD3vr87TPdel9LAxgQOT/klG+B4O1F14=', 'Turkey Publishers', x'0184e32c3c01f763f287d4a4f63c3293', x'0185cf6436b84d306262584b4c47beba', x'0184e32cccecff01c1205fa8dc9e20df', 'Beacon UI', 'Quarterdeck API/v1', '2023-02-14T20:12:49.394129Z', '2023-01-23T16:56:52.489425Z', '2023-01-23T16:56:52.489425Z')
;

INSERT INTO api_key_permissions (api_key_id, permission_id, created, modified) VALUES
    (x'01851c016a56163161a693c1445f36e6', 16, '2022-12-05T16:48:21.123332Z', '2022-12-05T16:48:21.123332Z'),
    (x'01851c016a56163161a693c1445f36e6', 19, '2022-12-05T16:48:21.123332Z', '2022-12-05T16:48:21.123332Z'),
    (x'01851c016a56163161a693c1445f36e6', 20, '2022-12-05T16:48:21.123332Z', '2022-12-05T16:48:21.123332Z'),
    (x'01851c016a56163161a693c1445f36e6', 21, '2022-12-05T16:48:21.123332Z', '2022-12-05T16:48:21.123332Z'),
    (x'01851c016a56163161a693c1445f36e6', 22, '2022-12-05T16:48:21.123332Z', '2022-12-05T16:48:21.123332Z'),
    (x'0185df760de75905e97378ef96e402eb', 21, '2023-01-23T16:29:06.313194Z', '2023-01-23T16:29:06.313194Z'),
    (x'0185df7b44c7e5a52768317d59d6fc6a', 22, '2023-01-23T16:34:35.847672Z', '2023-01-23T16:34:35.847672Z'),
    (x'0185df7d23d0c2697cc3c1db4ad751ee', 16, '2023-01-23T16:36:38.480154Z', '2023-01-23T16:36:38.480154Z'),
    (x'0185df7d23d0c2697cc3c1db4ad751ee', 17, '2023-01-23T16:36:38.480154Z', '2023-01-23T16:36:38.480154Z'),
    (x'0185df7d23d0c2697cc3c1db4ad751ee', 18, '2023-01-23T16:36:38.480154Z', '2023-01-23T16:36:38.480154Z'),
    (x'0185df7d23d0c2697cc3c1db4ad751ee', 19, '2023-01-23T16:36:38.480154Z', '2023-01-23T16:36:38.480154Z'),
    (x'0185df863ba78d259e96e0801bb8c6d9', 21, '2023-01-23T16:46:34.407159Z', '2023-01-23T16:46:34.407159Z'),
    (x'0185df863ba78d259e96e0801bb8c6d9', 22, '2023-01-23T16:46:34.407159Z', '2023-01-23T16:46:34.407159Z'),
    (x'0185df88147d9ae4df140fee00101ff7', 21, '2023-01-23T16:48:35.453808Z', '2023-01-23T16:48:35.453808Z'),
    (x'0185df88147d9ae4df140fee00101ff7', 22, '2023-01-23T16:48:35.453808Z', '2023-01-23T16:48:35.453808Z'),
    (x'0185df89176e6c03dcdbb756dee9614b', 21, '2023-01-23T16:49:41.743053Z', '2023-01-23T16:49:41.743053Z'),
    (x'0185df89176e6c03dcdbb756dee9614b', 22, '2023-01-23T16:49:41.743053Z', '2023-01-23T16:49:41.743053Z'),
    (x'0185df89f353ff23957450984cf2dc3b', 21, '2023-01-23T16:50:38.035102Z', '2023-01-23T16:50:38.035102Z'),
    (x'0185df89f353ff23957450984cf2dc3b', 22, '2023-01-23T16:50:38.035102Z', '2023-01-23T16:50:38.035102Z'),
    (x'0185df8b0f3ca162573fb9517e31fbe5', 21, '2023-01-23T16:51:50.71668Z', '2023-01-23T16:51:50.71668Z'),
    (x'0185df8b0f3ca162573fb9517e31fbe5', 22, '2023-01-23T16:51:50.71668Z', '2023-01-23T16:51:50.71668Z'),
    (x'0185df8c144d84a9014b52d8b9abc29f', 21, '2023-01-23T16:52:57.549172Z', '2023-01-23T16:52:57.549172Z'),
    (x'0185df8c144d84a9014b52d8b9abc29f', 22, '2023-01-23T16:52:57.549172Z', '2023-01-23T16:52:57.549172Z'),
    (x'0185df8d01d077fd25527f459c0b1708', 21, '2023-01-23T16:53:58.352279Z', '2023-01-23T16:53:58.352279Z'),
    (x'0185df8d01d077fd25527f459c0b1708', 22, '2023-01-23T16:53:58.352279Z', '2023-01-23T16:53:58.352279Z'),
    (x'0185df8df3713992326989e477458b85', 21, '2023-01-23T16:55:00.209135Z', '2023-01-23T16:55:00.209135Z'),
    (x'0185df8df3713992326989e477458b85', 22, '2023-01-23T16:55:00.209135Z', '2023-01-23T16:55:00.209135Z'),
    (x'0185df8ed0ed50a503cc34c925d48bef', 21, '2023-01-23T16:55:56.9096Z', '2023-01-23T16:55:56.9096Z'),
    (x'0185df8ed0ed50a503cc34c925d48bef', 22, '2023-01-23T16:55:56.9096Z', '2023-01-23T16:55:56.9096Z'),
    (x'0185df8faa094394e70fc8698af079b8', 19, '2023-01-23T16:56:52.489425Z', '2023-01-23T16:56:52.489425Z'),
    (x'0185df8faa094394e70fc8698af079b8', 20, '2023-01-23T16:56:52.489425Z', '2023-01-23T16:56:52.489425Z'),
    (x'0185df8faa094394e70fc8698af079b8', 21, '2023-01-23T16:56:52.489425Z', '2023-01-23T16:56:52.489425Z'),
    (x'0185df8faa094394e70fc8698af079b8', 22, '2023-01-23T16:56:52.489425Z', '2023-01-23T16:56:52.489425Z')
;

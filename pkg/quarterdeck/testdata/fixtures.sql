-- ULID example.com: 01GKHJRF01YXHZ51YMMKV3RCMK
-- ULID checkers.io: 01GQFQ14HXF2VC7C1HJECS60XX
-- ULID ghost.co: 01GYAVA5ARPRC5Y5CHRJDV34CT
INSERT INTO organizations (id, name, domain, created, modified) VALUES
    (x'0184e32c3c01f763f287d4a4f63c3293', 'Testing', 'example-com', '2022-12-05T16:43:57.825256Z', '2022-12-05T16:43:57.825256Z'),
    (x'0185df70923d78b6c3b03193999303bd', 'Checkers', 'checkers-io', '2023-01-23T16:22:54.781762Z', '2023-01-23T16:22:54.781762Z'),
    (x'018795b51558b6185f1591c49bb1919a', 'Ghost', 'ghost-co', '2023-04-18T18:51:25.400137Z', '2023-04-18T18:51:25.400137Z')
;

-- Jannel ULID: 01GKHJSK7CZW0W282ZN3E9W86Z
-- Jannel Password: theeaglefliesatmidnight
-- Edison ULID: 01GQFQ4475V3BZDMSXFV5DK6XX
-- Edison Password: supersecretssquirrel
-- Zendaya ULID: 01GQYYKY0ECGWT5VJRVR32MFHM
-- Zendaya Password: iseeallthings
-- Sophia Thompson ULID: 01GRKWY7MD5HFMZQ4HZZG16MYY
-- Sophia Thompson Password: livingthedream
-- Robert Millser ULID: 01GRM2319FNCQM7H5P6H9PJ68A
-- Robert Miller Password: weareonamission
INSERT INTO users (id, name, email, password, terms_agreement, privacy_agreement, email_verified, email_verification_expires, email_verification_token, email_verification_secret, last_login, created, modified) VALUES
    (x'0184e32cccecff01c1205fa8dc9e20df', 'Jannel P. Hudson', 'jannel@example.com', '$argon2id$v=19$m=65536,t=1,p=2$Ujy6FI2NBqRIUHmqH0YcQA==$f1lwLv4DpE4OTkMq3sTShZS3NHADg9UvnZNHtuUOmZ8=', 't', 't', 'f', '2023-03-01T16:53:45.641698Z', 'EpiLbYGb58xsOsjk2CWaNMOS0s-LCyW1VVvKrZNg7dI', x'06295afac37f555f0a2d812a16eab8eb1e6de95168c85cda8f0a8191af93e3fd6b3284f11b2b57e569a20f54ec07c40e0e0cb8ec2fa3b0a435a75c4aa379c4bf2532a495efb0ac0b91ac77b2d9f03049281d6977016dd63250c859ad8336a0bf28daa4163357cfe87e9914f752e19cdefc6eb0e4cdeec147a93af7adbac96971', '2022-12-13T01:22:39Z', '2022-12-05T16:44:34.924036Z', '2022-12-05T16:44:34.924036Z'),
    (x'0185df7210e5d8d7f6d33d7ecad99bbd', 'Edison Edgar Franklin', 'eefrank@checkers.io', '$argon2id$v=19$m=65536,t=1,p=2$x4Zh4ARSD4wK7uZFaauyjg==$eCkUszypW+rLvQ+D9lpfTgVwqPSKH13rCdmzV9vZ8cQ=', 't', 't', 't', '2088-03-02T17:03:05.849287Z', 'F4wS_AhNP1daEXtyhxVwBdDd1UkbCL4G0xSB0nSY_l4', x'b8d140b3f6a5ea538adc2581e40c4fb9c7f13e87097186fb1b18083d64f30e7f517ba9f9b1ff5703573e750999bc7002e33fbc3a1989ad2d5bd9dad98e6b186f676f06fbd75dba81168fa14d543f2fa491b0a8c27a2dca1325fe326a1ad5395e73b8e0aa0ca2664705e5da1f9ebdd26b1fdb2dc4f503692e7c4cbf2c20220cf6', '2023-02-14T14:48:08Z', '2023-01-23T16:24:32.741955Z', '2023-01-23T16:24:32.741955Z'),
    (x'0185fde9f80e6439a2ee58de062a3e34', 'Zendaya Longeye', 'zendaya@testing.io', '$argon2id$v=19$m=65536,t=1,p=2$rQMSo/Lksd+/DazFmcuu4Q==$GtZGSh9SajnzXp/Cd8h/zpzgXrw4coXhRz/DhnG7GEU=', 't', 't', 't', '2023-03-02T17:05:56.918467Z', 'zOeWtHf0DHnifwmVMp8uSTYCIPbPTi2HbFF7-D3uUt4', x'793edb6470a78fe744af7fa2abe778f6549a47380cd76db7143fb78893873bb93856b1a919d3dfdd6698968fca53b32c78ac78a8fde0bcc0f7992d094b374538790c6a0aaa3b7e818d2cdeb267d9f8bb5a55910275cafa32733ec611f20151f21321f847858acc76faffb6cc50805b200ffc93f7713f1bb5ca5150cafafd9358', '2023-02-14T08:09:48.739212Z', '2023-01-29T14:24:07.182624Z', '2023-01-29T14:24:07.182624Z'),
    (x'018627cf1e8d2c5f4fdc91ffe01353de', 'Sophia Thompson', 'sophia@checkers.io', '$argon2id$v=19$m=65536,t=1,p=2$SXvnfVpYqOkm+wFPr4VlaA==$1wSremsvEP6Bc2LjzlncjQZvJrO3sBxTLP3AIAu8weQ=', 't', 't', 't', '2023-03-02T17:07:51.759218Z', '8ON1eiY4W5mUzs-l-ltBRj_yi3YzicJHAV5ARoIK1qE', x'aeff87bd8a3ed9f13d3b418174f760e26e94b4a2da678bd545015e400c2074e0290f99fedd8f5ed45d997d4caecf79df16817a3b154258cd338848baf502af5bbbffdd5f2427b58c2bf9b18c0ddd5c7a241532763c4eceefc806a5343b9634843e5667f8519f38ad68fcd7b9f6a0b24510a7f080d3d8a2da36bcc4ecdc2916f5', '2023-02-06T05:28:22.09Z', '2023-02-04T17:38:50.637565Z', '2023-02-04T17:38:50.637565Z'),
    (x'01862821852fab2f43c4b6345369190a', 'Robert Miller', 'bob@checkers.io', '$argon2id$v=19$m=65536,t=1,p=2$VFBxqapTyGeHhf4fkJkRjA==$QhYvYeDp/3KWCIkQDhrPCTvx2RLzPW0P2oc8ST9/Vgo=', 't', 't', 't', '2023-03-02T17:08:40.742653Z', 'cPZuO3RprHqH6pAeV_fdWWVVitULbOrfEngPhEhp9yg', x'38e16ae7f7de1ee4abe990984c10e96b2141b309bf941864a2b12d4bf69cd9d3986308d590d114a5790a066b0a2b51a5277f14ea5bc556f231617d2c9f50c9138eaa72a70f8546fbe69bdb89c926baea2c25a286f47727ab447dfa4dc483621b6d4cf7085d01b0ec60057a62e0e049d4a0cffefb4f15f5d6f19b9dc7863c4a8c', '2023-02-06T05:29:23.08Z', '2023-02-04T19:08:50.863919Z', '2023-02-04T19:08:50.863919Z')
;

-- valid invite from Edison to joe@checkers.io for checkers.io
-- valid invite from Jannel to eefrank@checkers.io for example.com
-- expired invite from Jannel to eefrank@checkers.io for example.com
-- bad token invite from Jannel to eefrank@checkers.io for example.com
INSERT INTO user_invitations (user_id, organization_id, role, email, expires, token, secret, created_by, created, modified) VALUES
    (x'018752ffc190971d220d5af7479fbbd4', x'0185df70923d78b6c3b03193999303bd', 'Member', 'joe@checkers.io', '2073-04-12T15:39:54.493794-05:00', 's6jsNBizyGh_C_ZsUSuJsquONYa-KH_2cmoJZd-jnIk', x'e57b56bb1d56b4a0c6cdea77e3b5a6de98db8ac94cc9035e3102e08009d242fba3f45bc557b226dd19c012d22c1404b1628ec401347a6c5dbb8e9470332de1c4a551827ab8bf5b9ccf2efb2e5ad0f108e72be671b09ba8193e93c6e4eaa48d8058ae5f71b0b336454a93cd046287bfcd02206b2a7a5fed80f18f2af251cbe088', x'0185df7210e5d8d7f6d33d7ecad99bbd', '2022-12-05T16:44:35.00123Z', '2022-12-05T16:44:35.00123Z'),
    (x'0185df7210e5d8d7f6d33d7ecad99bbd', x'0184e32c3c01f763f287d4a4f63c3293', 'Admin', 'eefrank@checkers.io', '2073-04-12T17:28:55.810443-05:00', 'pUqQaDxWrqSGZzkxFDYNfCMSMlB9gpcfzorN8DsdjIA', x'8f08ad6485e7e12da2d4efb9bdde929b34228e6a9028201586fe38ec3ee04a94cfde2f87bc2be3810d72d01b2a0f2675efc44cc015e7d49c3ddc6adc1cf3380960090a647dbffab0ffdfc7a93f20cac011b391889d0e98d1261602968e8ad7dae8fdb6e6e3cf474ebdac2d3b337f37d6bcac4ea98c72014a0576e87be6b932a2', x'0184e32cccecff01c1205fa8dc9e20df', '2022-12-05T16:44:35.00123Z', '2022-12-05T16:44:35.00123Z'),
    (x'0185df7210e5d8d7f6d33d7ecad99bbd', x'0184e32c3c01f763f287d4a4f63c3293', 'Admin', 'eefrank@checkers.io', '2020-04-12T17:28:55.810443-05:00', 's6jsNBizyGh_C_ZsUSuJsquONYa--gpcfzorN8DsdjIA', x'8f08ad6485e7e12da2d4efb9bdde929b34228e6a9028201586fe38ec3ee04a94cfde2f87bc2be3810d72d01b2a0f2675efc44cc015e7d49c3ddc6adc1cf3380960090a647dbffab0ffdfc7a93f20cac011b391889d0e98d1261602968e8ad7dae8fdb6e6e3cf474ebdac2d3b337f37d6bcac4ea98c72014a0576e87be6b932a2', x'0184e32cccecff01c1205fa8dc9e20df', '2022-12-05T16:44:35.00123Z', '2022-12-05T16:44:35.00123Z'),
    (x'0185df7210e5d8d7f6d33d7ecad99bbd', x'0184e32c3c01f763f287d4a4f63c3293', 'Admin', 'eefrank@checkers.io', '2073-04-12T17:28:55.810443-05:00', 'pUqQaDxWrqSGZzkxFDYNfCMSMlB--gpcfzorN8DsdjIA', x'8f08ad6485e7e12da2d4efb9bdde929b34228e6a9028201586fe38ec3ee04a94cfde2f87bc2be3810d72d01b2a0f2675efc44cc015e7d49c3ddc6adc1cf3380960090a647dbffab0ffdfc7a93f20cac011b391889d0e98d1261602968e8ad7dae8fdb6e6e3cf474ebdac2d3b337f37d6bcac4ea98c72014a0576e87be6b932a2', x'0184e32cccecff01c1205fa8dc9e20df', '2022-12-05T16:44:35.00123Z', '2022-12-05T16:44:35.00123Z')
;

INSERT INTO organization_users (organization_id, user_id, role_id, delete_confirmation_token, last_login, created, modified) VALUES
    (x'0184e32c3c01f763f287d4a4f63c3293', x'0184e32cccecff01c1205fa8dc9e20df', 1, 'g6JpZMQQAYTjLMzs/wHBIF+o3J4g36ZzZWNyZXTZQEd0b1d5b3UzTkdxYUNHVm5TbGtDM3RHRjQ4OFJFTDlyaWkyQjhpelNyWDVqV1JDYnFhMnhQc2FUTFlDWG9nNDSqZXhwaXJlc19hdNf/iQ6MQGJge5g', '2022-12-05T16:44:35.00123Z', '2022-12-05T16:44:35.00123Z', '2022-12-05T16:44:35.00123Z'),
    (x'0185df70923d78b6c3b03193999303bd', x'0185df7210e5d8d7f6d33d7ecad99bbd', 1, 'g6JpZMQQAYTjLMzs/wHBIF+o3J4g36ZzZWNyZXTZQFVLNWJZYXJvc3F2OGFJU29Tb0dWWlVQeUl0cFZzb2lnd3c2aUlDTEo3RnBsVUpmM3VNRG84eEZQOUFQclpxbzSqZXhwaXJlc19hdNf/BBZjAMJO4y4', '2023-02-23T16:24:32.741955Z', '2023-01-23T16:24:32.741955Z', '2023-01-23T16:24:32.741955Z'),
    (x'0184e32c3c01f763f287d4a4f63c3293', x'0185fde9f80e6439a2ee58de062a3e34', 1, 'g6JpZMQQAYX96fgOZDmi7ljeBio+NKZzZWNyZXTZQEdhZVVweTlhMUo4TDNqOXBWOW5zZ05nS0JNRTE0WjN2M204TFZ5YTNocktsZkN2OE80YkhQanFYamdZeDhTemGqZXhwaXJlc19hdNf/wql/AMJO9I0', NULL, '2023-01-29T14:24:07.182624Z', '2023-01-29T14:24:07.182624Z'),
    (x'0185df70923d78b6c3b03193999303bd', x'0185fde9f80e6439a2ee58de062a3e34', 3, NULL, '2023-01-29T14:24:07.182624Z', '2023-01-29T14:24:07.182624Z', '2023-01-29T14:24:07.182624Z'),
    (x'0185df70923d78b6c3b03193999303bd', x'018627cf1e8d2c5f4fdc91ffe01353de', 3, NULL, '2023-02-04T17:38:50.637565Z', '2023-02-04T17:38:50.637565Z', '2023-02-04T17:38:50.637565Z'),
    (x'0185df70923d78b6c3b03193999303bd', x'01862821852fab2f43c4b6345369190a', 4, NULL, '2023-02-04T19:08:50.863919Z', '2023-02-04T19:08:50.863919Z', '2023-02-04T19:08:50.863919Z')
;

INSERT INTO organization_projects (organization_id, project_id, created, modified) VALUES
    (x'0184e32c3c01f763f287d4a4f63c3293', x'0185cf6436b84d306262584b4c47beba', '2022-12-05T16:43:57.825256Z', '2022-12-05T16:43:57.825256Z'),
    (x'0184e32c3c01f763f287d4a4f63c3293', x'0185df804e85c8b399220570106ddd33', '2022-12-05T16:43:57.825256Z', '2022-12-05T16:43:57.825256Z'),
    (x'0185df70923d78b6c3b03193999303bd', x'0185df763d89b0f27bff56f3f6b934ef', '2023-01-23T16:22:54.781762Z', '2023-01-23T16:22:54.781762Z'),
    (x'0184e32c3c01f763f287d4a4f63c3293', x'0187bd8cc70a5d419013ff25baa99b7b', '2023-04-26T12:32:12.554166Z', '2023-04-26T12:32:12.554166Z'),
    (x'0185df70923d78b6c3b03193999303bd', x'0187be2faf3ff55e48424969bec9ebb6', '2023-04-26T15:30:08.831172Z', '2023-04-26T15:30:08.831172Z'),
    (x'0185df70923d78b6c3b03193999303bd', x'0187be318fdcc861e40d73196d6e0070', '2023-04-26T15:32:11.868714Z', '2023-04-26T15:32:11.868714Z'),
    (x'0185df70923d78b6c3b03193999303bd', x'0187be31d308f084323e692de7673ed4', '2023-04-26T15:32:29.064377Z', '2023-04-26T15:32:29.064377Z'),
    (x'0185df70923d78b6c3b03193999303bd', x'0187be321ce6d26c0333f7696aa4ed03', '2023-04-26T15:32:47.974258Z', '2023-04-26T15:32:47.974258Z'),
    (x'0185df70923d78b6c3b03193999303bd', x'0187be3269f649c3130721c3c70c6418', '2023-04-26T15:33:07.702712Z', '2023-04-26T15:33:07.702712Z'),
    (x'0185df70923d78b6c3b03193999303bd', x'0187be32b3217ab3a0c87e870fa6cbca', '2023-04-26T15:33:26.433043Z', '2023-04-26T15:33:26.433043Z'),
    (x'0185df70923d78b6c3b03193999303bd', x'0187be33003675c343a505e0d5165fd5', '2023-04-26T15:33:46.16672Z', '2023-04-26T15:33:46.16672Z'),
    (x'0185df70923d78b6c3b03193999303bd', x'0187be333b8289c9ef2cbd4d284128e2', '2023-04-26T15:34:01.346351Z', '2023-04-26T15:34:01.346351Z'),
    (x'0185df70923d78b6c3b03193999303bd', x'0187be337ed72ea1d961bfa1a6b32b51', '2023-04-26T15:34:18.583513Z', '2023-04-26T15:34:18.583513Z')
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
    (x'01851c016a56163161a693c1445f36e6', 14, '2022-12-05T16:48:21.123332Z', '2022-12-05T16:48:21.123332Z'),
    (x'01851c016a56163161a693c1445f36e6', 17, '2022-12-05T16:48:21.123332Z', '2022-12-05T16:48:21.123332Z'),
    (x'01851c016a56163161a693c1445f36e6', 18, '2022-12-05T16:48:21.123332Z', '2022-12-05T16:48:21.123332Z'),
    (x'01851c016a56163161a693c1445f36e6', 19, '2022-12-05T16:48:21.123332Z', '2022-12-05T16:48:21.123332Z'),
    (x'01851c016a56163161a693c1445f36e6', 20, '2022-12-05T16:48:21.123332Z', '2022-12-05T16:48:21.123332Z'),
    (x'0185df760de75905e97378ef96e402eb', 19, '2023-01-23T16:29:06.313194Z', '2023-01-23T16:29:06.313194Z'),
    (x'0185df7b44c7e5a52768317d59d6fc6a', 20, '2023-01-23T16:34:35.847672Z', '2023-01-23T16:34:35.847672Z'),
    (x'0185df7d23d0c2697cc3c1db4ad751ee', 14, '2023-01-23T16:36:38.480154Z', '2023-01-23T16:36:38.480154Z'),
    (x'0185df7d23d0c2697cc3c1db4ad751ee', 15, '2023-01-23T16:36:38.480154Z', '2023-01-23T16:36:38.480154Z'),
    (x'0185df7d23d0c2697cc3c1db4ad751ee', 16, '2023-01-23T16:36:38.480154Z', '2023-01-23T16:36:38.480154Z'),
    (x'0185df7d23d0c2697cc3c1db4ad751ee', 17, '2023-01-23T16:36:38.480154Z', '2023-01-23T16:36:38.480154Z'),
    (x'0185df863ba78d259e96e0801bb8c6d9', 19, '2023-01-23T16:46:34.407159Z', '2023-01-23T16:46:34.407159Z'),
    (x'0185df863ba78d259e96e0801bb8c6d9', 20, '2023-01-23T16:46:34.407159Z', '2023-01-23T16:46:34.407159Z'),
    (x'0185df88147d9ae4df140fee00101ff7', 19, '2023-01-23T16:48:35.453808Z', '2023-01-23T16:48:35.453808Z'),
    (x'0185df88147d9ae4df140fee00101ff7', 20, '2023-01-23T16:48:35.453808Z', '2023-01-23T16:48:35.453808Z'),
    (x'0185df89176e6c03dcdbb756dee9614b', 19, '2023-01-23T16:49:41.743053Z', '2023-01-23T16:49:41.743053Z'),
    (x'0185df89176e6c03dcdbb756dee9614b', 20, '2023-01-23T16:49:41.743053Z', '2023-01-23T16:49:41.743053Z'),
    (x'0185df89f353ff23957450984cf2dc3b', 19, '2023-01-23T16:50:38.035102Z', '2023-01-23T16:50:38.035102Z'),
    (x'0185df89f353ff23957450984cf2dc3b', 20, '2023-01-23T16:50:38.035102Z', '2023-01-23T16:50:38.035102Z'),
    (x'0185df8b0f3ca162573fb9517e31fbe5', 19, '2023-01-23T16:51:50.71668Z', '2023-01-23T16:51:50.71668Z'),
    (x'0185df8b0f3ca162573fb9517e31fbe5', 20, '2023-01-23T16:51:50.71668Z', '2023-01-23T16:51:50.71668Z'),
    (x'0185df8c144d84a9014b52d8b9abc29f', 19, '2023-01-23T16:52:57.549172Z', '2023-01-23T16:52:57.549172Z'),
    (x'0185df8c144d84a9014b52d8b9abc29f', 20, '2023-01-23T16:52:57.549172Z', '2023-01-23T16:52:57.549172Z'),
    (x'0185df8d01d077fd25527f459c0b1708', 19, '2023-01-23T16:53:58.352279Z', '2023-01-23T16:53:58.352279Z'),
    (x'0185df8d01d077fd25527f459c0b1708', 20, '2023-01-23T16:53:58.352279Z', '2023-01-23T16:53:58.352279Z'),
    (x'0185df8df3713992326989e477458b85', 19, '2023-01-23T16:55:00.209135Z', '2023-01-23T16:55:00.209135Z'),
    (x'0185df8df3713992326989e477458b85', 20, '2023-01-23T16:55:00.209135Z', '2023-01-23T16:55:00.209135Z'),
    (x'0185df8ed0ed50a503cc34c925d48bef', 19, '2023-01-23T16:55:56.9096Z', '2023-01-23T16:55:56.9096Z'),
    (x'0185df8ed0ed50a503cc34c925d48bef', 20, '2023-01-23T16:55:56.9096Z', '2023-01-23T16:55:56.9096Z'),
    (x'0185df8faa094394e70fc8698af079b8', 17, '2023-01-23T16:56:52.489425Z', '2023-01-23T16:56:52.489425Z'),
    (x'0185df8faa094394e70fc8698af079b8', 18, '2023-01-23T16:56:52.489425Z', '2023-01-23T16:56:52.489425Z'),
    (x'0185df8faa094394e70fc8698af079b8', 19, '2023-01-23T16:56:52.489425Z', '2023-01-23T16:56:52.489425Z'),
    (x'0185df8faa094394e70fc8698af079b8', 20, '2023-01-23T16:56:52.489425Z', '2023-01-23T16:56:52.489425Z')
;

---
title: "Consumer Groups"
weight: 5
date: 2023-05-17T17:03:41-04:00
---

Consumer groups allow multiple subscribers in different processes to coordinate how they consume events from a topic.

<!--more-->

## Consumer Group Identification

All subscribers must specify the same consumer group ID or name in order to join the same group. Consumer groups are stored in Ensign with a 16 byte ID; therefore in order to ensure uniqueness, users have the following options:

1. Specify a 16 byte ID directly (e.g. a UUID or a ULID)
2. Specify a name or an ID of any length that is unique to the project.

In the first case, Ensign will not modify the ID at all, guaranteeing its uniqueness. However, in the second case, Ensign will use a murmur3 128 bit hash to ensure that the computed ID is 16 bytes. If the ID is specified it is hashed, otherwise the name of the group is hashed. It is strongly recommended that a name string is used for the hash.

While the murmur3 hash does create the possibility of collisions, this will only happen for consumer groups in the same project (e.g. a consumer group with the same name in a different project will not cause a conflict). Therefore the probability is very low that a collision will occur. However, if you are creating a large number of consumer groups, it is generally better to use a UUID or ULID as the ID.
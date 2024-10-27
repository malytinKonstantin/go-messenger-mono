CREATE CONSTRAINT user_id_unique IF NOT EXISTS FOR (u:User) REQUIRE u.user_id IS UNIQUE;
CREATE CONSTRAINT request_id_unique IF NOT EXISTS FOR (r:FriendRequest) REQUIRE r.request_id IS UNIQUE;
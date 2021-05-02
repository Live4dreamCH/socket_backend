use chat;

-- register
-- insert into users (psw, u_name)
-- values (?, ?);

-- login
select psw, u_name
from users
where u_id=?;

-- fetch new msgs since last login
select m.msg_id, m.sender, m.msg_time, c.con_type , c.con
from msgs m, contents c
where m.s_id=? and m.msg_id>? and m.con_id=c.con_id

-- search user
-- select u_id, u_name
-- from users
-- where u_id=? or u_name=?;

-- has friend
select fr_id
from friends
where my_id=? and fr_id=?;

-- add friend
insert into friends (my_id,fr_id)
values (?, ?);

-- update friend list
select f.fr_id, u.u_name
from friends f, users u
where f.my_id=? and f.fr_id=u.u_id;

-- rm friend

-- create a conversation
insert into convs (s_name, owner)
values (null, null);

insert into conv_members (s_id, mem_id)
values (?, ?), (?, ?);

-- create a group
insert into convs (s_name, owner)
values (?, ?);

insert into conv_members (s_id, mem_id)
values (?, ?);

-- join a group
-- test
select sm.s_id
from conv_members sm
where sm.mem_id=? and sm.s_id=?;

insert into conv_members (s_id, mem_id)
values (?, ?);

-- a conv's member
select sm.mem_id
from conv_members sm
where sm.s_id=?;

-- add a msg
insert into contents (con_type, con)
values (?, ?);

insert into msgs (sender, msg_time, con_id, s_id)
values (?, now(), ?, ?);
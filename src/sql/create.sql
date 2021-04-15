create table users(
	u_id	int			not null	auto_increment,
	psw		varchar(51)	not null,
	u_name	varchar(31)	not null,
	primary key	(u_id)
);

create table friends(
	my_id	int	not null,
	fr_id	int not null,
	primary key (my_id, fr_id),
	foreign key (my_id) references users(u_id) on delete cascade,
	foreign key (fr_id) references users(u_id) on delete cascade
);

create table sessions(
	s_id	int			not null	auto_increment,
	s_name	varchar(31)	null,	-- 当会话是好友间会话，会话名为null
	owner	int 		null,	-- 当会话是好友间会话，会话所有者为null
	primary key (s_id),
	foreign key (owner) references users(u_id)
	-- 当会话所有者被删除，应当任意指定一名新群主or系统来当群主，而非解散整个群聊
	-- 类似linux对孤儿进程的处理：parent=init
);

create table session_members(
	s_id	int not null,
	mem_id	int not null,
	primary key (s_id, mem_id),
	foreign key (s_id) references sessions(s_id) on delete cascade,
	foreign key (mem_id) references users(u_id) on delete cascade
);

create table msgs(
	msg_id		int 		not null	auto_increment,
	sender		int 		not null,	-- 消息转发后, sender变为转发者的ID
	msg_time	datetime	not null,
	primary key (msg_id),
	foreign key (sender) references users(u_id)
	-- 出于安全因素，删除用户不能删除其已有聊天记录
);

create table session_msgs(
	msg_id	int	not null,
	s_id	int not null,
	primary key (msg_id),
	foreign key (msg_id) references msgs(msg_id) on delete cascade,
	foreign key (s_id) references sessions(s_id) on delete cascade
);

create table contents(
	con_id		int 			not null	auto_increment,
	con_type	tinyint			not null,
	con			varchar(1005)	not null,
	primary key (con_id)
);

create table msg_contents(
	msg_id	int	not null,
	con_id	int not null,
	primary key (msg_id),
	foreign key (msg_id) references msgs(msg_id) on delete cascade,
	foreign key (con_id) references contents(con_id) on delete cascade
);
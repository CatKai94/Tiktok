create table tiktok.comments
(
    id           bigint auto_increment comment '评论id，自增主键'
        primary key,
    user_id      bigint            not null comment '评论发布用户id',
    video_id     bigint            not null comment '评论视频id',
    comment_text varchar(255)      not null comment '评论内容',
    create_date  datetime          not null comment '评论发布时间',
    cancel       tinyint default 0 not null comment '默认评论发布为0，取消后为1'
)
    comment '评论表' charset = utf8mb3;

create index videoIdIdx
    on tiktok.comments (video_id)
    comment '评论列表使用视频id作为索引-方便查看视频下的评论列表';

create table tiktok.follows
(
    id          bigint auto_increment comment '自增主键'
        primary key,
    user_id     bigint            not null comment '用户id',
    follower_id bigint            not null comment '关注的用户',
    cancel      tinyint default 0 not null comment '默认关注为0，取消关注为1',
    constraint userIdToFollowerIdIdx
        unique (user_id, follower_id)
)
    comment '关注表' charset = utf8mb3;

create index FollowerIdIdx
    on tiktok.follows (follower_id);

create table tiktok.likes
(
    id       bigint auto_increment comment '自增主键'
        primary key,
    user_id  bigint            not null comment '点赞用户id',
    video_id bigint            not null comment '被点赞的视频id',
    cancel   tinyint default 0 not null comment '默认点赞为0，取消赞为1',
    constraint userIdtoVideoIdIdx
        unique (user_id, video_id)
)
    comment '点赞表' charset = utf8mb3;

create index userIdIdx
    on tiktok.likes (user_id);

create index videoIdx
    on tiktok.likes (video_id);

create table tiktok.messages
(
    id          int auto_increment
        primary key,
    user_id     int          null,
    receiver_id int          null,
    msg_content varchar(255) null,
    created_at  int          null,
    have_get    int          null
);

create table tiktok.users
(
    id       int auto_increment comment '自增ID'
        primary key,
    username varchar(30)  not null comment '账号',
    password varchar(100) not null comment '密码'
);

create table tiktok.videos
(
    id           bigint auto_increment comment '自增主键，视频唯一id'
        primary key,
    author_id    bigint       not null comment '视频作者id',
    play_url     varchar(255) not null comment '播放url',
    cover_url    varchar(255) not null comment '封面url',
    publish_time datetime     not null comment '发布时间戳',
    title        varchar(255) null comment '视频名称'
)
    comment '
视频表' charset = utf8mb3;

create index author
    on tiktok.videos (author_id);

create index time
    on tiktok.videos (publish_time);



-- To show the RDBMS structure of the demo app, I've elected to include the table scripts

CREATE TABLE [dbo].[Users] (
    [Id]           BIGINT        NOT NULL,
    [Email]        VARCHAR (250) NOT NULL,
    [Username]     VARCHAR (15)  NOT NULL,
    [FirstName]    VARCHAR (50)  NOT NULL,
    [LastName]     VARCHAR (50)  NOT NULL,
    [PasswordHash] VARCHAR (10)  NOT NULL,
    [DateOfBirth]  DATETIME      NOT NULL,
    PRIMARY KEY CLUSTERED ([Id] ASC)
);

CREATE TABLE [dbo].[Tweets] (
    [Id]     BIGINT        IDENTITY (1, 1) NOT NULL,
    [UserId] BIGINT        NOT NULL,
    [Tweet]  VARCHAR (280) NOT NULL,
    PRIMARY KEY CLUSTERED ([Id] ASC),
    CONSTRAINT [FK_Tweets_ToUsers] FOREIGN KEY ([UserId]) REFERENCES [dbo].[Users] ([Id])
);


CREATE TABLE [dbo].[Followers] (
    [UserId]         BIGINT NOT NULL,
    [FollowerUserId] BIGINT NOT NULL,
    CONSTRAINT [FK_Followers_ToUsers] FOREIGN KEY ([FollowerUserId]) REFERENCES [dbo].[Users]([Id]),
    CONSTRAINT [FK_Followers_ToUsers2] FOREIGN KEY ([UserId]) REFERENCES [dbo].[Users]([Id]),
    CONSTRAINT [UC_Followers_UserId_FollowerUserId] UNIQUE ([UserId], [FollowerUserId])
);


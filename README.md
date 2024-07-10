
# Distributed Social Media Application

Designed a distributed microservices social media application (like TWITTER)


## Tech Stack

**Client:** React, React Query, React Hook Form, NodeJs, ShadCn UI, Tailwind Css, Axios

**Server:** Golang, Postgres, Redis, RabbitMQ, Docker Compose, Kubernetes, WebSocket

**Cloud:** Avien Postgres Node, Avien Redis Node (Production Deployment Only)



## Installation

Make sure you have following services installed:

- Golang
  - [Install Golang ubuntu](https://www.digitalocean.com/community/tutorials/how-to-install-go-on-ubuntu-20-04)
- Docker & Docker compose
  - [Install Docker Compose ubuntu](https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-compose-on-ubuntu-20-04)
- Make
  - sudo apt-get install build-essential


Run DS_SOCIAL APP with docker compose

```bash
  cd DS_SOCIAL-APP
  cd project

 (Make Change in project/docker-compose.yml according to your requirement)

  make up_build
```
    
## Services

- User Service
- Tweet Service
- Timeline Service
- Notification Service
- User Relation Service


## Screenshots

#### CDC (Change Data Capture)

![App Screenshot](https://res.cloudinary.com/dgrxzxtd8/image/upload/v1720552032/x4ts4i4wru1tovosel7r.png)

My approch to create a CDC pipline which will update redis cache on following events of postgres database:

- **Post on Twitter** (Update user tweet sorted set, monitor user timeline, add new tweet to hashset using tweetId key)
- **Favorite/Dislike Post** (Update user like tweet ids set and tweet hash likeCount)
- **User Follow/Unfollow** (Update the follower's timeline with the last three tweets from the person they are following and user following set)
- **Revise User Information** (Update every user's tweet hash with their updated information.)


### Client Application

![ScreenShot 1](https://res.cloudinary.com/dgrxzxtd8/image/upload/v1720553014/i3blk1zgfueuahtubg3e.png)


**Client User Home Timeline**

![](https://res.cloudinary.com/dgrxzxtd8/image/upload/v1720553014/vmg1rrjeyfhrpreoxcmd.png)

**Single Tweet Page**

![](https://res.cloudinary.com/dgrxzxtd8/image/upload/v1720553014/jskump9xjxijx8us4w68.png)

**Reply Thread**

![](https://res.cloudinary.com/dgrxzxtd8/image/upload/v1720553916/sakdt9tfiui5n1jkv1ar.png)

**Tweet Thread Page**

![](https://res.cloudinary.com/dgrxzxtd8/image/upload/v1720553013/s7jmejoj6onwd3nwyorr.png)

**User Profile Page**

![](https://res.cloudinary.com/dgrxzxtd8/image/upload/v1720553014/oitqeaaqtafeacidftfc.png)

**User Like Posts**





## API Reference

### User Service API

#### Create a User

```http
  POST /api/user
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `name` | `string` | **Required**. User Name |
| `email` | `string` | **Required**. Unique User Email |
| `password` | `string` | **Required**. User Account Password |

#### Login a User

```http
  POST /api/user/login
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `email`      | `string` | **Required**. Unique User Email |
| `password` | `string` | **Required**. User Account Password |



Returns userId,userName,email and accesstoken for validate user.

#### Get a User details from token

```http
  GET /api/user/details
```

| Header | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `accessToken`      | `string` | **Required**. User JWT AccessToken |



Returns user complete details using decoded userId from the token.

#### Update a User details from token

```http
  PATCH /api/user/details
```

| Header | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `accessToken`      | `string` | **Required**. User JWT AccessToken |


| Body(FormData) | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `body`      | `FormData` | User can pass name,email and profileImage to update |

Validate user JWT token and update user details according to body.

![](https://res.cloudinary.com/dgrxzxtd8/image/upload/v1720594402/xlckswjpwrcvttanckmc.png)

### Tweet Service API

#### Post a tweet

```http
  POST /api/tweet
```

| Body | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `userId` | `string` | **Required**. Unique userId |
| `content` | `string` | **Required**. Tweet content with max 140 char |
| `mediaFiles` | `fileBuffer[]` | Media files array (images only) |

![](https://res.cloudinary.com/dgrxzxtd8/image/upload/v1720592935/fpmxcubnrj75gjo9f9v7.png)

#### Get All tweets

```http
  POST /api/tweet
```

Returns all tweets from database.

#### Get tweet details

```http
  GET /api/tweet/details/:id
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Unique tweet Id |



Returns complete tweet details with user information.

#### Like/Dislike a tweet

```http
  PATCH /api/tweet/action
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `tweetId`      | `string` | **Required**. Unique Tweet Id |
| `userId`      | `string` | **Required**. Unique User Id |
| `isLike`      | `boolean` | **Required**. Tweet like/dislike flag |


Return user tweet like/dislike response.

![](https://res.cloudinary.com/dgrxzxtd8/image/upload/v1720593766/wonlkpqpskep64aoq1xi.png)

#### Get User Id who liked tweet

```http
  GET /api/tweet/action/:id
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Unique Tweet Id |



Return all user ids who liked the tweet

#### Get all tweets liked by user

```http
  GET /api/tweet/user/action/:id
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Unique User Id |


Return all tweet Ids liked by user.

#### Post a reply

```http
  POST /api/tweet/reply
```

| Body | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `userId` | `string` | **Required**. Unique userId |
| `content` | `string` | **Required**. Tweet content with max 140 char |
| `mediaFiles` | `fileBuffer[]` | Media files array (images only) |
| `tweetId` | `string` | **Required**. Unique tweetId |



#### Get all replies of tweet

```http
  GET /api/tweet/reply/:id
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Unique reply Id |

| Query | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `page`      | `number` | **Required**. Page number (default 1) |
| `size`      | `number` | **Required**. Doc Size (default 10) |

Returns all replies on tweet with user information based on pagination

#### Like/Dislike a reply

```http
  PATCH /api/tweet/reply/action
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `userId`      | `string` | **Required**. Unique User Id |
| `isLike`      | `boolean` | **Required**. Tweet like/dislike flag |
| `replyId`      | `string` | **Required**. Unique Reply Id |


Return user reply like/dislike response.

#### Get all replies liked by user

```http
  GET /api/tweet/user/reply/action/:id
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Unique User Id |


Return all like replies by a user.


### TimeLine Service API

#### Get user following TimeLine

```http
  GET /api/timeline/following/:id
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Unique User Id |

| Query | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `page`      | `number` | **Required**. Page number (default 1) |
| `size`      | `number` | **Required**. Doc Size (default 10) |

Returns user following feed based on pagination.

#### Get user home TimeLine

```http
  GET /api/timeline/home/:id
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Unique User Id |

| Query | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `page`      | `number` | **Required**. Page number (default 1) |
| `size`      | `number` | **Required**. Doc Size (default 10) |

Returns random tweets home timeline to user.


#### Get user liked tweets 

```http
  GET /api/timeline/like/:id
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Unique User Id |



Returns all liked tweet from user.

#### Get user posted tweets 

```http
  GET /api/timeline/posts/:id
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Unique User Id |



Returns all user posted tweets.



### User Relation Service API

#### Follow a User

```http
  PATCH /api/relation/follow
```

| Body | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `followeeId`      | `string` | **Required**. Unique Following User Id |
| `followerId`      | `string` | **Required**. Unique Follower User Id |

Returns user following payload

#### Unfollow a User

```http
  PATCH /api/relation/unFollow
```

| Body | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `followeeId`      | `string` | **Required**. Unique Following User Id |
| `followerId`      | `string` | **Required**. Unique Follower User Id |

Returns user unFollowing payload

![](https://res.cloudinary.com/dgrxzxtd8/image/upload/v1720594200/wiog010pusyzv3xb45ol.png)

#### Get all following user

```http
  GET /api/relation/following/:id
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Unique User Id |


Returns all following user list of given user id.

#### Get all follower user ids

```http
  GET /api/relation/follower_id/:id
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Unique User Id |


Returns all follower user ids list of given user id.


#### Get all following user

```http
  GET /api/relation/following_id/:id
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Unique User Id |


Returns all following user ids list of given user id.

#### Get all follower user

```http
  GET /api/relation/recommend/:id
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Unique User Id |


Returns unfollowed user list for "Who to follow" section.


### Notification Service


![](https://res.cloudinary.com/dgrxzxtd8/image/upload/v1720595376/brrtv654tn37hyjk6qpb.png)

A websocket-based **Notification service** was developed, and it would alert all active users in response to a client-generated websocket event trigger.

Notification The Websocket Service offers the subsequent events:

- **POST TWEET**: Send out a notification of the author's latest tweet to all currently following users.
- **LIKE TWEET**: Notify the active author when a user likes one of their posts.
- **FOLLOW USER**: Notify the user who is an active follower about the next event.
